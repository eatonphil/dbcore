open System

open FSharp.Configuration
open Npgsql

open Database

// See: https://stackoverflow.com/a/30903481/1507139
type Config =
    {
        
    }

let getConfig () : Config =

[<EntryPoint>]
let main args =
    let config = getConfig()

    let db = new DatabaseReader(
        config.Database.Host,
        config.Database.Database,
        config.Database.Username,
        config.Database.Password)
    let tables = db.getTables()
    0
