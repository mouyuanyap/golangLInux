#!/bin/sh

docker stop $(docker ps -aq)
#docker stop $(docker ps -aq)
#docker stop $(docker ps -aq)
docker load -i /root/dockerimage/gin.tar
docker run --publish 8080:8080 --rm gin
