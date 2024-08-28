#!/bin/bash

CONTAINER_NAME="test-db"
POSTGRES_IMAGE="pf-test-postgres"
POSTGRES_PASSWORD="password"

if [ "$(docker ps -a -q -f name=${CONTAINER_NAME})" ]; then
    echo "Removing existed container: ${CONTAINER_NAME}..."
    docker rm -f ${CONTAINER_NAME}
fi

echo "Starting new container: ${CONTAINER_NAME}..."
docker build -t ${POSTGRES_IMAGE} ./tests
docker run -d --name=${CONTAINER_NAME} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -p 5433:5432 ${POSTGRES_IMAGE}

docker ps -a | grep ${CONTAINER_NAME}