### Development

For getting elasticsearch up and running locally try these commands
```sh
docker run -d --name elasticsearch -p 9200:9200 -p 9300:9300 -v $HOME/var/data/elasticsearch:/usr/share/elasticsearch/data -e "discovery.type=single-node" elasticsearch:6.8.13
curl -X GET "localhost:9200/_cat/master?v&pretty"
curl -X GET "localhost:9200/_cat/nodes?h=ip,port,heapPercent,name&pretty"
curl -X PUT -H "Content-Type: application/json" "localhost:9200/items" -d '{"settings": {"index": {"number_of_shards": 2, "number_of_replicas": 0}}}'
curl -X GET "localhost:9200/_cat/indices?v&pretty"
```
