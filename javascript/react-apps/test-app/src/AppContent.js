import React, { Component, Fragment } from 'react';
import './AppContent.css';

export default class AppContent extends Component {
  // Initialize the state variable here if you do not want to have a constructor.
  // Otherwise initialize state inside the constructor func with `this.state = <whatever>`.
  // DO NOT declare your own variable called `state`.
  // Remember the state in this component is only visible to this components instances.
  // If you want to share a state between multiple components perform the lift state refactor to
  // lift the state to a common ancestor.
  //    state = {posts: []};

  constructor(props) {
    // Always call super in the constructor
    super(props);
    // Instead of referencing DOM elements by id or any unique selector for that matter, react uses the concept of refs.
    // That way a ref is a localized reference to a components elements.
    // And the reason this is important is because react's strength is to have reusable components and having an id
    // in an element which is in a component defeats the reusability purpose.
    // That said, use refs sparingly. There's always a better way to manage state.
    //    this.listRef = React.createRef();


    // bind the `this` in handlePostChange to the context of the AppContent component
    this.handlePostChange = this.handlePostChange.bind(this);
    // We would have initialized the state here if the parent component alone wasnt holding all the state
    //    this.state = {posts: []};
  }

  handlePostChange(posts) {
    // Call the handlePostChange function passed in the props of the AppContent component.
    this.props.handlePostChange(posts);
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
        // But the better way is to assign the `ref.current` to the variable.
        // So, the alternate version of the code to populate the ul referenced by this.listRef.current is as follows :
        //    const posts = this.listRef.current;
        //    data.forEach((obj) => {
        //      let li = document.createElement("li");
        //      li.appendChild(document.createTextNode(obj.title));
        //      posts.appendChild(li);
        //    });
        //  Correspondingly in the render function in the ul, you would have set the value of the listRef like so
        //    <ul ref={this.listRef}>
        //    </ul>
        //
        //  However there is an even better way to do this and that is via the state.
        //  Set the state with the data you receive back from the fetch call.
        //  And on state change do the ul > li change. This is a cleaner approach as it keeps the success callback cleaner.
        //  Remember that state has to be set with the setState method. Ofcourse this is provided that
        //  we are using some value from the state to be displayed in the UI. If the state only lives in the parent component then
        //  ofcourse this wont be necessary for us.
        //      this.setState({posts: data});

        // Call the parent components setState method indirectly via this
        this.handlePostChange(data);
      });
  }

  clickedItem = (x) => {
    console.log("clicked", x);
  }

  render() {
    // In JSX components, the onclick handler is provided in attributes as 'onClick'.
    // Other similar event handlers can be for example - onMouseEnter, onMouseLeave etc
    //
    // Note : If we did not use state and went with the ref, the ul element would have been assigned as the listRef here.
    return (
      <Fragment>
        <p>
          This is the app content
          <br />
          <button onClick={this.fetchList} className="btn btn-primary" href="#">Fetch Data</button>
        </p>
        <hr />
        {/* <p>{this.state.posts.length} items were fetched</p> */}
        <ul>
          {/* We are referring to posts via props because we arent maintaining state in this component any more. The state is passed down from the parent using props. */}
          {this.props.posts.map((elm) => (
            // Make sure to have the key attribute in the li element and value of each key needs to be unique.
            <li key={elm.id}>
              <a href="#!" onClick={() => this.clickedItem(elm.id)}>
                {elm.title}
              </a>
            </li>
          ))}
        </ul>
      </Fragment>
    );
  }
}
