import React, { Component, Fragment } from 'react';
import { BrowserRouter as Router, Link, Route, Switch } from 'react-router-dom';
import AddOrEditMovie from './components/AddOrEditMovie';
import Admin from './components/Admin';
import Genres from './components/Genres';
import Home from './components/Home';
import Login from './components/Login';
import Movies from './components/Movies';
import OneGenre from './components/OneGenre';
import OneMovie from './components/OneMovie';

export default class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      jwt: "",
    };
    this.handleJwtChange = this.handleJwtChange.bind(this);
    this.logout = this.logout.bind(this);
  }

  handleJwtChange = (jwt) => {
    this.setState({jwt: jwt});
  };

  logout = () => {
    // Remove the jwt token from the state
    this.setState({jwt: ""});
    // Remove the jwt token from the localStorage
    window.localStorage.removeItem("jwt");
  };

  componentDidMount() {
    // If there is a valid jwt token in localStorage and the jwt token in the current state is empty
    // then pull it out of localStorage and set it
    const token = JSON.parse(window.localStorage.getItem("jwt"));
    if (token && token.length > 0 && this.state.jwt.length === 0) {
      this.setState({jwt: token});
    }
  }

  render() {
    let loginLink;
    if (this.state.jwt.length > 0) {
      loginLink = <Link to="/logout" onClick={this.logout}>Logout</Link>
    } else {
      loginLink = <Link to="/login">Login</Link>
    }

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
            <div className="col mt-3">
              <h1 className="mt-3">go-movie</h1>
            </div>
            <div className="col mt-3 text-end">
              {loginLink}
            </div>
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

                  {this.state.jwt.length > 0 && (
                    <Fragment>
                      <li className="list-group-item">
                        <Link to="/admin/movie/0">Add/edit Movie</Link>
                      </li>
                      <li className="list-group-item">
                        <Link to="/admin">Manage Catalog</Link>
                      </li>
                    </Fragment>
                  )}

                </ul>
              </nav>
            </div>
          <div className="col-md-10">
            <Switch>
              {/* This is an example of the react router rendering a component */}
              <Route path="/movies/:id" component={OneMovie} />

              {/* This is another example of the react router rendering a component */}
              <Route path="/movies">
                <Movies />
              </Route>

              {/* Note the use of keyword `exact` here in route matching. It is used because the order matters in react router */}
              <Route exact path="/genres">
                <Genres />
              </Route>

              {
                /* This is an example of the react router rendering a component with properties
                * Notice how we are using the spread operator for the props that are passed on from the Route component
                * and then adding more properties of it's own to the Genre component, such as title.
                            <Route
                              exact
                              path="/genre/:genre_name"
                              render={(props) => <Genre {...props} title={`Drama`} />}
                            />
                            */
              }

              <Route exact path="/genre/:id" component={OneGenre} />

              {/* When adding a movie, the id would be 0 in this path and when editing a movie it will have a proper id. */}
              <Route path="/admin/movie/:id" component={(props) => <AddOrEditMovie {...props} jwt={this.state.jwt} />} />

              <Route path="/admin" component={(props) => <Admin {...props} jwt={this.state.jwt} />} />

              {/*
                This is an example of the react router rendering a component with properties passed to the component.
                We are gonna need this because the props will contain the history if we want to redirect the user somewhere else.
                Also, we need to pass a function via the props to share the jwt token and set the state for jwt in the parent component of Login, ie, App
                (through lifting state like we have done here with the handleJwtChange function here).
              */}
              <Route exact path="/login" component={(props) => <Login {...props} handleJwtChange={this.handleJwtChange} />} />

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
