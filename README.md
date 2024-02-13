```
# 1. 기존 컨테이너 제거

docker rm -f elasticsearch-docker

# 2. 새로운 컨테이너 실행

docker run -d -p 9200:9200 -p 9300:9300 \
-e "discovery.type=single-node" \
-e "ELASTIC_USERNAME="akaps"" \
-e "ELASTIC_PASSWORD="004"" \
--name elasticsearch-docker \
docker.elastic.co/elasticsearch/elasticsearch:7.14.0
```