// Import bootstrap css and js from the npm installed bootstrap library.
// After that install all other css and js.
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import AppContent from './AppContent';
import AppFooter from './AppFooter';
import AppHeader from './AppHeader';
import './index.css';

// Make sure to import class Component from react
class App extends Component {
  // Every react component must have a render function that returns some JSX
  render() {
    const headerProps = {
      title: "Test React App",
      subject: "Frontend apps",
      favColor: "red",
    };
    // class is a reserved word in JS, that's why in react components we cant use a "class" for css.
    // Some of these html attributes are slightly different in React. Instead of "class" for example we use "className".
    // Remember that we can pass multiple properties by passing an object with the JS spread operator `...`
    return (
      <div className="app">
        <AppHeader {...headerProps}/>
        <AppContent />
        <AppFooter />
      </div>
    );
  }
}

// ReactDOM.render is necessary for displaying the react component in the browser
ReactDOM.render(<App />, document.getElementById('root'));
