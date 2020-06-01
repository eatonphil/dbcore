module Template

open System.Diagnostics
open System.IO
open System.Text.RegularExpressions

open Scriban


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
    }


let private writeProjectToDisk(sourceDir: string, outDir: string, ctx: Context) =
    for f in getFiles(sourceDir) do
        let tpl = Template.Parse(File.ReadAllText(f), f)

        // Drop the SourceDir/ prefix
        let f = f.Substring(sourceDir.Length + 1)
        // Handle the special case where files should be enumerated per table
        let tableSubstitute = "DBCORE__tables__"
        let fsAndCtxs =
            if not (f.Contains tableSubstitute) then [(f, {| ctx with Table = ctx.Tables.[0] |})]
            else [ for t in ctx.Tables do
                       yield (f.Replace(tableSubstitute, t.Name),
                              {| ctx with Table = t |}) ]
        for (f, ctx) in fsAndCtxs do
            let outFile = Path.Combine(outDir, f)
            printfn "[DEBUG] Generating: %s" outFile

            // Create directory if not exists
            FileInfo(outFile).Directory.Create()

            File.WriteAllText(outFile, tpl.Render(ctx))


let private generate(templateDir: string, projectDir: string, cfg: Config.IConfig, ctx: Context) =
    // Required for distribution to get right base path (especially within single file executable)
    let baseDir = System.AppContext.BaseDirectory
    let sourceDir = Path.Combine(baseDir, templateDir, cfg.Template)
    let outDir = Path.Combine(projectDir, cfg.Template)
    writeProjectToDisk(sourceDir, outDir, ctx)

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
