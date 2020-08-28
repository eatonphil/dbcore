module Config

open System.Collections.Generic
open System.IO

open YamlDotNet.Serialization
open YamlDotNet.Serialization.NamingConventions


type DatabaseTableConfig() =
    let mutable label = ""

    member val Name = "" with get, set

    member this.Label
        with get() : string =
            if label <> "" then label else this.Name
        and set(value: string) = label <- value


type DatabaseConfig() =
    let mutable port = ""
    let mutable schema = ""

    member val Dialect = "postgres" with get, set
    member val Host = "localhost" with get, set
    member val Database = "" with get, set
    member val Username = "" with get, set
    member val Password = "" with get, set

    member val Tables: array<DatabaseTableConfig> = [||] with get, set

    member this.Schema
        with get() : string =
            if schema <> "" then schema else
                match this.Dialect with
                    | "postgres" -> "public"
                    | "mysql" -> this.Database
                    | "sqlite" -> ""
                    | _ -> failwith ("database.schema must be set for unknown dialect: " + this.Dialect)

    member this.Port
        with get() : string =
            if port <> "" then port else
                match this.Dialect with
                    | "postgres" -> "5432"
                    | "mysql" -> "3306"
                    | "sqlite" -> ""
                    | _ -> failwith ("database.port must be set for unknown dialect: " + this.Dialect)
        and set(value: string) = port <- value


type IConfig =
    abstract member OutDir : string with get, set
    abstract member Template : string with get, set


type ApiAuthAllowConfig() =
    member val Get = "" with get, set
    member val Put = "" with get, set
    member val Post = "" with get, set
    member val Delete = "" with get, set

type ApiAuthConfig() =
    member val Enabled = false with get, set
    member val Table = "users" with get, set
    member val Username = "username" with get, set
    member val Password = "password" with get, set
    member val Allow = Dictionary<string, ApiAuthAllowConfig>() with get, set


type ApiAuditConfig() =
    member val Enabled = false with get, set
    member val CreatedAt = "" with get, set
    member val UpdatedAt = "" with get, set
    member val DeletedAt = "" with get, set


type ApiConfig() =
    interface IConfig with
        member val OutDir = "api" with get, set
        member val Template = "go" with get, set

    member val Auth = ApiAuthConfig() with get, set
    member val RouterPrefix = "" with get, set
    member val Extra = Dictionary<string, obj>() with get, set
    member val Runtime = Dictionary<string, obj>() with get, set

    member val Audit = ApiAuditConfig() with get, set

    member this.Validate() =
        // TODO: should probably turn validation into a schema per template
        if (this :> IConfig).Template = "go" && not (this.Extra.ContainsKey "repo")
            then failwith "Repo is required: `api.extra.repo: $repo`"


type BrowserConfig() =
    interface IConfig with
        member val OutDir = "browser" with get, set
        member val Template = "react" with get, set

    member val DefaultRoute = "" with get, set


type CustomConfig() =
    interface IConfig with
        member val OutDir = "" with get, set
        member val Template = "" with get, set

    member val Extra = Dictionary<string, obj>() with get, set

type Config() =
    member val Project = "" with get, set
    member val CultureName = System.Globalization.CultureInfo.CurrentCulture.Name with get, set

    member val Database = DatabaseConfig() with get, set
    member val Api = ApiConfig() with get, set
    member val Browser = BrowserConfig() with get, set
    member val Custom : array<CustomConfig> = [| |] with get, set

    member this.Validate() =
        if this.Project = "" then failwith "Project name is required: `project: $name`"

        this.Api.Validate()

let GetConfig(f: string) : Config =
    use file = new FileStream(f, FileMode.Open, FileAccess.Read)
    use stream = new StreamReader(file)
    let deserializer =
        (DeserializerBuilder())
            .WithNamingConvention(CamelCaseNamingConvention.Instance)
            .Build()
    let config = deserializer.Deserialize<Config>(stream)

    config
