// Import bootstrap css and js from the npm installed bootstrap library.
// After that install all other css and js.
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import AppContent from './AppContent';
import AppFooter from './AppFooter';
import './index.css';


// Make sure to import class Component from react
class App extends Component {
  // Every react component must have a render function that returns some JSX
  render() {
    // class is a reserved word in JS, that's why in react components we cant use a "class" for css.
    // Some of these html attributes are slightly different in React. Instead of "class" for example we use "className".
    return (
      <div className="app">
        <AppContent />
        <AppFooter />
      </div>
    );
  }
}

// ReactDOM.render is necessary for displaying the react component in the browser
ReactDOM.render(<App />, document.getElementById('root'));
