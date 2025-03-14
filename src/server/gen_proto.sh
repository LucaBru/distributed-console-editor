#!/bin/bash

cd ../../
protoc --go_out=src/server/ --go_opt=paths=source_relative \
    --go-grpc_out=src/server --go-grpc_opt=paths=source_relative \
    protos/editor.proto