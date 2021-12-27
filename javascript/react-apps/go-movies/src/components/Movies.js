import React, { Component, Fragment } from 'react';
import { Link } from 'react-router-dom';

export default class Movies extends Component {
  state = {movies: [], isLoaded: false};

  // The componentDidMount function is part of the react life cycle.
  // It trigger after an instance of the React component has been mounted into the DOM, ie after the component has rendered on the DOM.
  // It is quite an useful function to initialize the UI of a component.
  // See this link for the React lifecycle diagram - https://projects.wojtekmaj.pl/react-lifecycle-methods-diagram/
  // Also, see this link for the details on lifecycle - https://reactjs.org/docs/react-component.html
  // There are some other lifecycle methods - componentDidMount, componentDidUpdate, componentWillUnmount
  componentDidMount() {
    fetch("http://localhost:4000/v1/movies")
      .then((resp) => (resp.json()))
      .then(data => {
        this.setState({
          movies: data.movies,
          isLoaded: true,
        })
      });
  }

  render() {
    // Easy way to multi assign values from map
    const {movies, isLoaded} = this.state;
    if (!isLoaded) {
      return (
        <Fragment>
          <p>Loading</p>
        </Fragment>
      );
    } else {
      return (
        <Fragment>
          <h2>Choose a movie</h2>
          <ul>
            {movies.map((m) => (
              <li key={m.id}>
                {/* Note the syntax of javascript templating here */}
                <Link to={`/movies/${m.id}`}>{m.title}</Link>
              </li>
            ))}
          </ul>
        </Fragment>
      );
    }
  }
}
