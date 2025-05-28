#!/bin/bash
for clientId in {0..9}; do
    go run . --doc-id=38a14415-8cb3-4d70-ab9f-e4d1a6463e7f --author=$clientId &
    done
sleep 50