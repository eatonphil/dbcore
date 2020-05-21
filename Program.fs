open System

open Npgsql

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

[<EntryPoint>]
let main args =
    let dsn = sprintf "Host=%s;Database=%s;Username=%s;Password=%s;" (args.[0]) (args.[1]) (args.[2]) (args.[3])
    let conn = new NpgsqlConnection(dsn)
    conn.Open()
    let cmd = new NpgsqlCommand("SELECT table_name FROM information_schema.tables WHERE table_schema='public'", conn)
    let dr = cmd.ExecuteReader()
    let tables = [| while dr.Read() do yield dr.GetString 0 |]
    dr.Close()

    for name in tables do
        let cmd = new NpgsqlCommand(sprintf "SELECT column_name, data_type FROM information_schema.columns WHERE table_schema='public' AND table_name='%s'" name, conn)
        let dr = cmd.ExecuteReader()
        let columns = [| while dr.Read() do yield { Name = dr.GetString 0; Type = dr.GetString 1 } |]
        let table: Table = { Name = name; Columns = columns; }
        printfn "%A" table
        dr.Close()

    0
