namespace Database


type Reader(cfg0: Config.DatabaseConfig) =
    let cfg = cfg0

    member this.GetTables() : Table[] =
        match cfg.Dialect with
            | "postgres" -> PostgresReader(cfg).GetTables()
            | "mysql" -> MySQLReader(cfg).GetTables()
            | d -> failwith ("Unsupported SQL dialect: " + d)

