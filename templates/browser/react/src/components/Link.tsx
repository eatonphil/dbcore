import React from 'react';
import { Link as RRLink } from 'react-router-dom';

interface Props {
  to: string;
  children: React.ReactNode;
}

export function Link({
  to,
  children,
}: Props) {
  return (
    <RRLink
      to={to}
      className="text-blue-500 hover:text-blue-800"
      children={children}
    />
  );
}
