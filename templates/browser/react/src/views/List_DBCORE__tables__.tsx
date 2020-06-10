import React from 'react';

import { Heading } from '../components/Heading';
import { Link } from '../components/Link';
import { List } from '../components/List';

export function List{{ table.name|string.capitalize }}() {
  const actions = (
    <Link to="/{{ table.name }}/create">Create</Link>
  );

  return (
    <>
      <Heading size="xl" actions={actions}>{{ table.name|string.capitalize }}</Heading>
      <List endpoint="{{ table.name }}" />
    </>
  );
}
