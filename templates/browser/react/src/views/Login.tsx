import React from 'react';

import { Button } from '../components/Button';
import { Form } from '../components/Form';
import { Input } from '../components/Input';

export function Login() {
  const [username, setUsername] = React.useState('');
  const [password, setPassword] = React.useState('');

  const [error, setError] = React.useState('');
  const handleSubmit = React.useCallback(async (e) => {
    e.preventDefault();
    setError('');

    try {
      const req = await window.fetch('http://localhost:9090/v1/session/start', {
        method: 'POST',
        body: JSON.stringify({
          username,
          password,
        }),
        headers: {
          'content-type': 'application/json',
        },
        credentials: 'include',
      });

      const rsp = await req.json();
      if (rsp.error) {
        setError(rsp.error);
        return;
      }

      window.location.href = '/';
    } catch (e) {
      console.error(e);
      return false;
    }
  });

  return (
    <div className="flex justify-center">
      <div className="w-full max-w-xs">
        <Form onSubmit={handleSubmit}>
          <div className="mb-4">
            <Input
              label="Username"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </div>
          <div className="mb-4">
            <Input
              label="Password"
              id="password"
              value={password}
              type="password"
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>
          <Button type="submit">Sign in</Button>
          {error && <div className="text-red-600 text-sm mt-4">{error}</div>}
        </Form>
      </div>
    </div>
  );
}
