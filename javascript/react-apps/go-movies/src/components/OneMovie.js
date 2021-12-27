import React, { Component, Fragment } from 'react';

export default class OneMovie extends Component {
  state = {movie: {}, isLoaded: false, error: null};

  componentDidMount() {
    // Notice how we retrieve the id from the url.
    // The react router makes it available to us with the property called match.
    fetch("http://localhost:4000/v1x/movie/"+this.props.match.params.id)
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
          movie: data.movie,
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
    const {movie, isLoaded, error} = this.state;

    if (!isLoaded) {
      return (
        <Fragment>
          <p>Loading...</p>
        </Fragment>
      );
    } else {
      if (!error) {
        return (
          <Fragment>
            <h2>Movie: {movie.title}</h2>
            <table className="table table-compact table-striped">
              <thead>
              </thead>
              <tbody>
                <tr>
                  <td><strong>Description:</strong></td>
                  <td>{movie.description}</td>
                </tr>
                <tr>
                  <td><strong>Runtime:</strong></td>
                  <td>{movie.runtime}</td>
                </tr>
                <tr>
                  <td><strong>Rating:</strong></td>
                  <td>{movie.rating}</td>
                </tr>
                <tr>
                  <td><strong>MPAA Rating:</strong></td>
                  <td>{movie.mpaa_rating}</td>
                </tr>
                <tr>
                  <td><strong>Genres:</strong></td>
                  <td>{Object.values(movie.movie_genres).join(",")}</td>
                </tr>
              </tbody>
            </table>
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
