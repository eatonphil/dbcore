namespace Reader

open System.Data

open Npgsql
open MySql.Data.MySqlClient

open Database


// This is for SQL implementations that follow the ANSI standard of
// using information_schema tables to store database metadata.
[<NoComparison>]
type InformationSchema(cfg0: Config.DatabaseConfig) =
    let cfg = cfg0
    let connFactory =
        match cfg.Dialect with
            | "mysql" -> fun (dsn) -> new MySqlConnection(dsn) :> IDbConnection
            | "postgres" -> fun (dsn) -> new NpgsqlConnection(dsn) :> IDbConnection
            | d -> failwith ("Unsupported SQL dialect: " + d)

    let getConn() : IDbConnection =
        let dsn =
            sprintf "Host=%s;Port=%s;Database=%s;Username=%s;Password=%s;"
                cfg.Host
                cfg.Port
                cfg.Database
                cfg.Username
                cfg.Password
        let conn = connFactory(dsn)
        conn.Open()
        conn

    let addStringParameter(cmd: IDbCommand, name: string, value: string) : unit =
        let p = cmd.CreateParameter()
        p.ParameterName <- name
        p.Value <- value
        cmd.Parameters.Add(p) |> ignore

    let getConstraints(conn: IDbConnection, table: string, typ: string) : Constraint[] =
        let query =
            match cfg.Dialect with
                | "postgres" ->
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
                | "mysql" ->
                    "SELECT
                         column_name,
                         COALESCE(referenced_table_name, ''),
                         COALESCE(referenced_column_name, '')
                     FROM
                         information_schema.table_constraints tc
                         JOIN information_schema.key_column_usage kcu
                             ON kcu.table_schema = tc.table_schema
                             AND kcu.table_name = tc.table_name
                             AND kcu.constraint_name = tc.constraint_name
                     WHERE tc.constraint_type=@type AND tc.table_name=@name"
                 | d -> failwith "Unknown dialect: " + d
        use cmd = conn.CreateCommand()
        cmd.CommandText <- query

        addStringParameter(cmd, "name", table)
        addStringParameter(cmd, "type", typ)
        cmd.Prepare()
        use dr = cmd.ExecuteReader()
        [|
            while dr.Read() do
                yield {
                    Column = dr.GetString(0)
                    Type = ""
                    ForeignTable = dr.GetString(1)
                    ForeignColumn = dr.GetString(2)
                }
        |]

    let getTable(conn: IDbConnection, name: string) : Table =
        let columns =
            let autoIncrementCheck =
                match cfg.Dialect with
                    | "postgres" -> "COALESCE(column_default, '') LIKE '%seq%' -- Not the greatest way to detect auto incrementing, but ok for now"
                    | "mysql" -> "extra LIKE '%auto_increment%'"
                    | d -> failwith ("Unsupported SQL dialect: " + d)
            let query =
                sprintf "SELECT
                     column_name,
                     data_type,
                     is_nullable <> 'NO',
                     %s
                 FROM
                     information_schema.columns
                 WHERE
                     table_schema=@schema AND table_name=@name" autoIncrementCheck
            use cmd = conn.CreateCommand()
            cmd.CommandText <- query
            addStringParameter(cmd, "name", name)
            addStringParameter(cmd, "schema", cfg.Schema)
            cmd.Prepare()
            use dr = cmd.ExecuteReader()
            [|
                while dr.Read() do
                    yield {
                        Name = dr.GetString(0)
                        Type = dr.GetString(1).ToLower()
                        Nullable = dr.GetBoolean(2)
                        AutoIncrement = dr.GetBoolean(3)
                    }
            |]

        if columns.Length = 0
            then failwith ("Expected more than 0 columns in table " + name)
            else "" |> ignore

        let foreignKeys = getConstraints(conn, name, "FOREIGN KEY")

        let primaryKey =
            let keys = getConstraints(conn, name, "PRIMARY KEY")
            if keys.Length = 0 then None
            else
                let typ = [ for c in columns do
                                if c.Name = keys.[0].Column then
                                    yield c.Type ].[0]
                Some ({ keys.[0] with Type = typ.ToLower() })

        {
            Name = name
            Columns = columns
            ForeignKeys = foreignKeys
            PrimaryKey = primaryKey
        }

    member this.GetTables() : Table[] =
        use conn = getConn()
    
        // Fetch names first so cmd can be closed
        let tableNames =
            let query =
                "SELECT
                    table_name
                 FROM
                     information_schema.tables
                 WHERE
                     table_schema=@schema"
            use cmd = conn.CreateCommand()
            cmd.CommandText <- query
            addStringParameter(cmd, "schema", cfg.Schema)

            use dr = cmd.ExecuteReader()
            [| while dr.Read() do yield dr.GetString 0 |]
    
        [| for name in tableNames do yield getTable(conn, name) |]
