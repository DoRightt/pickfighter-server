#!/bin/bash

CONTAINER_NAME="dev-consul"

if [ "$(docker ps -a -q -f name=${CONTAINER_NAME})" ]; then
    echo "Existed container removing: ${CONTAINER_NAME}..."
    docker rm -f ${CONTAINER_NAME}
fi

echo "Run new container: ${CONTAINER_NAME}..."
docker run -d -p 8500:8500 -p 8600:8600/udp --name=${CONTAINER_NAME} hashicorp/consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0
