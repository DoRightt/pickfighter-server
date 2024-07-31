#!/bin/bash

MOCKS_DIR="./gen/mocks"

mkdir -p "$MOCKS_DIR"

declare -A INTERFACES
INTERFACES=(
    ["internal/controller/fighters"]="controller"
    ["internal/repository/psql"]="db"
    ["internal/handler/grpc"]="grpc"
)

for PACKAGE in "${!INTERFACES[@]}"; do
    INTERFACE="${INTERFACES[$PACKAGE]}"
    DESTINATION="$MOCKS_DIR/mock_$(basename "$PACKAGE").go"
    
    echo "Generating mock for $INTERFACE in package $PACKAGE..."
    
    mockgen -source="$PACKAGE/${INTERFACE}.go" -destination="$DESTINATION" -package=mocks
    
    if [ $? -ne 0 ]; then
        echo "Error generating mock for $INTERFACE in package $PACKAGE"
        exit 1
    fi
done

echo "Mocks generated successfully in $MOCKS_DIR"