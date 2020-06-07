import React from 'react';

export function useListData(endpoint: string) {
  const [cols, setCols] = React.useState([]);
  const [rows, setRows] = React.useState([]);

  const [offset, setOffset] = React.useState(0);
  const [limit, setLimit] = React.useState(25);
  const [filter, setFilter] = React.useState('');
  React.useEffect(function () {
    async function fetchRows() {
      const req = await window.fetch(`http://localhost:9091${endpoint}?offset=${offset}&limit=${limit}&filter=${filter}`);
      const rsp = await req.json();
      setRows(rsp.data);
      setCols(Object.keys(rsp.data[0]));
    }

    fetchRows();
  }, [offset, limit, filter]);

  return {
    cols,
    rows,
    offset,
    setOffset,
    limit,
    setLimit,
    filter,
    setFilter,
  };
}
