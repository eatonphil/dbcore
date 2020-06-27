import React from 'react';
import { useHistory } from 'react-router-dom';

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

export function {{ table.label|dbcore_capitalize }}Create() {
  const [state, setState] = React.useState({
    {{~ for column in table.columns ~}}
    {{~ if column.auto_increment
          continue
        end ~}}
    '{{ column.name }}': '',
    {{~ end ~}}
  });

  const history = useHistory();
  const [error, setError] = React.useState('');
  const handleSubmit = React.useCallback(async (e) => {
    e.preventDefault();
    setError('');

    try {
      const rsp = await request('{{ table.label }}', {
        {{~ for column in table.columns ~}}
        {{~ if column.auto_increment || (api.audit.enabled && ([api.audit.created_at, api.audit.updated_at, api.audit.deleted_at] | array.contains column.name))
              continue
            end ~}}
        '{{ column.name }}': {{ javascriptValueify column.type }}(state['{{ column.name }}']),
        {{~ end ~}}
      });

      if (rsp.error) {
        setError(rsp.error);
        return false;
      }

      history.push('/{{ table.label }}');
    } finally {
      return false;
    }
  }, [
    history,
    {{~ for column in table.columns ~}}
    {{~ if column.auto_increment
          continue
        end ~}}
    state['{{ column.name }}'],
    {{~ end ~}}
  ]);

  return (
    <>
      <Link to="/{{ table.label }}">{{ table.label|dbcore_capitalize }}</Link>
      <Heading size="xl">Create</Heading>
      <Form error={error} buttonText="Create" onSubmit={handleSubmit}>
        {{~ for column in table.columns ~}}
        {{~ if column.auto_increment || (api.audit.enabled && ([api.audit.created_at, api.audit.updated_at, api.audit.deleted_at] | array.contains column.name))
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
              setState(s => ({ ...s, ['{{ column.name }}']: value }))
            }}
          />
        </div>
        {{ end }}
      </Form>
    </>
  );
}
