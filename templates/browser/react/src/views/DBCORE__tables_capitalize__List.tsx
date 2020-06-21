import React from 'react';
{{~ if table.primary_key.value ~}}
import { useHistory } from 'react-router-dom';
{{~ end ~}}

import { {{ table.name|dbcore_capitalize }} } from '../api';
import { Heading } from '../components/Heading';
import { Link } from '../components/Link';
import { List } from '../components/List';
import { useListData } from '../hooks/useListData';

export function {{ table.name|dbcore_capitalize }}List() {
  {{~ if table.primary_key.value ~}}
  const history = useHistory();
  {{~ end ~}}
  const actions = (
    <Link to="/{{ table.name }}/create">Create</Link>
  );

  const data = useListData<{{ table.name|dbcore_capitalize }}>("{{ table.name }}");

  return (
    <>
      <Heading
        size="xl"
        actions={actions}
      >{{ table.name|dbcore_capitalize }}</Heading>
      <List
        data={data}
        {{~ if table.primary_key.value ~}}
        onRowClick={(row: {{ table.name|dbcore_capitalize }}) =>
          history.push("/{{ table.name }}/_/"+row["{{ table.primary_key.value.column }}"])}
        {{~ end ~}}
      />
    </>
  );
}
