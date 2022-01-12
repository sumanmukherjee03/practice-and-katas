import React, { Fragment, useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

function OneGenreFunc(props) {
  // This const declaration follows the pattern
  //    const [<state_variable>, <settter_func_name_for_state_variable>] = useState(<initialization_value_of_state_variable>);
  const [movies, setMovies] = useState([]);
  const [genreName, setGenreName] = useState("");
  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);

  // useEffect is a React hook that has the same functionality as componentDidMount and componentDidUpdate.
  // As in, if you performed the same side effect in both of those lifecycle methods in React Component class then
  // the equivalent of that would go in the useEffect hook when using the functional verison of React.
  // useEffect runs after the component renders - both during the first load and after an update.
  // useEffect takes an optional argument for optimization purposes. This is similar to `prevProps` and `prevState` in Component class componentDidUpdate.
  // There you only perform a side effect if a certain key in the state is different from the previous state.
  // This parameter is called the dependency param. Since useEffect is called on every render,
  // if you want to render only once you pass it an empty dependency array. Otherwise you can pass it an array of dependencies like
  // [props, state] or [props.match.id] or [state.movies] etc.
  // useEffect here compares these dependencies with the current state, props, depencies etc
  // with the previous props or state and if that hasnt changed, then dont re-run this function again.
  //
  // Read more about the useEffect hook here - https://reactjs.org/docs/hooks-effect.html
  useEffect(() => {
    // Notice how we retrieve the genre_name from the url.
    // The react router makes it available to us with the property called match.
    fetch(`${process.env.REACT_APP_API_URL}/v1/genre/`+props.match.params.id+"/movies")
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
        setMovies(data.movies);
        setIsLoaded(true);
        setGenreName(props.location.genreName);
      }, (error) => {
        // This is the error callback based on the returned http status
        // NOTE : Notice how the genreName came from the Link to property in Genres.js
        setError(error);
        setIsLoaded(true);
        setGenreName(props.location.genreName);
      });
  }, [props]); // Note : We are passing the entire props in the dependencies array because we extract the id from url match as well as genreName from location

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

export default OneGenreFunc;
