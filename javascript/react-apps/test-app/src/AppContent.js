import React, { Component, Fragment } from 'react';
import './AppContent.css';

export default class AppContent extends Component {
  render() {
    return (
      <Fragment>
        <div>
          <h1>Hello, world!</h1>
        </div>
        <p>
          This is the app content
          <br />
          <button className="btn btn-primary" href="#">Test Button</button>
        </p>
      </Fragment>
    );
  }
}
