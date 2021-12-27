import React, { Component, Fragment } from 'react';
import { Link } from 'react-router-dom';

export default class Genres extends Component {
  state = {genres: [], isLoaded: false, error: null};

  // The componentDidMount function is part of the react life cycle.
  // It trigger after an instance of the React component has been mounted into the DOM, ie after the component has rendered on the DOM.
  // It is quite an useful function to initialize the UI of a component.
  // See this link for the React lifecycle diagram - https://projects.wojtekmaj.pl/react-lifecycle-methods-diagram/
  // Also, see this link for the details on lifecycle - https://reactjs.org/docs/react-component.html
  // There are some other lifecycle methods - componentDidMount, componentDidUpdate, componentWillUnmount
  componentDidMount() {
    fetch("http://localhost:4000/v1/genres")
      .then((response) => {
        let status = parseInt(response.status);
        if (status >= 400) {
          const contentType = response.headers.get("content-type");
          if (contentType && contentType.indexOf("application/json") < 0) {
            return response.text()
              .then(() => ({error_type: "ERROR", message: "Encountered an error with status code - " + status }))
              .then(Promise.reject.bind(Promise));
          } else {
            return response.json()
              .then((result) => result.error)
              .then(Promise.reject.bind(Promise));
          }
        }
        return response.json();
      })
      .then((data) => {
        // This is the success callback based on the returned http status
        this.setState({
          genres: data.genres,
          isLoaded: true
        });
      }, (error) => {
        // This is the error callback based on the returned http status
        this.setState({
          isLoaded: true,
          error
        });
      });
  }


  render() {
    // Easy way to multi assign values from map
    const {genres, isLoaded, error} = this.state;

    if (!isLoaded) {
      return (
        <Fragment>
          <p>Loading genres ...</p>
        </Fragment>
      );
    } else {
      if (!error) {
        return (
          <Fragment>
            <h2>Choose a genre</h2>
            <ul>
              {genres.map((g) => (
                <li key={g.id}>
                  {/* Note the syntax of javascript templating here */}
                  <Link to={`/genre/${g.id}`}>{g.genre_name}</Link>
                </li>
              ))}
            </ul>
          </Fragment>
        );
      } else {
        return (
          <Fragment>
            <p>{error.message}</p>
          </Fragment>
        );
      }
    }
  }
}
