module Database

open Npgsql


type Column =
    {
        Name: string
        Type: string
    }


type Table =
    {
        Name: string
        Columns: Column[]
    }


[<NoComparison>]
type DatabaseReader =
    {
        Cfg: Config.DatabaseConfig
    }

    member private this.GetConn() : NpgsqlConnection =
        let dsn = sprintf
                      "Host=%s;Port=%s;Database=%s;Username=%s;Password=%s;"
                      this.Cfg.Host
                      this.Cfg.Port
                      this.Cfg.Database
                      this.Cfg.Username
                      this.Cfg.Password
        let conn = new NpgsqlConnection(dsn)
        conn.Open()
        conn

    member private this.GetTable(conn: NpgsqlConnection, name: string) : Table =
        let query = sprintf "SELECT column_name, data_type
                             FROM information_schema.columns
                             WHERE table_schema='public' AND table_name='%s'" name
        let cmd = new NpgsqlCommand(query, conn)
        let dr = cmd.ExecuteReader()
        let columns = [| while dr.Read() do yield { Name = dr.GetString 0; Type = dr.GetString 1 } |]
        dr.Close()
        { Name = name; Columns = columns }

    member this.GetTables () : Table[] =
        let conn = this.GetConn()
        let query = "SELECT table_name
                     FROM information_schema.tables
                     WHERE table_schema='public'"
        let cmd = new NpgsqlCommand(query, conn)
        let dr = cmd.ExecuteReader()
        let tableNames = [| while dr.Read() do yield dr.GetString 0 |]
        dr.Close()

        let tables = [| for name in tableNames do yield this.GetTable(conn, name) |]
        conn.Close()

        tables


let MakeDatabaseReader (cfg: Config.DatabaseConfig) : DatabaseReader =
    { Cfg = cfg }
