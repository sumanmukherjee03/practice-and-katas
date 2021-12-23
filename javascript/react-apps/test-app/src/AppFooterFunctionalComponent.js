import React, { Fragment } from 'react';
import './AppFooter.css';

// Because this is a functional react component, without the use of React hooks
// we wont be able to able to access `this` keyword or the state like a normal React component.
// But then these kind of components are useful if all you want is just a dumb react component.
export default function AppFooterFunctionalComponent(props) {
  const currentYear = new Date().getFullYear();
  // Fragment is a special react component that itself wont be generated as html but serves as a container
  // for multiple html components so that the jsx is valid for returning. Otherwise, we would have used some alternative
  // like a div. But that would have gotten displayed in html and might not be the html that we desire.
  return (
    <Fragment>
      <hr />
      <p className="footer">Copyright &copy; {currentYear} {props.appName} TLD.</p>
    </Fragment>
  );
}
