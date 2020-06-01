module Database


type Column =
    {
        Name: string
        Type: string
        GoType: string
        AutoIncrement: bool
    }


type Constraint =
    {
        Column: string
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
