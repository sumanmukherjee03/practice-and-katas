import React, { Component, Fragment } from 'react';
import './AppFooter.css';

export default class AppFooter extends Component {
  render() {
    const currentYear = new Date().getFullYear();
    // Fragment is a special react component that itself wont be generated as html but serves as a container
    // for multiple html components so that the jsx is valid for returning. Otherwise, we would have used some alternative
    // like a div. But that would have gotten displayed in html and might not be the html that we desire.
    return (
      <Fragment>
        <hr />
        <p className="footer">Copyright &copy; {currentYear} Test React App TLD.</p>
      </Fragment>
    );
  }
}
