#!/bin/bash

# This section handles command-line arguments using a while loop.
# It parses options preceded by -p or --path to specify the path where the .proto files are located.
# It stores the path in the variable PROTO_PATH.
POSITIONAL=()
while [[ $# -gt 0 ]]
do
KEY="$1"

case $KEY in
    -p|--path)
    PROTO_PATH="$2"
    shift # past argument
    shift # past value
    ;;
esac
done
set -- "${POSITIONAL[@]}" # restore positional parameters

# This section checks if the PROTO_PATH variable is set.
# If not, it outputs an error message and exits the script with a non-zero status code.
if [[ ! $PROTO_PATH ]]; then
    echo "Please specify proto folder with -p option"
    exit 1
fi

# Run protoc to generate Go files from your .proto files
protoc -I proto/ \
--go_out=./internal/proto --go_opt=paths=source_relative \
--go-grpc_out=./internal/proto --go-grpc_opt=paths=source_relative \
proto/$PROTO_PATH/*.proto

# This part finds all generated _grpc.pb.go files in the output directory and iterates over them using a while loop.
# For each file, it extracts the directory and filename, constructs a destination path for the mock file, and generates mocks using mockgen.
find ./internal/proto/$PROTO_PATH -maxdepth 1 -name "*_grpc.pb.go" | while read -r file; do \
  dir=$(dirname "$file"); \
  base=$(basename "$file" _grpc.pb.go); \
  dest="./test/mocks/proto/$PROTO_PATH/mock_${base}.go"; \
  mockgen -source="$file" -destination="$dest"; \
done