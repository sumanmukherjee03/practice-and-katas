import React, { Fragment, useEffect, useState } from 'react';
import { confirmAlert } from 'react-confirm-alert';
import 'react-confirm-alert/src/react-confirm-alert.css';
import { Link } from 'react-router-dom';
import './AddOrEditMovie.css';
import Input from './form-components/Input';
import Select from './form-components/Select';
import TextArea from './form-components/TextArea';
import Alert from './ui-components/Alert';

function AddOrEditMovieFunc(props) {
  const defaultMovie = {
    id: 0,
    title: "",
    release_date: "",
    runtime: "",
    mpaa_rating: "",
    rating: "",
    description: "",
  };
  const mpaaOptions = {
    G: "G",
    PG: "PG",
    PG13: "PG13",
    R: "R",
    NC17: "NC17",
  };
  const defaultAlert = {
    type: "d-none",
    message: "",
  };

  // This const declaration follows the pattern
  //    const [<state_variable>, <settter_func_name_for_state_variable>] = useState(<initialization_value_of_state_variable>);
  const [movie, setMovie] = useState(defaultMovie);
  const [errors, setErrors] = useState([]);
  const [loadingError, setLoadingError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [alert, setAlert] = useState(defaultAlert);

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

    const id = props.match.params.id;
    if (id > 0) {
      // Notice how we retrieve the id from the url.
      // The react router makes it available to us with the property called match.
      fetch(`${process.env.REACT_APP_API_URL}/v1/movie/`+id)
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
                .then((result) => result.loadingError)
                .then(Promise.reject.bind(Promise));
            }
          }
          return response.json();
        })
        .then((data) => {
          // This is the success callback based on the returned http status
          const releaseDate = new Date(data.movie.release_date);
          data.movie.release_date = releaseDate.toISOString().split("T")[0];
          setMovie(data.movie);
          setIsLoaded(true);
        }, (error) => {
          // This is the error callback based on the returned http status
          setIsLoaded(true);
          setLoadingError(error);
        });
    } else {
      setIsLoaded(true);
    }
  }, [props]);




  const handleSubmit = (ev) => {
    ev.preventDefault();

    // Client side form validation - check for errors in input elements of the form
    // and if there are any then update the state with that info
    let errors = [];
    if (movie.title.length === 0) {
      errors.push("title");
    }
    if (movie.release_date.length === 0) {
      errors.push("release_date");
    }
    if (movie.runtime.length === 0) {
      errors.push("runtime");
    }
    if (movie.description.length === 0) {
      errors.push("description");
    }
    setErrors(errors);
    if (errors.length > 0) {
      return false;
    }

    // First of all, get the form data from the form element
    const data = new FormData(ev.target);
    // Then convert the form data into a javascript object
    const payload = Object.fromEntries(data.entries());
    // Add request headers for authentication
    const headers = new Headers();
    headers.append("Content-Type", "application/json");
    headers.append("Authorization", "Bearer " + props.jwt);

    // For POST requests use this request options, so that it can be passed to the fetch call
    const reqOptions = {
      method: 'POST',
      body: JSON.stringify(payload),
      headers: headers,
    };

    fetch(`${process.env.REACT_APP_API_URL}/v1/admin/movie/edit`, reqOptions)
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
        // NOTICE here how we are redirecting the user to a different screen - the manage catalogue screen.
        // This is how we do redirects in a single page app.
        props.history.push({
          pathname: "/admin",
        });
      }, (error) => {
        // This is the error callback based on the returned http status
        const errorMsg = error.error_type + " : " + error.message;
        setAlert({
          type: "alert-danger",
          message: errorMsg,
        });
      });
  };

  const handleChange = (ev) => {
    let val = ev.target.value;
    let name = ev.target.name;
    // Notice how we are using the javascript spread operator here to grab the values from the previous version of the state variable movie
    // and using that to update the value of the current state of the variable movie with the name and value of the input that just changed.
    setMovie({
      ...movie,
      [name]: val,
    });
  };

  const hasError = (key) => {
    return errors.indexOf(key) !== -1;
  };

  const handleDelete = () => {
    const id = props.match.params.id;

    // Add request headers for authentication
    const headers = new Headers();
    headers.append("Content-Type", "application/json");
    headers.append("Authorization", "Bearer " + props.jwt);

    const reqOpts = {
      method: 'DELETE',
      headers: headers,
    };

    fetch(`${process.env.REACT_APP_API_URL}/v1/admin/movie/${id}/delete`, reqOpts)
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
        // NOTICE here how we are redirecting the user to a different screen - the manage catalogue screen.
        // This is how we do redirects in a single page app.
        props.history.push({
          pathname: "/admin",
        });
      }, (error) => {
        // This is the error callback based on the returned http status
        const errorMsg = error.error_type + " : " + error.message;
        setAlert({
          type: "alert-danger",
          message: errorMsg,
        });
      });
  };

  const confirmDelete = (ev) => {
    confirmAlert({
      title: 'Delete movie',
      message: 'Are you sure you want to proceed with the delete?',
      buttons: [
        {
          label: 'Yes',
          onClick: () => handleDelete()
        },
        {
          label: 'No',
          onClick: () => console.log("Traitor - you escaped this time")
        }
      ]
    });
  };


  // Code for rendering the component
  if (!isLoaded) {
    return (
      <Fragment>
        <p>Loading...</p>
      </Fragment>
    );
  } else {
    if (!loadingError) {
      return (
        <Fragment>
          <h2>Add/Edit Movie</h2>
          <Alert alertType={alert.type} alertMessage={alert.message} />
          <hr />
          {/* Note : We arent using a method post on the form because we want the post to be controlled by React */}
          <form onSubmit={handleSubmit}>
            <input type="hidden" id="id" name="id" value={movie.id} onChange={handleChange} />

            <Input
              title={"Title"}
              type={"text"}
              name={"title"}
              value={movie.title}
              handleChange={handleChange}
              className={hasError("title") ? "is-invalid" : ""}
              errorDiv={hasError("title") ? "text-danger" : "d-none"}
              errorMsg={"Please enter a title"}
            />

            <Input
              title={"Release Date"}
              type={"text"}
              name={"release_date"}
              value={movie.release_date}
              handleChange={handleChange}
              className={hasError("release_date") ? "is-invalid" : ""}
              errorDiv={hasError("release_date") ? "text-danger" : "d-none"}
              errorMsg={"Please enter a valid release date"}
            />

            <Input
              title={"Runtime"}
              type={"text"}
              name={"runtime"}
              value={movie.runtime}
              handleChange={handleChange}
              className={hasError("runtime") ? "is-invalid" : ""}
              errorDiv={hasError("runtime") ? "text-danger" : "d-none"}
              errorMsg={"Please enter a valid runtime"}
            />

            <Select
              title={"MPAA Rating"}
              name={"mpaa_rating"}
              value={movie.mpaa_rating}
              handleChange={handleChange}
              placeholder="Choose..."
              options={mpaaOptions}
            />

            <Input
              title={"Rating"}
              type={"text"}
              name={"rating"}
              value={movie.rating}
              handleChange={handleChange}
            />

            <TextArea
              title={"Description"}
              name={"description"}
              value={movie.description}
              handleChange={handleChange}
              className={hasError("description") ? "is-invalid" : ""}
              errorDiv={hasError("description") ? "text-danger" : "d-none"}
              errorMsg={"Please enter a description"}
            />

            <hr />

            <button className="btn btn-primary">Save</button>
            <Link to="/admin" className="btn btn-warning ms-1">Cancel</Link>

            {/* NOTICE the way conditional rendering is being done in React. Also, note that we arent using the Link tag. Instead using the "a" tag. */}
            {movie.id > 0 && (
                <a href="#!" onClick={() => confirmDelete()} className="btn btn-danger ms-1">Delete</a>
            )}
        </form>

          {/*
            One easy way to visualize the current state is via this :
                <div className="mt-3">
                  <pre>{JSON.stringify(this.state, null, 3)}</pre>
                </div>
                */}
        </Fragment>
      );
    } else {
      return (
        <Fragment>
          <p>{loadingError.message}</p>
        </Fragment>
      );
    }
  }
}

export default AddOrEditMovieFunc;
