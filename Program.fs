open System.IO

open Database


[<EntryPoint>]
let main (args: string []): int =
    let projectDir = if args.Length > 0
                         then args.[0]
                         else failwith "Expected project directory"

    // TODO: validate file
    let config = Config.GetConfig(Path.Combine(projectDir, "genapp.yml"))

    let db = Database.MakeDatabaseReader(config.Database)
    let tables = db.GetTables()

    let ctx: Template.Context = {
        Project = config.Project
        Api = config.Api
        Browser = config.Browser
        Tables = tables
    }
    Template.Generate(projectDir, config.Api, ctx)
    Template.Generate(projectDir, config.Browser, ctx)

    0
