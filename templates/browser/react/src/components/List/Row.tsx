import React from 'react';

interface Props<T> {
  children: React.ReactNode;
  header: (keyof T)[];
  key: string;
  onClick?: () => void;
  row: T;
}

export function Row<T>({
  children,
  header,
  key,
  onClick,
  row,
}: Props<T>) {
  const className = "text-gray-700 p-3 border-b" +
    (onClick ? " cursor-pointer" : "");
  return header.map(col => (
    <div
      key={`${key}-${col}`}
      className={className}
      onClick={onClick ? onClick : () => { }}
    >{row[col]}</div>
  ));
}
