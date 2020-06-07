import React from 'react';

export function useListData(endpoint: string) {
  const [cols, setCols] = React.useState([]);
  const [rows, setRows] = React.useState([]);
  const [error, setError] = React.useState('');
  const [offset, setOffset] = React.useState(0);
  const [limit, setLimit] = React.useState(25);
  const [filter, setFilter] = React.useState('');
  React.useEffect(function () {
    async function fetchRows() {
      setError('');

      const req = await window.fetch(`http://localhost:9090${endpoint}?offset=${offset}&limit=${limit}&filter=${filter}`);
      const rsp = await req.json();
      if (rsp.error) {
        setError(error);
        return;
      }

      setRows(rsp.data);
      setCols(Object.keys(rsp.data[0]));
    }

    fetchRows();
  }, [offset, limit, filter]);

  return {
    cols,
    rows,
    error,
    offset,
    setOffset,
    limit,
    setLimit,
    filter,
    setFilter,
  };
}
