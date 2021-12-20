import React, { Component } from 'react';
import ReactDOM from 'react-dom';
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
        <div>
          <h1>Hello, world!</h1>
        </div>
        <AppFooter />
      </div>
    );
  }
}

// ReactDOM.render is necessary for displaying the react component in the browser
ReactDOM.render(<App />, document.getElementById('root'));
