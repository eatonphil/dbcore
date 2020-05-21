open System

open Npgsql

[<EntryPoint>]
let main args =
    let dsn = sprintf "Host=%s;Database=%s;Username=%s;Password=%s;" (args.[0]) (args.[1]) (args.[2]) (args.[3])
    let conn = new NpgsqlConnection(dsn);
    conn.Open();
    let cmd = new NpgsqlCommand("SELECT table_name FROM information_schema.tables WHERE table_schema='public'", conn);
    let dr = cmd.ExecuteReader();
    while dr.Read() do
        printfn "%s" (string (dr.[0]));

    0
