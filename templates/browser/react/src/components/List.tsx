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
        className="w-full"
      />

      {error && <div className="text-red-600">{error}</div>}

      <div className={`grid grid-cols-${cols.length}`}>
        {cols.map(col => <div className="grid-header-col">{col}</div>)}
        {rows}
      </div>
    </div>
  );
}
