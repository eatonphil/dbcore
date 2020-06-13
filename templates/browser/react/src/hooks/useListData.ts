import React from 'react';

import { request } from '../api';

export function useListData(endpoint: string) {
  const [cols, setCols] = React.useState([]);
  const [rows, setRows] = React.useState([]);
  const [total, setTotal] = React.useState(0);
  const [error, setError] = React.useState('');
  const [offset, setOffset] = React.useState(0);
  const [limit, setLimit] = React.useState(25);
  const [filter, setFilter] = React.useState('');
  const [sortColumn, setSortColumn] = React.useState('id');
  const [sortOrder, setSortOrder] = React.useState('asc');

  React.useEffect(function () {
    async function fetchRows() {
      setError('');

      const params = Object.entries({ offset, limit, filter, sortColumn, sortOrder })
                           .map(([key, value]) => `${key}=${value}`)
                           .join('&');
      const rsp = await request(endpoint + '?' + params);
      if (rsp.error) {
        setRows([]);
        setCols([]);
        setTotal(0);
        setError(rsp.error);
        return;
      }

      setTotal(rsp.total);
      setRows(rsp.data);
      setCols(rsp.data.length ? Object.keys(rsp.data[0]) : []);
    }

    fetchRows();
  }, [offset, limit, filter, sortColumn, sortOrder]);

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
    total,
  };
}
