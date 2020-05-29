module Config

open System.Collections.Generic
open System.IO

open YamlDotNet.Serialization
open YamlDotNet.Serialization.NamingConventions


type DatabaseConfig() =
    member val Dialect = "postgres" with get, set
    member val Host = "localhost" with get, set
    member val Port = "5432" with get, set
    member val Database = "" with get, set
    member val Username = "" with get, set
    member val Password = "" with get, set


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

    member val Auth = new ApiAuthConfig() with get, set
    member val RouterPrefix = "" with get, set
    member val Extra = new Dictionary<string, string>() with get, set


type BrowserConfig() =
    interface IConfig with
        member val OutDir = "browser" with get, set
        member val Template = "react-ts" with get, set


type CustomConfig() =
    interface IConfig with
        member val OutDir = "" with get, set
        member val Template = "" with get, set

    // TODO: allow any nested structure, unsure how to type
    member val Extra = new Dictionary<string, string>() with get, set

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
        (new DeserializerBuilder())
            .WithNamingConvention(CamelCaseNamingConvention.Instance)
            .Build()
    let config = deserializer.Deserialize<Config>(stream)

    // TODO: validate config
    config
