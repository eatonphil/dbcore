namespace Database

open Npgsql


[<NoComparison>]
type PostgresReader(cfg0: Config.DatabaseConfig) =
    let cfg = cfg0

    member private this.GetConn() : NpgsqlConnection =
        let dsn =
            sprintf "Host=%s;Port=%s;Database=%s;Username=%s;Password=%s;"
                cfg.Host
                cfg.Port
                cfg.Database
                cfg.Username
                cfg.Password
        let conn = new NpgsqlConnection(dsn)
        conn.Open()
        conn

    member private this.GetConstraints(conn: NpgsqlConnection, table: string, typ: string) : Constraint[] =
        let query =
            "SELECT
                 kcu.column_name,
                 ccu.table_name,
                 ccu.column_name
             FROM
                 information_schema.table_constraints AS tc
                 JOIN information_schema.key_column_usage AS kcu
                   ON tc.constraint_name = kcu.constraint_name
                   AND tc.table_schema = kcu.table_schema
                 JOIN information_schema.constraint_column_usage AS ccu
                   ON ccu.constraint_name = tc.constraint_name
                   AND ccu.table_schema = tc.table_schema
             WHERE tc.constraint_type=@type AND tc.table_name=@name"
        use cmd = new NpgsqlCommand(query, conn)
        cmd.Parameters.AddWithValue("name", table) |> ignore
        cmd.Parameters.AddWithValue("type", typ) |> ignore
        cmd.Prepare()
        use dr = cmd.ExecuteReader()
        [|
            while dr.Read() do
                yield {
                    Column = dr.GetString 0
                    ForeignTable = dr.GetString 1
                    ForeignColumn = dr.GetString 2
                }
        |]

    member private this.SqlToGoType(typ: string, nullable: bool) : string =
        match typ with
        | "int" -> if nullable then "*int32" else "int32"
        | "bigint" -> if nullable then "null.Int" else "int64"
        | "text" | "varchar" | "char" -> if nullable then "null.String" else "string"
        | "bool" -> if nullable then "null.Bool" else "bool"
        | "timestamp" -> if nullable then "null.Time" else "time.Time"
        | _ -> failwith ("Unsupported PostgreSQL type" + typ)

    member private this.GetTable(conn: NpgsqlConnection, name: string) : Table =
        let columns =
            let query =
                "SELECT
                     column_name,
                     data_type,
                     is_nullable <> 'NO',
                     COALESCE(column_default, '') LIKE '%seq%' -- Not the greatest way to detect auto incrementing, but ok for now
                 FROM
                     information_schema.columns
                 WHERE
                     table_schema='public' AND table_name=@name"
            use cmd = new NpgsqlCommand(query, conn)
            cmd.Parameters.AddWithValue("name", name) |> ignore
            cmd.Prepare()
            use dr = cmd.ExecuteReader()
            [|
                while dr.Read() do
                    yield {
                        Name = dr.GetString 0
                        Type = dr.GetString 1
                        GoType = this.SqlToGoType (dr.GetString 1, dr.GetBoolean 2)
                        AutoIncrement = dr.GetBoolean 3
                    }
            |]

        if columns.Length = 0
            then failwith ("Expected more than 0 columns in table " + name)
            else "" |> ignore

        let foreignKeys = this.GetConstraints(conn, name, "FOREIGN KEY")

        let primaryKey =
            let keys = this.GetConstraints(conn, name, "PRIMARY KEY")
            if keys.Length > 0
                then Some (keys.[0])
                else None

        {
            Name = name
            Columns = columns
            ForeignKeys = foreignKeys
            PrimaryKey = primaryKey
        }

    member this.GetTables() : Table[] =
        use conn = this.GetConn()
    
        // Fetch names first so cmd can be closed
        let tableNames =
            let query =
                "SELECT
                table_name
             FROM
                 information_schema.tables
             WHERE
                 table_schema='public'"
            use cmd = new NpgsqlCommand(query, conn)
            use dr = cmd.ExecuteReader()
            [| while dr.Read() do yield dr.GetString 0 |]
    
        [| for name in tableNames do yield this.GetTable(conn, name) |]
