import React from 'react';
{{~ if table.primary_key.value ~}}
import { useHistory } from 'react-router-dom';
{{~ end ~}}

import { Heading } from '../components/Heading';
import { Link } from '../components/Link';
import { List } from '../components/List';

export function {{ table.name|string.capitalize }}List() {
  {{~ if table.primary_key.value ~}}
  const history = useHistory();
  {{~ end ~}}
  const actions = (
    <Link to="/{{ table.name }}/create">Create</Link>
  );

  return (
    <>
      <Heading
        size="xl"
        actions={actions}
      >{{ table.name|string.capitalize }}</Heading>
      <List
        {{~ if table.primary_key.value ~}}
        onRowClick={(row) => history.push("/{{ table.name }}/_/"+row["{{ table.primary_key.value.column }}"])}
        {{~ end ~}}
        endpoint="{{ table.name }}"
      />
    </>
  );
}
