import React from 'react';

import { Button } from '../components/Button';
import { Form } from '../components/Form';
import { Input } from '../components/Input';

export function Login() {
  const [username, setUsername] = React.useState('');
  const [password, setPassword] = React.useState('');

  const [error, setError] = React.useState('');
  const handleSubmit = React.useCallback(async () => {
    setError('');

    const req = await window.fetch('http://localhost:9091/v1/session/create', {
      method: 'POST',
      body: JSON.stringify({
        username,
        password,
      }),
      headers: {
        'content-type': 'application/json',
      },
    });

    const rsp = await req.json();
    if (rsp.error) {
      setError(rsp.error);
      return;
    }

    window.location.href = '/';
  });

  return (
    <div className="w-full max-w-xs">
      <Form onSubmit={handleSubmit}>
        <Input
          <div className="mb-4">
          label="Username"
          id="username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          />
        </div>
        <div className="mb-6">
          <Input
            label="Password"
            id="password"
            value={password}
            type="password"
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        {error && <div className="text-red-600">{error}</div>}
        <Button>Sign in</Button>
      </Form>
    </div>
  );
}
