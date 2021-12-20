import React, { Component, Fragment } from 'react';
import './AppHeader.css';

export default class AppHeader extends Component {
  render() {
    // Fragment is a special react component that itself wont be generated as html but serves as a container
    // for multiple html components so that the jsx is valid for returning. Otherwise, we would have used some alternative
    // like a div. But that would have gotten displayed in html and might not be the html that we desire.
    return (
      <Fragment>
        <h1>{this.props.title}</h1>
        <hr />
      </Fragment>
    );
  }
}
