import React, { Fragment, useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

function AdminFunc(props) {
  // This const declaration follows the pattern
  //    const [<state_variable>, <settter_func_name_for_state_variable>] = useState(<initialization_value_of_state_variable>);
  const [movies, setMovies] = useState([]);
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
    // This is to redirect if the user is not logged in
    // DO NOT use componentWillMount lifecycle hook for doing this because that is getting deprecated.
    if (props.jwt.length === 0) {
      props.history.push({
        pathname: "/login",
      });
      return;
    }
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
  });


  // React render logic
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
              <Link key={m.id} to={`/admin/movie/${m.id}`} className="list-group-item list-group-item-action">{m.title}</Link>
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

export default AdminFunc;
