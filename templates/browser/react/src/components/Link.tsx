import React from 'react';
import { Link as RRLink } from 'react-router-dom';

interface Props {
  children: React.ReactNode;
  className: string;
  to: string;
}

export function Link({
  children,
  className,
  to,
}: Props) {
  return (
    <RRLink
      to={to}
      className={`text-blue-500 hover:text-blue-800 ${className}`}
      children={children}
    />
  );
}
