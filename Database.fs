type Column =
    {
        Name: string
        Type: string
    }

type Table =
    {
        Name: string
        Columns: Column []
    }

type DatabaseReader =
    {
        Host: string
        Database: string
        Username: string
        Password: string
    }

    method this.getConn () : NpgsqlConnection =
        let dsn = sprintf "Host=%s;Database=%s;Username=%s;Password=%s;" this.Host this.Database this.Username this.Password
        let conn = new NpgsqlConnection(dsn)
        conn.Open()
        let tables = getTables(conn)
        conn.Close()

    method this.getTables () : Column [] =
        let conn = this.getConn()
        let cmd = new NpgsqlCommand("SELECT table_name FROM information_schema.tables WHERE table_schema='public'", conn)
        let dr = cmd.ExecuteReader()
        let tables = [| while dr.Read() do yield dr.GetString 0 |]
        dr.Close()

        for name in tables do
            let cmd = new NpgsqlCommand(sprintf "SELECT column_name, data_type FROM information_schema.columns WHERE table_schema='public' AND table_name='%s'" name, conn)
            let dr = cmd.ExecuteReader()
            let columns = [| while dr.Read() do yield { Name = dr.GetString 0; Type = dr.GetString 1 } |]
            dr.Close()
            yield { Name = name; Columns = columns; }

        conn.Close()
