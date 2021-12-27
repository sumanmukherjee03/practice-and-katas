import React from 'react';
import { BrowserRouter as Router, Link, Route, Switch } from 'react-router-dom';
import Admin from './components/Admin';
import Genres from './components/Genres';
import Home from './components/Home';
import Movies from './components/Movies';
import OneGenre from './components/OneGenre';
import OneMovie from './components/OneMovie';

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
                  <Link to="/genres">Genres</Link>
                </li>
                <li className="list-group-item">
                  <Link to="/admin">Manage Catalog</Link>
                </li>
              </ul>
            </nav>
          </div>
          <div className="col-md-10">
            <Switch>
              {/* This is an example of the react router rendering a component */}
              <Route path="/movies/:id" component={OneMovie} />

              <Route path="/movies">
                <Movies />
              </Route>

              {/* Note the use of keyword `exact` here in route matching. It is used because the order matters in react router */}
              <Route exact path="/genres">
                <Genres />
              </Route>

              {
                /* This is an example of the react router rendering a component with properties
                <Route
                  exact
                  path="/genre/:genre_name"
                  render={(props) => <Genre {...props} title={`Drama`} />}
                />
                */
              }

              <Route exact path="/genre/:genre_name" component={OneGenre} />

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

// This is a stub function for a component called movie
// function Movie() {
  // Get the path and url from the route matching
  // path lets us build routes that are relative to the parent route
  // and url instead lets us build relative links
  //   let {path, url} = useRouteMatch();
  //
  // Get the id of the movie from the params
  // let {id} = useParams();
  // return (
    // <h2>Movie id {id}</h2>
  // );
// }
