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

    member this.Write(ctx: obj) =
        for f in getFiles(this.SourceDir) do
            let tpl = Template.Parse(File.ReadAllText(f), f)

            let templateFileName = Regex.Replace(f, "GENAPP__([a-zA-Z0-9]*)__(.*)", "{{ for a in $1 }}{{ a }}$2\n{{ end }}")
            let fs = Template.Parse(templateFileName).Render(ctx).Split("\n")
            for f in fs do
                let outFile = Path.Combine(this.OutDir, f)
                printfn "[DEBUG] Generating: %s" outFile
                // Create directory if not exists
                (new FileInfo(outFile)).Directory.Create()
                //File.WriteAllText(outFile, tpl.Render(ctx))


let MakeEngine(sourceDir: string, outDir: string) : Engine =
    {
        SourceDir = sourceDir
        OutDir = outDir
    }
