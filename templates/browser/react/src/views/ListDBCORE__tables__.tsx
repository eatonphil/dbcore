import React from 'react';

export function {{ table|string.capitalize }}() {
  const { cols, rows } = useListData(`/v1/{{table}}`);

  return (
    <>
      <h2>{{ table|string.capitalize }}</h2>
      <div class={`grid grid-cols-${cols.length}`}>
        {cols.map(col => <div class="grid-header-col">{col}</div>)}
        {rows}
      </div>
    </>
  );
}
