import React from 'react';
import ReactDOM from 'react-dom';
import {
  BrowserRouter as Router,
  Switch,
  Route
} from 'react-router-dom';

import { Header } from './components/Header';
import { Home } from './views/Home';
{{ if api.auth.enabled }}
import { Login, Logout } from './views/Login';
{{ end }}
{{~ for table in tables ~}}
import { {{ table.name | dbcore_capitalize }}List } from './views/{{ table.name|dbcore_capitalize }}List';
{{~ if !table.primary_key.value
      continue
    end
~}}
import { {{ table.name | dbcore_capitalize }}Create } from './views/{{ table.name|dbcore_capitalize }}Create';
import { {{ table.name | dbcore_capitalize }}Update } from './views/{{ table.name|dbcore_capitalize }}Update';
import { {{ table.name | dbcore_capitalize }}Details } from './views/{{ table.name|dbcore_capitalize }}Details';
{{~ end ~}}

function App() {
  {{ if api.auth.enabled }}
  const [pageLoaded, setPageLoaded] = React.useState(false);
  React.useEffect(() => {
    const match = document.cookie.match(new RegExp('(^| )au=([^;]+)'));
    const sessionToken = match ? match[2] : '';
    const isLogin = window.location.pathname === '/login';
    if (!sessionToken && !isLogin) {
      window.location.href = '/login?return=' + encodeURIComponent(window.location.href);
      return;
    } else if (sessionToken && isLogin) {
      window.location.href = '/';
      return;
    }

    setPageLoaded(true);
  }, []);

  if (!pageLoaded) {
    return null;
  }
  {{ end }}

  return (
    <Router>
      <div>
        <Header />

        <div className="container mx-auto">
          <Switch>
            <Route exact path="/">
              <Home />
            </Route>
            {{ if api.auth.enabled }}
            <Route exact path="/login">
              <Login />
            </Route>
            <Route exact path="/logout">
              <Logout />
            </Route>
            {{ end }}
          
            {{~ for table in tables ~}}
            <Route exact path="/{{ table.name }}">
              <{{ table.name | dbcore_capitalize}}List />
            </Route>
            {{~ if !table.primary_key.value
                  continue
                end
            ~}}
            <Route exact path="/{{ table.name }}/create">
              <{{ table.name | dbcore_capitalize}}Create />
            </Route>
            <Route exact path="/{{ table.name }}/_/:key/update">
              <{{ table.name | dbcore_capitalize}}Update />
            </Route>
            <Route exact path="/{{ table.name }}/_/:key">
              <{{ table.name | dbcore_capitalize}}Details />
            </Route>
            {{ end }}
          </Switch>
        </div>
      </div>
    </Router>
  );
}

window.onload = function () {
  ReactDOM.render(<App />, document.getElementById('root'));
}
