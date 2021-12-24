import React, { Component, Fragment } from 'react';

export default class OneMovie extends Component {
  state = {movie: {}};

  componentDidMount() {
    // Notice how we retrieve the id from the url.
    // The react router makes it available to us with the property called match.
    this.setState({
      movie: {
        id: this.props.match.params.id,
        title: 'Some movie',
        runtime: 150
      }
    });
  }

  render() {
    return (
      <Fragment>
        <h2>Movie: {this.state.movie.id}</h2>
        <table className="table table-compact table-striped">
          <thead>
          </thead>
          <tbody>
            <tr>
              <td><strong>Title:</strong></td>
              <td>{this.state.movie.title}</td>
            </tr>
            <tr>
              <td><strong>Runtime:</strong></td>
              <td>{this.state.movie.runtime}</td>
            </tr>
          </tbody>
        </table>
      </Fragment>
    );
  }
}
