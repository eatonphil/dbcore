import React from 'react';
import { Link } from 'react-router-dom';

export function Header() {
    return (
	<header style={{ display: 'flex' }}>
	    <Link to="/"><h1>Todo</h1></Link>
	    {/* TODO: define this */}
	    <Link>Logout</Link>
	</header>
    );
}
