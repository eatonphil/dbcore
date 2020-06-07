module Config

open System.Collections.Generic
open System.IO

open YamlDotNet.Serialization
open YamlDotNet.Serialization.NamingConventions


type DatabaseConfig() =
    let mutable port = ""
    let mutable schema = ""

    member val Dialect = "postgres" with get, set
    member val Host = "localhost" with get, set
    member val Database = "" with get, set
    member val Username = "" with get, set
    member val Password = "" with get, set

    member this.Schema
        with get() : string =
            if schema <> "" then schema else
                match this.Dialect with
                    | "postgres" -> "public"
                    | "mysql" -> this.Database
                    | _ -> failwith ("database.schema must be set for unknown dialect: " + this.Dialect)

    member this.Port
        with get() : string =
            if port <> "" then port else
                match this.Dialect with
                    | "postgres" -> "5432"
                    | "mysql" -> "3306"
                    | _ -> failwith ("database.port must be set for unknown dialect: " + this.Dialect)
        and set(value: string) = port <- value


type IConfig =
    abstract member OutDir : string with get, set
    abstract member Template : string with get, set


type ApiAuthConfig() =
    member val Enabled = false with get, set
    member val Table = "users" with get, set
    member val Username = "username" with get, set
    member val Password = "password" with get, set


type ApiConfig() =
    interface IConfig with
        member val OutDir = "api" with get, set
        member val Template = "go" with get, set

    member val Auth = ApiAuthConfig() with get, set
    member val RouterPrefix = "" with get, set
    member val Extra = Dictionary<string, object>() with get, set


type BrowserConfig() =
    interface IConfig with
        member val OutDir = "browser" with get, set
        member val Template = "react" with get, set


type CustomConfig() =
    interface IConfig with
        member val OutDir = "" with get, set
        member val Template = "" with get, set

    member val Extra = Dictionary<string, object>() with get, set

type Config() =
    member val Project = "" with get, set

    member val Database = DatabaseConfig() with get, set
    member val Api = ApiConfig() with get, set
    member val Browser = BrowserConfig() with get, set
    member val Custom : array<CustomConfig> = [| |] with get, set


let GetConfig(f: string) : Config =
    use file = new FileStream(f, FileMode.Open, FileAccess.Read)
    use stream = new StreamReader(file)
    let deserializer =
        (DeserializerBuilder())
            .WithNamingConvention(CamelCaseNamingConvention.Instance)
            .Build()
    let config = deserializer.Deserialize<Config>(stream)

    // TODO: validate config
    config
