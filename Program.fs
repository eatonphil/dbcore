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

    let template = Template.MakeEngine(
                       Path.Combine("templates", config.Api.Language),
                       Path.Combine(projectDir, config.Api.Outdir))
    template.Write({| Tables = tables; Api = config.Api |})

    0
