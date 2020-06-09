import React from 'react';

import { useListData } from '../hooks/useListData';

interface Props {
  endpoint: string;
}

function PageButton({
  disabled,
  onClick,
  children,
}: HTMLButtonElement) {
  return (
    <button
      className={'p-1 ' + (disabled ? 'text-gray-300' : '')}
      title={disabled ? 'Disabled' : null}
      disabled={disabled}
      onClick={onClick}
      children={children}
    />
  );
}

export function List(props: Props) {
  const { cols, error, filter, offset, limit, rows, setFilter, setOffset, total } = useListData(props.endpoint);

  const pageInfo = (
    <div className="text-right text-gray-700 p-2">
      <PageButton
        disabled={offset === 0}
        onClick={() => setOffset(o => o - limit)}
      >&#x2BC7;</PageButton>
      Showing {offset + 1}-{offset + Math.min(limit, rows.length)} of {total}
      <PageButton
        disabled={Math.floor(offset / limit) === Math.floor(total / limit)}
        onClick={() => setOffset(o => o + limit)}
      >&#x2BC8;</PageButton>
    </div>
  );

  return (
    <div>
      <textarea
        onChange={(e) => setFilter(e.target.value)}
        value={filter}
        className="border w-full py-2 px-3 text-gray-700 mb-4 leading-tight"
        placeholder="id = 2"
      />

      {error && <div className="bg-red-600 text-white p-2 mb-2">{error}</div>}

      {total === 0 ? 'No results found.' : (
        <>
          {pageInfo}
          <div className={`grid grid-cols-${cols.length} border-l border-r border-t`}>
            {cols.map(col => <div key={`header-${col}`} className="font-semibold text-gray-700 p-3 bg-gray-100 border-b">{col}</div>)}
            {rows.map((row, i) => cols.map(col => <div key={`cell-${i}-${col}`} className="text-gray-700 p-3 border-b">{row[col]}</div>))}
          </div>
          {pageInfo}
        </>
      )}
    </div>
  );
}
