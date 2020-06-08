import React from 'react';
import { Link } from 'react-router-dom';

import { Heading } from '../components/Heading';

export function Home() {
  return (
    <>
      <Heading size="xl">Home!</Heading>
      <ul>
        {{~ for table in tables ~}}
        {{~ if table.primary_key.is_none
              continue
            end
        ~}}
        <li>
          <Link to="/{{ table.name }}">{{ table.name|string.capitalize }}</Link>
        </li>
        {{~ end ~}}
      </ul>
    </>
  );
}
