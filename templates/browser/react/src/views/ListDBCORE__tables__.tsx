import React from 'react';

import { List } from '../components/List';

export function {{ table.name|string.capitalize }}() {
  return (
    <>
      <h2>{{ table.name|string.capitalize }}</h2>
      <List
        endpoint="/v1/{{ table.name }}"
      />
    </>
  );
}
