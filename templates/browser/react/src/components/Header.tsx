import React from 'react';

import { Link } from '../components/Link';

export function Header() {
  return (
    <div className="border-b p-4 mb-4">
      <nav className="container mx-auto flex">
        <div className="text-sm lg:flex-grow">
          <h1>
            <Link to="/">
              {{ project|string.capitalize }}
            </Link>
          </h1>
        </div>
        {window.location.pathname !== "/login" && (
          <div className="flex">
            <Link to="#">Logout</Link>
          </div>
        )}
      </nav>
    </div>
  );
}
