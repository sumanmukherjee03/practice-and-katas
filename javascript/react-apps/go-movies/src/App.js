import React from 'react';
import { HashRouter as Router, Link, Route, Switch } from 'react-router-dom';

export default function App() {
  // Everything that you want to be routed to needs to be enclosed in the <Router></Router> tag
  // Inside that you would have Links and each Link would route to a path provided by a Route tag enclosed in Switch.
  // Ofcourse this mechanism of routing is the browser routing.
  // As you click on the links in the UI look at how the url in the url bar changes. You can bookmark the urls and it will show the correct UI.
  // If you use a HashRouter then the url in the url browser then the urls would look like http://localhost:3000/#/admin.
  // With a hash router the backend webserver doesnt need any changes. It is strictly frontend.
  // With the BrowserRouter you can configure the urls in your webserver to respond in ways you want.
  return (
    <Router>
      <div className="container">
        <div className="row">
          <h1 className="mt-3">
            go-movie
          </h1>
          <hr className="mb-3"></hr>
        </div>

        <div className="row">
          <div className="col-md-2">
            <nav>
              <ul className="list-group">
                <li className="list-group-item">
                  {/* Instead of using this syntax for links <a href="/">Home</a> , use the syntax for <Link to="/">Home</link> as shown below */}
                  <Link to="/">Home</Link>
                </li>
                <li className="list-group-item">
                  <Link to="/movies">Movies</Link>
                </li>
                <li className="list-group-item">
                  <Link to="/admin">Manage Catalog</Link>
                </li>
              </ul>
            </nav>
          </div>
          <div className="col-md-10">
            <Switch>
              <Route path="/movies">
                <Movies />
              </Route>
              <Route path="/admin">
                <Admin />
              </Route>
              <Route path="/">
                <Home />
              </Route>
            </Switch>
          </div>
        </div>
      </div>
    </Router>
  );
}

function Home() {
  return <h2>Home</h2>;
}
function Movies() {
  return <h2>Movies</h2>;
}
function Admin() {
  return <h2>Manage Catalog</h2>;
}
