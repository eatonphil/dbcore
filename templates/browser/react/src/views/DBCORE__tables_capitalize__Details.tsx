import React from 'react';
import { useRouteMatch } from 'react-router-dom';

import { Heading } from '../components/Heading';
import { Link } from '../components/Link';
import { request } from '../api';

export function {{ table.label|dbcore_capitalize }}Details() {
  const { params: { key } } = useRouteMatch();
  const actions = (
    <Link to={`/{{ table.label }}/_/${key}/update`}>Update</Link>
  );

  const [loaded, setLoaded] = React.useState(false);
  const [state, setState] = React.useState<{ [key: string]: string }>({
    {{~ for column in table.columns ~}}
    {{~ if column.auto_increment
          continue
        end ~}}
    '{{ column.name }}': '',
    {{~ end ~}}
  });

  const [error, setError] = React.useState('');
  React.useEffect(function() {
    async function fetchRow() {
      setError(error);
      const rsp = await request(`{{ table.label }}/${key}`);

      if (rsp.error) {
        setError(rsp.error);
        return;
      }

      setState(rsp);
      setLoaded(true);
    }

    fetchRow();
  }, [key]);

  if (!loaded) {
    return null;
  }

  return (
    <>
      <Link to="/{{ table.label }}">{{ table.label|dbcore_capitalize }}</Link>
      <Heading size="xl" actions={actions}>{key}</Heading>
      {{~ for column in table.columns ~}}
      <div className="mb-4 border-b">
        <div className="block text-gray-700 text-sm font-bold mb-2 uppercase text-sm">
          {{ column.name }}
        </div>
        <div>{state["{{ column.name }}"]}</div>
      </div>
      {{ end }}
    </>
  );
}
