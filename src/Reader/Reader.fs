namespace Reader


type Reader(cfg0: Config.DatabaseConfig) =
    let cfg = cfg0

    member this.GetTables() : Database.Table[] =
        match cfg.Dialect with
            | "mysql" | "postgres" -> InformationSchema(cfg).GetTables()
            | d -> failwith ("Unsupported SQL dialect: " + d)

