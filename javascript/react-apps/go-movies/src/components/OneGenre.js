import React, { Component, Fragment } from 'react';
import { Link } from 'react-router-dom';

export default class OneGenre extends Component {
  state = {genre: {}, isLoaded: false, error: null};

  componentDidMount() {
    // Notice how we retrieve the genre_name from the url.
    // The react router makes it available to us with the property called match.
    fetch("http://localhost:4000/v1/genre/"+this.props.match.params.genre_name)
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
        this.setState({
          genre: data.genre,
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
    const {genre, isLoaded, error} = this.state;

    if (!isLoaded) {
      return (
        <Fragment>
          <p>Loading genre... </p>
        </Fragment>
      );
    } else {
      if (!error) {
        return (
          <Fragment>
            <h2>Genre: {genre.genre_name}</h2>
            <ul>
              {genre.movie_genres.map((mg) => (
                <li key={mg.movie.id}>
                  {/* Note the syntax of javascript templating here */}
                  <Link to={`/movies/${mg.movie.id}`}>{mg.movie.title}</Link>
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
