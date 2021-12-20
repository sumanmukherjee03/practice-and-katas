import React, { Component } from 'react';
import ReactDOM from 'react-dom';

// Make sure to import class Component from react
class App extends Component {
  // Every react component must have a render function that returns some JSX
  render() {
    return <div>
      <h1>Hello, world!</h1>
    </div>;
  }
}

// ReactDOM.render is necessary for displaying the react component in the browser
ReactDOM.render(<App />, document.getElementById('root'));
