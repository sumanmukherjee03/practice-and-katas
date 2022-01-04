import React, { Component, Fragment } from 'react';
import { Link } from 'react-router-dom';
import Input from './form-components/Input';
import Alert from './ui-components/Alert';

export default class Login extends Component {
  constructor(props) {
    super(props);
    this.state = {
      email: "",
      password: "",
      error: null, // Mainly to hold submit error from backend
      errors: [], // Mainly to hold form errors
      alert: {
        type: "d-none", // Initially do not display the alert notification div
        message: "",
      },
    };

    this.handleInputElmChange = this.handleInputElmChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleJwtChange = this.handleJwtChange.bind(this);
  }

  handleInputElmChange = (ev) => {
    let value = ev.target.value;
    let name = ev.target.name;

    // NOTICE how the variable `name` was used as the key for the map here
    // Also notice the use of the callback func when using setState.
    this.setState((prevState) => ({
      ...prevState,
      [name]: value,
    }));
  }

  // This function is called by the login success handler and this func in turn calls
  // the parent App components handleJwtChange func passed down to this component through props.
  // This is necessary to maintain the state of the jwt token in the parent component.
  // That way the jwt token can be passed down to other components from App through props.
  handleJwtChange = (token) => {
    this.props.handleJwtChange(token);
  }

  handleSubmit = (ev) => {
    ev.preventDefault();

    // Client side form validation - check for errors in input elements of the form
    // and if there are any then update the state with that info
    let errors = [];
    if (this.state.email.length === 0) {
      errors.push("email");
    }
    if (this.state.password.length === 0) {
      errors.push("password");
    }
    this.setState({errors: errors});
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
    fetch("http://localhost:4000/v1/signin", reqOptions)
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
        this.handleJwtChange(data.response);
        // NOTICE here how we are redirecting the user to a different screen - the manage catalogue screen.
        // This is how we do redirects in a single page app.
        this.props.history.push({
          pathname: "/admin",
        });
      }, (error) => {
        // This is the error callback based on the returned http status
        const message = error.error_type + " : " + error.message;
        this.setState({
          alert: {
            type: "alert-danger",
            message: message,
          },
        });
      });
  }

  hasError = (key) => {
    return this.state.errors.indexOf(key) >= 0;
  }

  render() {
    const {email, password, alert} = this.state;

    return (
      <Fragment>
        <h2>Login</h2>
        <hr />

        <Alert alertType={alert.type} alertMessage={alert.message} />

        <form className="pt-3" onSubmit={this.handleSubmit}>
          <Input
            title={"Email"}
            type={"text"}
            name={"email"}
            value={email}
            handleChange={this.handleInputElmChange}
            className={this.hasError("email") ? "is-invalid" : ""}
            errorDiv={this.hasError("email") ? "text-danger" : "d-none"}
            errorMsg={"Please enter a valid email"}
          />

          <Input
            title={"Password"}
            type={"text"}
            name={"password"}
            value={password}
            handleChange={this.handleInputElmChange}
            className={this.hasError("password") ? "is-invalid" : ""}
            errorDiv={this.hasError("password") ? "text-danger" : "d-none"}
            errorMsg={"Please enter a valid password"}
          />

          <hr />

          <button className="btn btn-primary">Login</button>
          <Link to="/admin" className="btn btn-warning ms-1">Cancel</Link>
        </form>
      </Fragment>
    );
  }
}
