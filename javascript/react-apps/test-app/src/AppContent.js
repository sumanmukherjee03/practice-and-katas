import React, { Component, Fragment } from 'react';
import './AppContent.css';

export default class AppContent extends Component {
  constructor(props) {
    // Always call super in the constructor
    super(props);
    // Instead of referencing DOM elements by id or any unique selector for that matter, react uses the concept of refs.
    // That way a ref is a localized reference to a components elements.
    // And the reason this is important is because react's strength is to have reusable components and having an id
    // in an element which is in a component defeats the reusability purpose.
    // That said, use refs sparingly. There's always a better way to manage state.
    this.listRef = React.createRef();
  }

  // The correct syntax to define functions in the class that can be called as methods is
  // funcName = () => {}. Only then can from another function in the class we can call this.funcName().
  // If the function was defined as
  // funcName() {}
  // then we cant call it from within another class method.
  fetchList = () => {
    fetch('https://jsonplaceholder.typicode.com/posts')
      .then((response) => response.json())
      .then((data) => {
        // This is how you would get access to the posts ul element without a ref.
        //    const posts = document.getElementById("post-list");
        const posts = this.listRef.current;
        data.forEach((obj) => {
          let li = document.createElement("li");
          li.appendChild(document.createTextNode(obj.title));
          posts.appendChild(li);
        });
      });
  }

  render() {
    // in JSX components, the onclick handler is provided in attributes as 'onClick'.
    // Other similar event handlers can be for example - onMouseEnter, onMouseLeave etc
    return (
      <Fragment>
        <p>
          This is the app content
          <br />
          <button onClick={this.fetchList} className="btn btn-primary" href="#">Fetch Data</button>
        </p>
        <hr />
        <ul ref={this.listRef}>
        </ul>
      </Fragment>
    );
  }
}
