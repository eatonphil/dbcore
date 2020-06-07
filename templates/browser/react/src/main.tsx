import React from 'react';
import ReactDOM from 'react-dom';
import {
  BrowserRouter as Router,
  Switch,
  Route
} from 'react-router-dom';

import { Header } from './components/Header';
import { Home } from './views/Home';
import { Login } from './views/Login';
{{~ for table in tables ~}}
import { {{ table | string.capitalize }} } from './views/List{{ table }}';
{{~ end ~}}

function App() {
  const [pageLoaded, setPageLoaded] = React.useState(false);
  React.useEffect(() => {
    async function fetchToken() {
      const sessionToken = await browser.cookie.get({ name: 'au' });
      if (!sessionToken) {
        window.location = Login.path;
        return;
      }

      setPageLoaded(true);
    }

    fetchToken();
  });

  if (!pageLoaded) {
    return null;
  }

  return (
    <Router>
      <div>
        <Header />

        <Switch>
          <Route exact path="/">
            <Home />
          </Route>
          <Route exact path="/login">
            <Login />
          </Route>

          {{~ for table in tables ~}}
          <Route exact path="/{{ table }}">
            <{{ table | string.capitalize}} />
          </Route>
          {{~ end ~}}
        </Switch>
      </div>
    </Router>
  );
}

window.onload = function () {
  ReactDOM.render(<App />, document.getElementById('root'));
}
