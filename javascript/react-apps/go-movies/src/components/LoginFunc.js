import React, { Fragment, useState } from 'react';
import { Link } from 'react-router-dom';
import Input from './form-components/Input';
import Alert from './ui-components/Alert';

function LoginFunc(props) {
  const defaultAlert = {
    type: "d-none",
    message: "",
  };

  // This const declaration follows the pattern
  //    const [<state_variable>, <settter_func_name_for_state_variable>] = useState(<initialization_value_of_state_variable>);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState(null);
  const [errors, setErrors] = useState([]);
  const [alert, setAlert] = useState(defaultAlert);

  const handleInputElmChange = (ev) => {
    let value = ev.target.value;
    let name = ev.target.name;
    switch(name) {
      case "email":
        setEmail(value);
        break;
      case "password":
        setPassword(value);
        break;
      default:
        console.log("no element to handle change");
    }
  };

  // This function is called by the login success handler and this func in turn calls
  // the parent App components handleJwtChange func passed down to this component through props.
  // This is necessary to maintain the state of the jwt token in the parent component.
  // That way the jwt token can be passed down to other components from App through props.
  const handleJwtChange = (token) => {
    props.handleJwtChange(token);
  };

  const handleSubmit = (ev) => {
    ev.preventDefault();

    // Client side form validation - check for errors in input elements of the form
    // and if there are any then update the state with that info
    let errors = [];
    if (email.length === 0) {
      errors.push("email");
    }
    if (password.length === 0) {
      errors.push("password");
    }
    setErrors(errors);
    if (errors.length > 0) {
      return false;
    }

    // First of all, get the form data from the form element
    const data = new FormData(ev.target);
    // Then convert the form data into a javascript object
    const payload = Object.fromEntries(data.entries());
    // For POST requests use this request options, so that it can be passed to the fetch call
    const reqOptions = {
      method: 'POST',
      body: JSON.stringify(payload),
    }
    fetch(`${process.env.REACT_APP_API_URL}/v1/signin`, reqOptions)
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
        const token = data.response;
        // Call the handleJwtChange func for this component which will in turn populate
        // the jwt token in the state of the parent component, ie App
        handleJwtChange(token);
        // Store the login information, ie the token in localStorage
        // We could have also used a cookie, but most modern browsers have localStorage now a days
        window.localStorage.setItem("jwt", JSON.stringify(token));
        // NOTICE here how we are redirecting the user to a different screen - the manage catalogue screen.
        // This is how we do redirects in a single page app.
        props.history.push({
          pathname: "/admin",
        });
      }, (error) => {
        // This is the error callback based on the returned http status
        const message = error.error_type + " : " + error.message;
        setError({
          type: "alert-danger",
          message: message,
        });
      });
  };

  const hasError = (key) => {
    return errors.indexOf(key) >= 0;
  };

  return (
    <Fragment>
      <h2>Login</h2>
      <hr />

      <Alert alertType={alert.type} alertMessage={alert.message} />

      <form className="pt-3" onSubmit={handleSubmit}>
        <Input
          title={"Email"}
          type={"text"}
          name={"email"}
          value={email}
          handleChange={handleInputElmChange}
          className={hasError("email") ? "is-invalid" : ""}
          errorDiv={hasError("email") ? "text-danger" : "d-none"}
          errorMsg={"Please enter a valid email"}
        />

        <Input
          title={"Password"}
          type={"text"}
          name={"password"}
          value={password}
          handleChange={handleInputElmChange}
          className={hasError("password") ? "is-invalid" : ""}
          errorDiv={hasError("password") ? "text-danger" : "d-none"}
          errorMsg={"Please enter a valid password"}
        />

        <hr />

        <button className="btn btn-primary">Login</button>
        <Link to="/admin" className="btn btn-warning ms-1">Cancel</Link>
      </form>
    </Fragment>
  );
}

export default LoginFunc;
