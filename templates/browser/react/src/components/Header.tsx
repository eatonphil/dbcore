import React from 'react';

import { Link } from './Link';

export function Header() {
  return (
    <div className="border-b p-4 mb-4">
      <nav className="container mx-auto flex">
        <div className="text-lg flex-grow">
          <h1>
            <Link to="/">
              {{ project|dbcore_capitalize }}
            </Link>
          </h1>
        </div>
        {{ if api.auth.enabled }}
        {window.location.pathname !== "/login" && (
          <div className="flex">
            <Link className="text-sm text-black-600" to="/logout">Logout</Link>
          </div>
        )}
        {{ end }}
      </nav>
    </div>
  );
}
