import React, { Component, Fragment } from 'react';
import { Link } from 'react-router-dom';

export default class OneGenre extends Component {
  state = {movies: [], isLoaded: false, error: null, genreName: ""};

  componentDidMount() {
    // Notice how we retrieve the genre_name from the url.
    // The react router makes it available to us with the property called match.
    fetch("http://localhost:4000/v1/genre/"+this.props.match.params.id+"/movies")
      .then((response) => {
        const status = parseInt(response.status);
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
        // NOTE : Notice how the genreName came from the Link to property in Genres.js
        this.setState({
          movies: data.movies,
          isLoaded: true,
          genreName: this.props.location.genreName
        });
      }, (error) => {
        // This is the error callback based on the returned http status
        // NOTE : Notice how the genreName came from the Link to property in Genres.js
        this.setState({
          genreName: this.props.location.genreName,
          isLoaded: true,
          error
        });
      });
  }

  render() {
    // Easy way to multi assign values from map
    const {movies, isLoaded, error, genreName} = this.state;

    if (!isLoaded) {
      return (
        <Fragment>
          <p>Loading movies for genre... </p>
        </Fragment>
      );
    } else {
      if (!error) {
        return (
          <Fragment>
            <h2>Genre: {genreName}</h2>
            <div className="list-group">
              { /* Note the syntax of javascript templating here. Also, each child in a list has to have a unique value for the key property. */ }
              {movies.map((m) => (
                <Link key={m.id} to={`/movies/${m.id}`} className="list-group-item list-group-item-action">{m.title}</Link>
              ))}
            </div>
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
