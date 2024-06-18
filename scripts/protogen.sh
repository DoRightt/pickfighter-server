#!/bin/bash

# Define the path to the .proto file
PROTO_FILE="./api/fightbettr.proto"

# Check if the .proto file exists
if [ ! -f "$PROTO_FILE" ]; then
  echo "Proto file not found at path: $PROTO_FILE"
  exit 1
fi

# Run the protoc command to generate Go and gRPC code
protoc -I=./api --go_out=. --go-grpc_out=. "$PROTO_FILE"

# Check if the command was successful
if [ $? -eq 0 ]; then
  echo "Protobuf files generated successfully."
else
  echo "Failed to generate protobuf files."
  exit 1
fi
