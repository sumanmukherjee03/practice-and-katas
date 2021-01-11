### How to start a postgres database

```sh
go get -u github.com/lib/pq
docker run --name postgres -e PASSWD='abcd1234' -p 5432:5432 -d upsumtech/postgres:12.4
export PGPASSWORD=abcd1234
psql -h localhost -U root -c "create database foo;"
psql -h localhost -U root -c "create user bar with encrypted password 'pqrs6789';"
psql -h localhost -U root -c "grant all privileges on database foo to bar;"
unset PGPASSWORD
psql -h localhost -U root -c "grant all privileges on database foo to bar;"
export PGPASSWORD=pqrs6789
psql -h localhost -U bar foo -c "create table accounts (user_id serial primary key, username varchar(50) unique not null, created_at timestamp not null);"
psql -h localhost -U bar foo -c "insert into accounts (username, created_at) values ('John Doe', CURRENT_TIMESTAMP);"
psql -h localhost -U bar foo -c "insert into accounts (username, created_at) values ('Jane Doe', CURRENT_TIMESTAMP);"
unset PGPASSWORD
```
