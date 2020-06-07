import React from 'react';
import { Link } from 'react-router-dom';

export function Header() {
  return (
    <nav className="flex items-center justify-between flex-wrap p-4 mb-4 border-b">
      <div className="text-sm lg:flex-grow">
        <h1>
          <Link
            to="/"
            className="text-blue-500 hover:text-blue-800"
          >
            {{ project|string.capitalize }}
          </Link>
        </h1>
      </div>
      <div>
        {window.location.pathname !== "/login" &&
          <Link to="#" className="text-blue-500 hover:text-blue-800">Logout</Link>}
      </div>
    </nav>
  );
}
