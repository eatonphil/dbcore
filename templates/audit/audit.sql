{{~ for table in tables ~}}
{{~ name = api.audit.table_prefix + table.name + api.audit.table_suffix ~}}
{{~ if !table.primary_key.value
      continue
    end
~}}
CREATE TABLE "{{ name }}" (
    kind TEXT NOT NULL,
    body JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    id {{ table.primary_key.value.type }},
    FOREIGN KEY (id) REFERENCES "{{ table.name }}" ("{{ table.primary_key.value.name }}");
);
CREATE INDEX "{{ name }}_created_at_idx" ON "{{ name }}" ("id");

{{ end }}
