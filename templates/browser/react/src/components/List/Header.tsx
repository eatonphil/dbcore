import React from 'react';

interface Props<T> {
  column: keyof T;
}

export function Header<T>({
  column,
}: Props<T>) {
  const className = "font-semibold text-gray-700 p-3 bg-gray-100 border-b uppercase text-sm";
  return <div className={className} children={column} />;
}
