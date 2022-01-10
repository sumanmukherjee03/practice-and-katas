import React, { Component, Fragment } from 'react';
import { Link } from 'react-router-dom';
import Input from './form-components/Input';
import Alert from './ui-components/Alert';

export default class Graphql extends Component {
  constructor(props) {
    super(props);
    this.state = {
      movies: [],
      isLoaded: false,
      error: null,
      alert: {
        type: "d-none",
        message: "",
      },
      searchTerm: "",
    };
    this.loadAllMovies = this.loadAllMovies.bind(this);
    this.searchMovie = this.searchMovie.bind(this);
    this.performSearch = this.performSearch.bind(this);
  }

  componentDidMount() {
    this.loadAllMovies();
  }

  loadAllMovies = () => {
    const payload = `
    {
      list {
        id
        title
        runtime
        year
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
        let movies = Object.values(data.data.list);
        this.setState({
          movies: movies,
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

  searchMovie = (ev) => {
    let val = ev.target.value;
    this.setState((prevState) => ({
      searchTerm: val,
    }));
    if (val.length > 0) {
      this.performSearch();
    } else {
      this.loadAllMovies();
    }
  }

  performSearch = () => {
    // NOTICE how we are using javascript templating here for generating the payload with arguments
    const payload = `
    {
      search(titleContains: "${this.state.searchTerm}") {
        id
        title
        runtime
        year
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
        let movies = Object.values(data.data.search);
        if (movies.length > 0) {
          this.setState({
            movies: movies,
          });
        } else {
          this.setState({
            movies: [],
          });
        }
      }, (error) => {
        // This is the error callback based on the returned http status
        const errorMsg = error.error_type + " : " + error.message;
        this.setState({
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
    const {movies, isLoaded, error} = this.state;

    if (!isLoaded) {
      return (
        <Fragment>
          <h2>Movies via graphql</h2>
          <p>Loading...</p>
        </Fragment>
      );
    } else {
      if (!error) {
        return (
          <Fragment>
            <h2>Movies via graphql</h2>
            <Input
              title={"Search"}
              type={"text"}
              name={"search"}
              value={this.state.searchTerm}
              handleChange={this.searchMovie}
            />
            <div className="list-group">
              {/* Note the syntax of javascript templating here */}
              {movies.map((m) => (
                <Link key={m.id} to={`/graphql/movies/${m.id}`} className="list-group-item list-group-item-action">{m.title}</Link>
              ))}
            </div>
          </Fragment>
        );
      } else {
        return (
          <Fragment>
            <h2>Movies via graphql</h2>
            <Alert alertType={alert.type} alertMessage={alert.message} />
          </Fragment>
        );
      }
    }
  }
}
