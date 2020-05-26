module Config

open System.IO

open YamlDotNet.Serialization
open YamlDotNet.Serialization.NamingConventions


type DatabaseConfig() =
    // TODO: only postgres is supported
    member val Dialect = "postgres" with get, set
    member val Host = "localhost" with get, set
    member val Port = "5432" with get, set
    member val Database = "" with get, set
    member val Username = "" with get, set
    member val Password = "" with get, set
    member val Parameters = "" with get, set


type ApiConfig() =
    let mutable outdir = ""
    member this.OutDir
        with get() = if outdir = "" then this.Language else ""
        and set(value) = outdir <- value

    member val Language = "go" with get, set
    member val Repo = "" with get, set
    member val Project = "" with get, set
    member val Address = "" with get, set


type Config() =
    member val Database = DatabaseConfig() with get, set
    member val Api = ApiConfig() with get, set


let GetConfig(f: string) : Config =
    let file = new FileStream(f, FileMode.Open, FileAccess.Read)
    let stream = new StreamReader(file)
    let deserializer = (new DeserializerBuilder()).WithNamingConvention(CamelCaseNamingConvention.Instance).Build()
    let config = deserializer.Deserialize<Config>(stream)
    stream.Close()

    // TODO: validate config
    config
