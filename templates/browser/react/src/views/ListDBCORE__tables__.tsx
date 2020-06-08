import React from 'react';

import { Heading } from '../components/Heading';
import { List } from '../components/List';

export function {{ table.name|string.capitalize }}() {
  return (
    <>
      <Heading size="xl">{{ table.name|string.capitalize }}</Heading>
      <List
        endpoint="/v1/{{ table.name }}"
      />
    </>
  );
}
