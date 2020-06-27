namespace Reader

open System.Data

open Microsoft.Data.Sqlite

open Database

type SQLite(cfg0: Config.DatabaseConfig) =
    let cfg = cfg0

    let getForeignKeys(conn: SqliteConnection, name: string) : Constraint[] =
        use cmd = conn.CreateCommand()
        cmd.CommandText <- """SELECT "from", "table", "to" FROM pragma_foreign_key_list($schema)"""
        let _ = cmd.Parameters.AddWithValue("$schema", cfg.Schema)
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

    let getTable(conn: SqliteConnection, name: string) : Table =
        let mutable primaryKey : Option<Constraint> = None
        let columns =
            use cmd = conn.CreateCommand()
            cmd.CommandText <- """SELECT
                     name,
                     type,
                     "notnull",
                     pk = (SELECT 1 FROM sqlite_master WHERE tbl_name=$name AND sql LIKE '%%AUTOINCREMENT%%'),
                     pk
                 FROM
                     pragma_table_info($name)"""
            let _ = cmd.Parameters.AddWithValue("$name", name)
            cmd.Prepare()
            use dr = cmd.ExecuteReader()
            [|
                while dr.Read() do
                    let name = dr.GetString(0)
                    let typ = dr.GetString(1).ToLower()
                    let pk = dr.GetBoolean(4)
                    if pk then
                        primaryKey <- Some {
                            Column = name
                            Type = typ
                            ForeignTable = ""
                            ForeignColumn = ""
                        }

                    yield {
                        Name = name
                        Type = typ
                        Nullable = not (pk || dr.GetBoolean 2)
                        AutoIncrement = dr.GetBoolean(3)
                    }
            |]

        if columns.Length = 0
            then failwith ("Expected more than 0 columns in table " + name)
            else "" |> ignore

        let foreignKeys = getForeignKeys(conn, name)

        {
            Name = name
            Label = name
            Columns = columns
            ForeignKeys = foreignKeys
            PrimaryKey = primaryKey
        }

    member this.GetTables() : Table[] =
        use conn = new SqliteConnection("Data Source=" + cfg.Database)
        conn.Open()
        let tableNames =
            use cmd = conn.CreateCommand()
            cmd.CommandText <- "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'"

            use dr = cmd.ExecuteReader()
            [| while dr.Read() do yield dr.GetString 0 |]

        [| for name in tableNames do yield getTable(conn, name) |]
