import React, { Component, Fragment } from 'react';

export default class OneMovieGraphql extends Component {
  state = {movie: {}, isLoaded: false, error: null};

  componentDidMount() {

    const payload = `
    {
      movie(id: ${this.props.match.params.id}){
        id
        title
        description
        runtime
        year
        release_date
        rating
        mpaa_rating
        created_at
        updated_at
        poster
      }
    }
    `;

    const reqHeaders = new Headers();
    reqHeaders.append("Content-Type", "application/json");

    const reqOpts = {
      method: 'POST',
      body: payload,
      headers: reqHeaders,
    };

    fetch("http://localhost:4000/v1/graphql", reqOpts)
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
        let movie = data.data.movie;
        this.setState({
          movie: movie,
          isLoaded: true,
        });
      }, (error) => {
        // This is the error callback based on the returned http status
        const errorMsg = error.error_type + " : " + error.message;
        this.setState({
          isLoaded: true,
          error: errorMsg,
          alert: {
            type: "alert-danger",
            message: errorMsg,
          },
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
            <h2>Movie: {movie.title} ({movie.year})</h2>
            {movie.poster.length > 0 && (
              <div>
                <img src={`https://image.tmdb.org/t/p/w200${movie.poster}`} alt="poster" />
              </div>
            )}
            <div className="float-start">
              <small>Rating: {movie.mpaa_rating}</small>
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

