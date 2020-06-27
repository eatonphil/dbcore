import React from 'react';
import { useHistory, useRouteMatch } from 'react-router-dom';

import { Form } from '../components/Form';
import { Heading } from '../components/Heading';
import { Input } from '../components/Input';
import { Link } from '../components/Link';
import { request } from '../api';

{{~
  func javascriptValueify
    case $0
      when "integer", "int", "bigint", "smallint", "decimal", "numeric", "real", "double precision"
        "Number"
      when "boolean"
        "Boolean"
      else
        ""
    end
  end
~}}

export function {{ table.name|dbcore_capitalize }}Update() {
  const [state, setState] = React.useState<{ [key: string]: string }>({
    {{~ for column in table.columns ~}}
    {{~ if column.auto_increment
          continue
        end ~}}
    '{{ column.name }}': '',
    {{~ end ~}}
  });

  const history = useHistory();
  const [error, setError] = React.useState('');
  const { params: { key } } = useRouteMatch();
  const handleSubmit = React.useCallback(async (e) => {
    e.preventDefault();
    setError('');

    try {
      const rsp = await request(`{{ table.name }}/${key}`, {
        {{~ for column in table.columns ~}}
        {{~ if column.auto_increment
              continue
            end ~}}
        '{{ column.name }}': {{ javascriptValueify column.type }}(state['{{ column.name }}']),
        {{~ end ~}}
      }, 'PUT');

      if (rsp.error) {
        setError(rsp.error);
        return false;
      }

      history.push(`/{{ table.name }}/_/${key}`);
    } finally {
      return false;
    }
  }, [
    key,
    {{~ for column in table.columns ~}}
    {{~ if column.auto_increment
          continue
        end ~}}
    state['{{ column.name }}'],
    {{~ end ~}}
  ]);

  const [loaded, setLoaded] = React.useState(false);
  React.useEffect(function() {
    async function fetchRow() {
      setError(error);
      const rsp = await request(`{{ table.name }}/${key}`);

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
      <Link to="/{{ table.name }}">{{ table.name|dbcore_capitalize }}</Link> / {key}
      <Heading size="xl">Update</Heading>
      <Form error={error} buttonText="Update" onSubmit={handleSubmit}>
        {{~ for column in table.columns ~}}
        {{~ if column.auto_increment
              continue
            end ~}}
        <div className="mb-4">
          <Input
            label="{{ column.name }}"
            id="{{ column.name }}"
            value={state['{{ column.name }}']}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
              // e.target.value is not available within the setState callback, so copy it.
              // https://duncanleung.com/fixing-react-warning-synthetic-events-in-setstate/
              const { value } = e.target;
              setState((s: { [key: string]: string }) => ({ ...s, ['{{ column.name }}']: value }))
            }}
          />
        </div>
        {{ end }}
      </Form>
    </>
  );
}
