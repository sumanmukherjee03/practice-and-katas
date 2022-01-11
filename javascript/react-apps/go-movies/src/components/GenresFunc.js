import React, { Fragment, useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

function GenresFunc(props) {
  const [genres, setGenres] = useState([]);
  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);

  useEffect(() => {
    fetch(`${process.env.REACT_APP_API_URL}/v1/genres`)
      .then((response) => {
        let status = parseInt(response.status);
        if (status >= 400) {
          const contentType = response.headers.get("content-type");
          if (contentType && contentType.indexOf("application/json") < 0) {
            return response.text()
              .then(() => ({error_type: "ERROR", message: "Encountered an error with status code - " + status }))
              .then(Promise.reject.bind(Promise));
          } else {
            return response.json()
              .then((result) => result.error)
              .then(Promise.reject.bind(Promise));
          }
        }
        return response.json();
      })
      .then((data) => {
        // This is the success callback based on the returned http status
        setGenres(data.genres);
        setIsLoaded(true);
      }, (error) => {
        // This is the error callback based on the returned http status
        setError(error);
        setIsLoaded(true);
      });
  }, [])


  if (!isLoaded) {
    return (
      <Fragment>
        <p>Loading genres ...</p>
      </Fragment>
    );
  } else {
    if (!error) {
      return (
        <Fragment>
          <h2>Choose a genre</h2>
          <div className="list-group">
            {/*
              Note the syntax of javascript templating here.
              Also, worth noting is the fact that the Link to property can be an object.
              The "pathname" key/value in that object is mandatory for the link to work.
              However, you can add other additional pproperties as necessary.
              In the component that Link calls, these values are available under `this.props.location`.
            */}
            {genres.map((g) => (
              <Link
                to={
                  {
                    pathname: `/genre/${g.id}`,
                    genreName: g.genre_name,
                  }
                }
                className="list-group-item list-group-item-action"
                key={g.id}
              >{g.genre_name}</Link>
            ))}
          </div>
        </Fragment>
      );
    } else {
      return (
        <Fragment>
          <p>{error.message}</p>
        </Fragment>
      );
    }
  }
}

export default GenresFunc
