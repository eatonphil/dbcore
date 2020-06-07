import React from 'react';

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
    <form onSubmit={handleSubmit}>
      <label>
        Username:
        <input value={username} onChange={(e) => setUsername(e.target.value)} />
      </label>
      <label>
        Password:
        <input value={password} onChange={(e) => setPassword(e.target.value)} type="password" />
      </label>
      {error && <div className="text-red-600">{error}</div>}
    </form>
  );
}
