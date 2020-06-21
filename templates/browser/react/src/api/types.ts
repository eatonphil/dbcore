{{~
  func toTypeScriptType
    case $0
      when "integer", "int", "bigint", "smallint", "decimal", "numeric", "real", "double precision"
        "number"
      when "boolean"
        "boolean"
      else
        "string"
    end
  end
~}}
{{~ for table in tables ~}}
export interface {{ table.name|dbcore_capitalize }} {
  {{~ for column in table.columns ~}}
  {{ column.name }}{{ if column.nullable }}?{{ end }}: {{ toTypeScriptType column.type }};
  {{~ end ~}}
}
{{~ end ~}}
