#!/bin/bash

packages=(
    "./cmd"
    "./internal/service/fighters"
    "./internal/handler/grpc"
    "./internal/controller/fighters"
    "./internal/repository/psql"
    "./pkg/cfg"
    "./pkg/errors"
    "./pkg/model"
)

cmd="go test -v -coverprofile=coverage.out"

if [ $# -eq 1 ]; then
  cmd+=" $1"
else
  for pkg in "${packages[@]}"; do
    cmd+=" $pkg"
done
fi

eval "$cmd"