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
{{~ if table.primary_key.is_none
      continue
    end
~}}
import { {{ table.name | string.capitalize }} } from './views/List{{ table.name }}';
{{~ end ~}}

function App() {
  const [pageLoaded, setPageLoaded] = React.useState(false);
  React.useEffect(() => {
    const match = document.cookie.match(new RegExp('(^| )au=([^;]+)'));
    const sessionToken = match ? match[2] : '';
    if (!sessionToken && window.location.pathname !== '/login') {
      window.location.href = '/login';
      return;
    }

    setPageLoaded(true);
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
          {{~ if table.primary_key.is_none
                continue
              end
          ~}}
          <Route exact path="/{{ table.name }}">
            <{{ table.name | string.capitalize}} />
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
