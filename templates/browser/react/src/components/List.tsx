import React from 'react';

import { useListData } from '../hooks/useListData';

interface Props {
  endpoint: string;
}

export function List(props: Props) {
  const { cols, error, filter, rows, setFilter } = useListData(props.endpoint);

  return (
    <div>
      {/* TODO: debounce this */}
      <textarea
        onChange={(e) => setFilter(e.target.value)}
        value={filter}
        className="border rounded w-full py-2 px-3 text-gray-700 mb-4 leading-tight"
        placeholder="id = 2"
      />

      {error && <div className="text-red-600">{error}</div>}

      <div className={`grid grid-cols-${cols.length} border-l border-r border-t`}>
        {cols.map(col => <div className="font-semibold text-gray-700 p-3 bg-gray-100 border-b">{col}</div>)}
        {rows.map(row => cols.map(col => <div className="text-xs text-gray-700 p-3 border-b">{row[col]}</div>))}
      </div>
    </div>
  );
}
