import React from 'react';
{{~ if table.primary_key.value ~}}
import { useHistory } from 'react-router-dom';
{{~ end ~}}

import { {{ table.label|dbcore_capitalize }} } from '../api';
import { Heading } from '../components/Heading';
import { Link } from '../components/Link';
import { List } from '../components/List';
import { useListData } from '../hooks/useListData';

export function {{ table.label|dbcore_capitalize }}List() {
  {{~ if table.primary_key.value ~}}
  const history = useHistory();
  {{~ end ~}}
  const actions = (
    <Link to="/{{ table.label }}/create">Create</Link>
  );

  const data = useListData<{{ table.label|dbcore_capitalize }}>("{{ table.label }}");

  return (
    <>
      <Heading
        size="xl"
        actions={actions}
      >{{ table.label|dbcore_capitalize }}</Heading>
      <List
        data={data}
        {{~ if table.primary_key.value ~}}
        onRowClick={(row: {{ table.label|dbcore_capitalize }}) =>
          history.push("/{{ table.label }}/_/"+row["{{ table.primary_key.value.column }}"])}
        {{~ end ~}}
      />
    </>
  );
}
