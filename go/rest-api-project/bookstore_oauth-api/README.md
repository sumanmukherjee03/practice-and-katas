### Development

To get started with cassandra locally you can use the following commands

```sh
mkdir -p ~/var/data/cassandra
docker run --name cassandra -v ~/var/data/cassandra:/var/lib/cassandra -p 7000:7000 -p 9042:9042 -d cassandra:3.11
```

For the cassandra client
```sh
brew install cassandra
```

After connecting to the cassandra server through cqlsh
```sh
describe keyspaces;
CREATE KEYSPACE oauth WITH REPLICATION = {'class':'SimpleStrategy', 'replication_factor':1};
USE oauth;
DESCRIBE tables;
CREATE TABLE access_tokens(access_token VARCHAR PRIMARY KEY, user_id BIGINT, client_id BIGINT, expires BIGINT);
SELECT * FROM access_tokens;
```
