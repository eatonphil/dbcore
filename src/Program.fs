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
    let tables: Template.Table[] =
        let notAuditOrAutoIncrement(c: Database.Column) : bool =
            if c.AutoIncrement then false
                else if not cfg.Api.Audit.Enabled then true
                    else if (c.Name = cfg.Api.Audit.CreatedAt ||
                             c.Name = cfg.Api.Audit.UpdatedAt ||
                             c.Name = cfg.Api.Audit.DeletedAt) then false
                        else true

        let notAutoIncrement(c: Database.Column) : bool = not c.AutoIncrement

        [|
            for table in db.GetTables() do
                let mutable label = table.Name
                for t in cfg.Database.Tables do
                    if t.Name = table.Name then label <- t.Label

                let columnsNoAudit = table.Columns |> Array.filter notAuditOrAutoIncrement
                let columnsNoAutoIncrement = table.Columns |> Array.filter notAutoIncrement

                let t: Template.Table = {
                    Label = label
                    Name = table.Name
                    Columns = table.Columns
                    ForeignKeys = table.ForeignKeys
                    PrimaryKey = table.PrimaryKey

                    ColumnsNoAutoIncrement = columnsNoAutoIncrement
                    ColumnsNoAudit = columnsNoAudit
                }
                t
        |]

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
