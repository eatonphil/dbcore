import React from 'react';

import { List } from '../components/List';

export function {{ table|string.capitalize }}() {
  return (
    <>
      <h2>{{ table|string.capitalize }}</h2>
      <List
        endpoint="/v1/{{ table }}"
      />
    </>
  );
}
