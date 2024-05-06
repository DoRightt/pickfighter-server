#!/bin/bash

packages=(
    "./internal/services"
    "./internal/services/auth"
    "./internal/services/common"
    "./internal/repo/auth"
    "./internal/repo/common"
    "./internal/repo/fighters"
    "./cmd"
    "./pkg/cfg"
    "./pkg/errors"
    "./pkg/httplib"
    "./pkg/logger"
    "./pkg/model"
    "./pkg/pgxs"
    "./pkg/sigx"
    "./pkg/utils"
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