#!/bin/bash

cd ../../
export PREFIX="editor-service"
protoc --go_out=src/server/ --go_opt=module=$PREFIX \
    --go-grpc_out=src/server/  --go-grpc_opt=module=$PREFIX \
    protos/editor.proto

protoc --go_out=src/server/  --go_opt=module=$PREFIX \
    --go-grpc_out=src/server/  --go-grpc_opt=module=$PREFIX \
    src/server/node/log.proto
