module Database


type Column =
    {
        Name: string
        Type: string
        Nullable: bool
        AutoIncrement: bool
    }


type Constraint =
    {
        Column: string
        Type: string
        ForeignTable: string
        ForeignColumn: string
    }


type Table =
    {
        Name: string
        Columns: Column[]
        ForeignKeys: Constraint[]
        PrimaryKey: Option<Constraint>
    }
