import * as Router from 'preact-router';
import * as React from 'preact';
import { TextField } from "@fluentui/react/lib/TextField";


function Login() {
  const [username, setUsername] = React.useState("");
  const [password, setPassword] = React.useState("");

  return (
    <>
      <form>
        <TextField
          label="Username"
          value={username}
          onChange={(_, value) => setUsername(value)}
        />
        <TextField
          label="Password"
          type="password"
          value={password}
          onChange={(_, value) => setPassword(value)}
        />
      </form>
    </>
  )
}


function App() {
  return (
    <></>
  );
};

React.render(<App />, document.getElementById('root'));
