## Development

To dump the database seed file first create the DB called `go_movies`. Then from the cli execute
```
psql -h localhost -U root -d go_movies -f go_movies.sql
```

To start the application
```
go run cmd/api/*.go --db-password <db_password> --the-movie-db-api-key=c79fd5e6585b53309e4474cdaa0a8b7
```

Also, start the npm server using `npm start` from the javascript react app.

Log into the UI using - john.doe@example.org/password
