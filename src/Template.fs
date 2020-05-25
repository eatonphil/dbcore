module Template

open System.IO
open System.Text.RegularExpressions

open Scriban


let rec private getFiles(dir: string) : seq<string> =
    seq {
        yield! Directory.EnumerateFiles(dir, "*.*")
        for d in Directory.EnumerateDirectories(dir) do
            yield! getFiles(d)
    }


type Engine =
    {
        SourceDir: string
        OutDir: string
    }

    member this.Write(ctx: {| Tables: Database.Table[]; Api: Config.ApiConfig |}) =
        for f in getFiles(this.SourceDir) do
            let tpl = Template.Parse(File.ReadAllText(f), f)

            // Drop the SourceDir/ prefix
            let f = f.Substring(this.SourceDir.Length + 1)
            // Handle the special case where files should be enumerated per table
            let tableSubstitute = "GENAPP__tables__"
            let fsAndCtxs =
                if not (f.Contains tableSubstitute) then [(f, {| ctx with Table = ctx.Tables.[0] |})]
                else [ for t in ctx.Tables do
                           yield (Path.GetFileName(f).Replace(tableSubstitute, t.Name),
                                  {| ctx with Table = t |}) ]
            for (f, ctx) in fsAndCtxs do
                let outFile = Path.Combine(this.OutDir, f)
                printfn "[DEBUG] Generating: %s" outFile

                // Create directory if not exists
                (new FileInfo(outFile)).Directory.Create()

                File.WriteAllText(outFile, tpl.Render(ctx))


let MakeEngine(sourceDir: string, outDir: string) : Engine =
    {
        SourceDir = sourceDir
        OutDir = outDir
    }
