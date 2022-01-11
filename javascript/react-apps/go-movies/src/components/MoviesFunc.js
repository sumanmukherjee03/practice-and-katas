import React, { Fragment, useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

function MoviesFunc(props) {
  // This const declaration follows the pattern
  //    const [<state_variable>, <settter_func_name_for_state_variable>] = useState(<initialization_value_of_state_variable>);
  const [movies, setMovies] = useState([]);
  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);

  // useEffect is a React hook that has the same functionality as componentDidMount.
  // It runs before the component renders
  // useEffect takes an optional argument for a default value.
  // Here we must set a default value for useEffect to empty array.
  useEffect(() => {
    fetch(`${process.env.REACT_APP_API_URL}/v1/movies`)
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
        setMovies(data.movies);
        setIsLoaded(true);
      }, (error) => {
        // This is the error callback based on the returned http status
        setError(error);
        setIsLoaded(true);
      });
  }, []);

  // useEffect does what componentDidMount does. It's a hook that runs before the component renders.
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
          <h2>Choose a movie</h2>
          <div className="list-group">
            {/* Note the syntax of javascript templating here */}
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

export default MoviesFunc;
