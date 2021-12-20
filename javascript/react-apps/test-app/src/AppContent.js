import React, { Component, Fragment } from 'react';
import './AppContent.css';

export default class AppContent extends Component {
  fetchList = () => {
    fetch('https://jsonplaceholder.typicode.com/posts')
      .then((response) => response.json())
      .then((data) => {
        let posts = document.getElementById("post-list");
        data.forEach((obj) => {
          let li = document.createElement("li");
          li.appendChild(document.createTextNode(obj.title));
          posts.appendChild(li);
        });
      });
  }

  render() {
    // in JSX components, the onclick handler is provided in attributes as 'onClick'.
    return (
      <Fragment>
        <p>
          This is the app content
          <br />
          <button onClick={this.fetchList} className="btn btn-primary" href="#">Fetch Data</button>
        </p>
        <hr />
        <ul id="post-list">
        </ul>
      </Fragment>
    );
  }
}
