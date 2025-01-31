#!/bin/bash

CONTAINER_NAME="dev-redis"

if [ "$(docker ps -a -q -f name=${CONTAINER_NAME})" ]; then
    echo "Existed container removing: ${CONTAINER_NAME}..."
    docker rm -f ${CONTAINER_NAME}
fi

echo "Run new container: ${CONTAINER_NAME}..."
docker run -d -p 8500:6379 --name=${CONTAINER_NAME} redis