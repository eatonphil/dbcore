open System.IO


[<EntryPoint>]
let main (args: string []): int =
    let projectDir = if args.Length > 0
                         then args.[0]
                         else failwith "Expected project directory"

    let cfg = Config.GetConfig(Path.Combine(projectDir, "dbcore.yml"))
    cfg.Validate()

    if cfg.Database.Dialect = "sqlite" then
        cfg.Database.Database <- Path.Combine(projectDir, cfg.Database.Database)
    let db = Reader.Reader(cfg.Database)
    let tables = db.GetTables()

    let ctx: Template.Context = {
        Project = cfg.Project
        Database = {| Dialect = cfg.Database.Dialect |}
        Api = cfg.Api
        Browser = cfg.Browser
        Tables = tables
        Template = ""
        OutDir = ""
        CultureName = cfg.CultureName
    }
    Template.GenerateApi(projectDir, cfg.Api, ctx)
    Template.GenerateBrowser(projectDir, cfg.Browser, ctx)

    0
