namespace Database


[<NoComparison>]
type MySQLReader(cfg0: Config.DatabaseConfig) =
    let cfg = cfg0

    member this.GetTables() : Table[] = [| |]
