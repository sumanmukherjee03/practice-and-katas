import React, { Component, Fragment } from 'react';

export default class OneMovie extends Component {
  state = {movie: {}, isLoaded: false, error: null};

  componentDidMount() {
    // Notice how we retrieve the id from the url.
    // The react router makes it available to us with the property called match.
    fetch(`${process.env.REACT_APP_API_URL}/v1/movie/`+this.props.match.params.id)
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
            <div className="float-start">
              <small>Rating: {movie.mpaa_rating}</small>
            </div>
            <div className="float-end">
              {Object.values(movie.movie_genres).map((m, index) => (
                <span className="badge bg-secondary me-1" key={index}>
                  {m}
                </span>
              ))}
            </div>
            <div className="clearfix"></div>
            <hr></hr>
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
