import React from 'react';

import { Header } from './Header';
import { Row } from './Row';
import { useListData } from '../../hooks/useListData';

interface Props<T> {
  onRowClick?: (row: T) => void;
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

export function List<T>({
  endpoint,
  onRowClick,
}: Props<T>) {
  const {
    cols,
    error,
    filter,
    offset,
    limit,
    rows,
    setFilter,
    setOffset,
    total,
  } = useListData(endpoint);

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

  const handleFilter = React.useCallback((e) => {
    setFilter(e.target.value);
    // Also reset page back to zero on filter change.
    setOffset(0);
  });

  return (
    <div>
      <textarea
        id="filter"
        onChange={handleFilter}
        value={filter}
        className="border w-full py-2 px-3 text-gray-700 mb-4 leading-tight"
        placeholder="Type to filter (e.g. id = 2)"
      />

      {error && <div className="bg-red-600 text-white p-2 mb-2">{error}</div>}

      {total === 0 ? 'No results found.' : (
        <>
          {pageInfo}
          <div className={`grid grid-cols-${cols.length} border-l border-r border-t`}>
            {cols.map(col => (
              <Header key={`header-${col}`} column={col} />
            ))}
            {rows.map((row, i) => (
              <Row
                onClick={onRowClick ? () => onRowClick(rows[i]) : null}
                key={`row-${i}`}
                row={row}
                header={cols}
              />
            ))}
          </div>
          {pageInfo}
        </>
      )}
    </div>
  );
}
