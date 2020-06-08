import React from 'react';
import { Link } from 'react-router-dom';

export function Home() {
  return (
    <>
      <h2>Home!</h2>
      {{~ for table in tables ~}}
      <Link to="/{{ table.name }}">{{ table.name|string.capitalize }}</Link>
      {{~ end ~}}
    </>
  );
}
