import React from 'react';
import { Link } from 'react-router-dom';

export function Home() {
  return (
    <>
      <h2>Home!</h2>
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
