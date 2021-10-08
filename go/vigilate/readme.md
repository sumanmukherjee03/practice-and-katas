
# Vigilate

This is the source code for the second project in the Udemy course Working with Websockets in Go (Golang).

A dead simple monitoring service, intended to replace things like Nagios.

## Build

Build in the normal way on Mac/Linux:

```
go build -o vigilate cmd/web/*.go
```

## Requirements

Vigilate requires:
- Postgres 11 or later (db is set up as a repository, so other databases are possible)
- An account with [Pusher](https://pusher.com/), or a Pusher alternative
(like [ipê](https://github.com/dimiro1/ipe))

For running migrations you will be needing `soda` from https://gobuffalo.io/en/docs/db/toolbox/
If the mac installation via brew seems to be causing problems download the source code from
https://github.com/tsawler/vigilate/releases/tag/v1 and move the binary from soda directory in that project to a bin dir

## Run

First, make sure ipê is running (if you're using ipê):

On Mac/Linux
```
cd ipe
./ipe
```

Run the migrations with soda
```
soda migrate
```

Run with flags:

```
./vigilate \
-dbuser='root' \
-dbpass='some_password' \
-pusherHost='localhost' \
-pusherPort='4001' \
-pusherKey='abc123' \
-pusherSecret='123abc' \
-pusherApp="1" \
-pusherSecure=false
```

## All Flags

```
Usage of ./vigilate:
  -db string
        database name (default "vigilate")
  -dbhost string
        database host (default "localhost")
  -dbport string
        database port (default "5432")
  -dbssl string
        database ssl setting (default "disable")
  -dbuser string
        database user
  -dbpass string
        database password
  -domain string
        domain name (e.g. example.com) (default "localhost")
  -identifier string
        unique identifier (default "vigilate")
  -port string
        port to listen on (default ":4000")
  -production
        application is in production
  -pusherApp string
        pusher app id (default "9")
  -pusherHost string
        pusher host
  -pusherKey string
        pusher key
  -pusherPort string
        pusher port (default "443")
  -pusherSecret string
        pusher secret
   -pusherSecure
        pusher server uses SSL (true or false)
```

### Some soda commands

For the migrations, you can follow the fizz format and documentation for it is available here :
  - https://github.com/gobuffalo/fizz

These are some helpful soda commands to help with models and migrations.
```
soda migrate
soda generate fizz CreateHostsTable
```
