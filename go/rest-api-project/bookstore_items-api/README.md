### Development

For getting elasticsearch up and running locally try these commands
```sh
docker run -d --name elasticsearch -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" elasticsearch:7.10.1
curl -X GET "localhost:9200/_cat/master?v&pretty"
curl -X GET "localhost:9200/_cat/nodes?h=ip,port,heapPercent,name&pretty"
```
