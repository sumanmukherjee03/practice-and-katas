import React, { Component, Fragment } from 'react';
import './AddOrEditMovie.css';
// import { Link } from 'react-router-dom';

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
      isLoaded: false,
      error: null,
    };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleSubmit = (ev) => {
    console.log("Form was submitted");
    ev.preventDefault();
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

  componentDidMount() {
  }

  render() {
    let {movie} = this.state;
    return (
      <Fragment>
        <h2>Add/Edit Movie</h2>
        <hr />
        {/* Note : We arent using a method post on the form because we want the post to be controlled by React */}
        <form onSubmit={this.handleSubmit}>
          <input type="hidden" id="id" name="id" value={movie.id} onChange={this.handleChange} />

          <div className="mb-3">
            {/* Notice the use of htmlFor instead of "for" in JSX */}
            <label htmlFor="title" className="form-label">Title</label>
            <input type="text" className="form-control" id="title" name="title" value={movie.title} onChange={this.handleChange} />
          </div>

          <div className="mb-3">
            <label htmlFor="release_date" className="form-label">Release Date</label>
            <input type="text" className="form-control" id="release_date" name="release_date" value={movie.release_date} onChange={this.handleChange} />
          </div>

          <div className="mb-3">
            <label htmlFor="runtime" className="form-label">Runtime</label>
            <input type="text" className="form-control" id="runtime" name="runtime" value={movie.runtime} onChange={this.handleChange} />
          </div>

          <div className="mb-3">
            <label htmlFor="mpaa_rating" className="form-label">MPAA Rating</label>
            <select className="form-select" id="mpaa_rating" name="mpaa_rating" value={movie.mpaa_rating} onChange={this.handleChange} >
              <option className="form-select">Choose...</option>
              <option className="form-select" value="G">G</option>
              <option className="form-select" value="PG">PG</option>
              <option className="form-select" value="PG13">PG13</option>
              <option className="form-select" value="R">R</option>
              <option className="form-select" value="NC17">NC17</option>
            </select>
          </div>

          <div className="mb-3">
            <label htmlFor="rating" className="form-label">Rating</label>
            <input type="text" className="form-control" id="rating" name="rating" value={movie.rating} onChange={this.handleChange} />
          </div>

          <div className="mb-3">
            <label htmlFor="description" className="form-label">Description</label>
            <textarea className="form-control" id="description" name="description" rows="3" onChange={this.handleChange} value={movie.description} />
          </div>

          <hr />

          <button className="btn btn-primary">Save</button>
        </form>

        <div className="mt-3">
          <pre>{JSON.stringify(this.state, null, 3)}</pre>
        </div>
      </Fragment>
    );
  }
}
