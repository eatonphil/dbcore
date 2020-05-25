module Database

open Npgsql


type Column =
    {
        Name: string
        Type: string
    }


type Constraint =
    {
        Column: string
        ForeignTable: string
        ForeignColumn: string
    }


type Table =
    {
        Name: string
        Columns: Column[]
        ForeignKeys: Constraint[]
        PrimaryKey: Option<Constraint>
    }


[<NoComparison>]
type DatabaseReader =
    {
        Cfg: Config.DatabaseConfig
    }

    member private this.GetConn() : NpgsqlConnection =
        let dsn =
            sprintf "Host=%s;Port=%s;Database=%s;Username=%s;Password=%s;"
                this.Cfg.Host
                this.Cfg.Port
                this.Cfg.Database
                this.Cfg.Username
                this.Cfg.Password
        let conn = new NpgsqlConnection(dsn)
        conn.Open()
        conn

    member private this.GetConstraints(conn: NpgsqlConnection, table: string, typ: string) : Constraint[] =
        let query =
            "SELECT
                 kcu.column_name,
                 ccu.table_name,
                 ccu.column_name,
             FROM
                 information_schema.table_constraints AS tc
                 JOIN information_schema.key_column_usage AS kcu
                   ON tc.constraint_name = kcu.constraint_name
                   AND tc.table_schema = kcu.table_schema
                 JOIN information_schema.constraint_column_usage AS ccu
                   ON ccu.constraint_name = tc.constraint_name
                   AND ccu.table_schema = tc.table_schema
             WHERE tc.constraint_type = '@type' AND tc.table_name = '@name'"
        use cmd = new NpgsqlCommand(query, conn)
        cmd.Parameters.AddWithValue("name", table)
        cmd.Parameters.AddWithValue("type", typ)
        use dr = cmd.ExecuteReader()
        [|
            while dr.Read() do
                yield {
                    Column = dr.GetString 0
                    ForeignTable = dr.GetString 1
                    ForeignColumn = dr.GetString 2
                }
        |]

    member private this.GetTable(conn: NpgsqlConnection, name: string) : Table =
        let columns =
            let query =
                "SELECT
                     column_name, data_type
                 FROM
                     information_schema.columns
                 WHERE
                     table_schema='public' AND table_name='@name'"
            use cmd = new NpgsqlCommand(query, conn)
            cmd.Parameters.AddWithValue("name", name)
            use dr = cmd.ExecuteReader()
            [| while dr.Read() do
                   yield { Name = dr.GetString 0; Type = dr.GetString 1 } |]

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

    member this.GetTables () : Table[] =
        use conn = this.GetConn()
        let query =
            "SELECT
                 table_name
             FROM
                 information_schema.tables
             WHERE
                 table_schema='public'"
        use cmd = new NpgsqlCommand(query, conn)
        use dr = cmd.ExecuteReader()

        [| while dr.Read() do yield this.GetTable(conn, dr.GetString 0) |]


let MakeDatabaseReader (cfg: Config.DatabaseConfig) : DatabaseReader =
    { Cfg = cfg }
