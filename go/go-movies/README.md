## Development

To dump the database seed file first create the DB called `go_movies`. Then from the cli execute
```
psql -h localhost -U root -d go_movies -f go_movies.sql
```

To generate the jwt secret follow this url - https://go.dev/play/p/s8KlqJIOWej

To start the application follow the instructions below
```
export GO_MOVIES_JWT_SECRET=<the_jwt_secret_obtained_from_go_platground_url>
go run cmd/api/*.go --db-password <db_password> --the-movie-db-api-key=c79fd5e6585b53309e4474cdaa0a8b7
```

Also, start the npm server using `npm start` from the javascript react app.

Log into the UI using - john.doe@example.org/password

To build the code for production use the example command below
```
env GOOS=linux GOARCH=amd64 go build -o gomovies ./cmd/api
```
