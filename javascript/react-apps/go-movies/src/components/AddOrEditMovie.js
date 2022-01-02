import React, { Component, Fragment } from 'react';
import './AddOrEditMovie.css';
import Input from './form-components/Input';
import Select from './form-components/Select';
import TextArea from './form-components/TextArea';
import Alert from './ui-components/Alert';

export default class AddOrEditMovie extends Component {
  constructor(props) {
    super(props);
    this.state = {
      movie: {
        id: 0,
        title: "",
        release_date: "",
        runtime: "",
        mpaa_rating: "",
        rating: "",
        description: "",
      },
      mpaa_options: {
        G: "G",
        PG: "PG",
        PG13: "PG13",
        R: "R",
        NC17: "NC17",
      },
      isLoaded: false,
      loadingError: null,
      alert: {
        type: "d-none",
        message: "",
      },
      errors: [],
    };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.hasError = this.hasError.bind(this);
  }

  handleSubmit = (ev) => {
    ev.preventDefault();

    // Client side form validation - check for errors in input elements of the form
    // and if there are any then update the state with that info
    let errors = [];
    if (this.state.movie.title.length === 0) {
      errors.push("title");
    }
    if (this.state.movie.release_date.length === 0) {
      errors.push("release_date");
    }
    if (this.state.movie.runtime.length === 0) {
      errors.push("runtime");
    }
    if (this.state.movie.description.length === 0) {
      errors.push("description");
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
    fetch("http://localhost:4000/v1/admin/movie/edit", reqOptions)
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
        this.setState({
          movie: data.movie,
          isLoaded: true,
          alert: {
            type: "alert-success",
            message: "Successfully submitted data to the backend",
          },
        });
      }, (error) => {
        // This is the error callback based on the returned http status
        const errorMsg = error.error_type + " : " + error.message;
        this.setState({
          isLoaded: true,
          alert: {
            type: "alert-danger",
            message: errorMsg,
          },
        });
      });
  }

  handleChange = (ev) => {
    let val = ev.target.value;
    let name = ev.target.name;
    // NOTE : setState can take a value of a state directly, or it can take a callback.
    // And the callbac function takes the previous state as an argument. You can use the previous state
    // to update the value of the current state.
    // Notice how we are using the javascript spread operator here to grab the values from the previous state
    // and using that to update the value of the current state with the name and value of the input that just changed.
    this.setState((prevState) => ({
      movie: {
        ...prevState.movie,
        [name]: val,
      },
    }));
  }

  hasError = (key) => {
    return this.state.errors.indexOf(key) !== -1;
  }

  componentDidMount() {
    const id = this.props.match.params.id;
    if (id > 0) {
      // Notice how we retrieve the id from the url.
      // The react router makes it available to us with the property called match.
      fetch("http://localhost:4000/v1/movie/"+id)
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
          this.setState({
            movie: data.movie,
            isLoaded: true
          });
        }, (error) => {
          // This is the error callback based on the returned http status
          this.setState({
            isLoaded: true,
            loadingError: error,
          });
        });
    } else {
      this.setState({isLoaded: true})
    }
  }

  render() {
    // Easy way to multi assign values from map
    const {movie, isLoaded, loadingError, alert} = this.state;

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
            <form onSubmit={this.handleSubmit}>
              <input type="hidden" id="id" name="id" value={movie.id} onChange={this.handleChange} />

              <Input
                title={"Title"}
                type={"text"}
                name={"title"}
                value={movie.title}
                handleChange={this.handleChange}
                className={this.hasError("title") ? "is-invalid" : ""}
                errorDiv={this.hasError("title") ? "text-danger" : "d-none"}
                errorMsg={"Please enter a title"}
              />

              <Input
                title={"Release Date"}
                type={"text"}
                name={"release_date"}
                value={movie.release_date}
                handleChange={this.handleChange}
                className={this.hasError("title") ? "is-invalid" : ""}
                errorDiv={this.hasError("title") ? "text-danger" : "d-none"}
                errorMsg={"Please enter a valid release date"}
              />

              <Input
                title={"Runtime"}
                type={"text"}
                name={"runtime"}
                value={movie.runtime}
                handleChange={this.handleChange}
                className={this.hasError("title") ? "is-invalid" : ""}
                errorDiv={this.hasError("title") ? "text-danger" : "d-none"}
                errorMsg={"Please enter a valid runtime"}
              />

              <Select
                title={"MPAA Rating"}
                name={"mpaa_rating"}
                value={movie.mpaa_rating}
                handleChange={this.handleChange}
                placeholder="Choose..."
                options={this.state.mpaa_options}
              />

              <Input
                title={"Rating"}
                type={"text"}
                name={"rating"}
                value={movie.rating}
                handleChange={this.handleChange}
              />

              <TextArea
                title={"Description"}
                name={"description"}
                value={movie.description}
                handleChange={this.handleChange}
                className={this.hasError("title") ? "is-invalid" : ""}
                errorDiv={this.hasError("title") ? "text-danger" : "d-none"}
                errorMsg={"Please enter a description"}
              />

              <hr />

              <button className="btn btn-primary">Save</button>
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
}
