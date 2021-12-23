// Import bootstrap css and js from the npm installed bootstrap library.
// After that install all other css and js.
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import AppContent from './AppContent';
import AppFooterFunctionalComponent from './AppFooterFunctionalComponent';
import AppHeader from './AppHeader';
import './index.css';

// Make sure to import class Component from react
class App extends Component {
  constructor(props) {
    super(props);
    // bind the `this` in handlePostChange to the context of the App component
    this.handlePostChange = this.handlePostChange.bind(this);
    this.state = {posts: []};
  }

  handlePostChange(posts) {
    this.setState({posts: posts});
  }

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
      // For the AppHeader component set the posts attribute to the posts prop defined in the App component.
      // This way the change in posts in the App component will get trickled down to the AppHeader component.
      //
      // For the AppContent pass the handlePostChange func of the App component to the AppContent component as an attribute.
      // That way when the success handler of the fetch in AppContent fires, it can execute the handlePostChange
      // function in the App component which in turn can set the state of the App component.
      // And then that state property posts can be trickled down to the AppHeader component.
      <div className="app">
        <AppHeader {...headerProps} posts={this.state.posts} />
        <AppContent posts={this.state.posts} handlePostChange={this.handlePostChange} />
        {/* <AppFooter /> */}
        <AppFooterFunctionalComponent appName="Test React App" />
      </div>
    );
  }
}

// ReactDOM.render is necessary for displaying the react component in the browser
ReactDOM.render(<App />, document.getElementById('root'));
