#!/bin/bash

if [ $# -eq 0 ]; then
  set -- "func"
elif [ $# -ne 1 ]; then
  echo "Using: $0 <cover type>"
  exit 1
fi

cover_type="$1"

if [ "$cover_type" != "func" ] && [ "$cover_type" != "html" ]; then
  echo "Unexpected cover type, use 'func' or 'html'."
  exit 1
fi

go tool cover "-$cover_type" coverage.out