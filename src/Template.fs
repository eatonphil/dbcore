module Template

open System.Diagnostics
open System.IO
open System.Text.RegularExpressions


let rec private getFiles(dir: string) : seq<string> =
    seq {
        yield! Directory.EnumerateFiles(dir, "*.*")
        for d in Directory.EnumerateDirectories(dir) do
            yield! getFiles(d)
    }


[<NoComparison>]
type Context =
    {
        Project: string
        Database: {| Dialect: string |}
        Tables: Database.Table[]
        Api: Config.ApiConfig
        Browser: Config.BrowserConfig
        OutDir: string
        Template: string
        CultureName: string
    }


let private writeProjectToDisk(sourceDir: string, outDir: string, ctx: Context) =
    let ti = System.Globalization.CultureInfo(ctx.CultureName, false).TextInfo

    // Register helpers
    let helpers =
        [
            ("dbcore_capitalize", System.Converter<string, string>(fun a -> ti.ToTitleCase(a)));
        ]
    let tctx = Scriban.TemplateContext()
    let so = Scriban.Runtime.ScriptObject()
    for (name, fn) in helpers do
        Scriban.Runtime.ScriptObjectExtensions.Import(so, name, fn)
    tctx.PushGlobal(so)

    // Copy/render each template
    for f in getFiles(sourceDir) do
        let tpl = Scriban.Template.Parse(File.ReadAllText(f), f)

        // Drop the SourceDir/ prefix
        let f = f.Substring(sourceDir.Length + 1)

        // Order by substring length descending
        let pathTemplateExpansions =
            [
                ("tables_capitalize",
                 fun (path: string, sub: string) ->
                     [ for t in ctx.Tables do
                           yield (path.Replace(sub, ti.ToTitleCase(t.Name)),
                                  {| ctx with Table = t |}) ]);
                ("tables",
                 fun (path: string, sub: string) ->
                     [ for t in ctx.Tables do
                           yield (path.Replace(sub, t.Name),
                                  {| ctx with Table = t |}) ]);
                ("",
                 fun (path: string, sub: string) ->
                     [(path, {| ctx with Table = ctx.Tables.[0] |})])
            ]

        let mutable fsAndCtxs = []
        for (path, expand) in pathTemplateExpansions do
            let escapedPath = sprintf "DBCORE__%s__" path
            // Only expand once on the longest match.
            if fsAndCtxs.Length = 0 && (f.Contains(escapedPath) || path = "") then
                fsAndCtxs <- expand(f, escapedPath)

        for (f, ctx) in fsAndCtxs do
            let outFile = Path.Combine(outDir, f)
            printfn "[DEBUG] Generating: %s" outFile

            // Create directory if not exists
            FileInfo(outFile).Directory.Create()

            tctx.PushGlobal(Scriban.Runtime.ScriptObject.From(ctx))
            File.WriteAllText(outFile, tpl.Render(tctx))


let private generate(templateDir: string, projectDir: string, cfg: Config.IConfig, ctx: Context) =
    // Required for distribution to get right base path (especially within single file executable)
    let baseDir = System.AppContext.BaseDirectory
    let sourceDir = Path.Combine(baseDir, templateDir, cfg.Template)
    let outDir = Path.Combine(projectDir, cfg.OutDir)
    writeProjectToDisk(sourceDir, outDir, { ctx with OutDir=cfg.OutDir; Template=cfg.Template })

    printfn "[DEBUG] Running post install script: %s"
        (Path.Combine(outDir, "scripts/post-generate.sh"))
    let processInfo = new ProcessStartInfo(
                          FileName = "bash",
                          Arguments = "scripts/post-generate.sh",
                          WorkingDirectory = outDir)
    use p = Process.Start(processInfo)
    p.WaitForExit()


let GenerateApi(projectDir: string, cfg: Config.IConfig, ctx: Context) =
    generate("templates/api", projectDir, cfg, ctx)


let GenerateBrowser(projectDir: string, cfg: Config.IConfig, ctx: Context) =
    generate("templates/browser", projectDir, cfg, ctx)
