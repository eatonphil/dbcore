import React from 'react';

import { Heading } from '../components/Heading';
import { Link } from '../components/Link';

export function Home() {
  return (
    <>
      <Heading size="xl">Home</Heading>
      <ul>
        {{~ for table in tables ~}}
        {{~ if !table.primary_key.value
              continue
            end
        ~}}
        <li>
          <Link to="/{{ table.label }}">{{ table.label|dbcore_capitalize }}</Link>
        </li>
        {{~ end ~}}
      </ul>
    </>
  );
}
