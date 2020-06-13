import React from 'react';

import { Form } from '../components/Form';
import { Input } from '../components/Input';
import { request } from '../api';

export function Logout() {
  React.useEffect(() => {
    async function stop() {
      // Should plan to redirect to /login even if the request fails
      // for some reason, but does need to be await-ed so the request
      // is not cancelled on navigation.
      try {
        await request('session/stop', {
          method: 'POST',
          credentials: 'include',
        });
      } finally {
        window.location.href = '/login';
      }
    }

    stop();
  });

  return null;
}

export function Login() {
  const [username, setUsername] = React.useState('');
  const [password, setPassword] = React.useState('');

  const [error, setError] = React.useState('');
  const handleSubmit = React.useCallback(async (e) => {
    e.preventDefault();
    setError('');

    try {
      const rsp = await request('session/start', {
        username,
        password,
      });

      if (rsp.error) {
        setError(rsp.error);
        return;
      }

      const params = new URLSearchParams(window.location.search);
      window.location.href = params.get('return') || '/';
    } catch (e) {
      // Need the try-catch so we can return false here.
      console.error(e);
      return false;
    }
  });

  return (
    <div className="flex justify-center">
      <div className="w-full max-w-xs">
        <Form buttonText="Sign In" error={error} onSubmit={handleSubmit}>
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
        </Form>
      </div>
    </div>
  );
}
