open System.IO


[<EntryPoint>]
let main (args: string []): int =
    let projectDir = if args.Length > 0
                         then args.[0]
                         else failwith "Expected project directory"

    // TODO: validate file
    let cfg = Config.GetConfig(Path.Combine(projectDir, "dbcore.yml"))

    let db = Database.Reader(cfg.Database)
    let tables = db.GetTables()

    let ctx: Template.Context = {
        Project = cfg.Project
        Api = cfg.Api
        Browser = cfg.Browser
        Tables = tables
    }
    Template.GenerateApi(projectDir, cfg.Api, ctx)
    Template.GenerateBrowser(projectDir, cfg.Browser, ctx)

    0
