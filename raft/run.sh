#!/bin/bash

PORTS=(50051 50052 50053) 
for PORT in "${PORTS[@]}"; do
    PID=$(lsof -i :$PORT -t 2>/dev/null)
    if [ -n "$PID" ]; then
        kill -9 $PID
    fi
done

rm -rf /tmp/my-raft-cluster
mkdir /tmp/my-raft-cluster
mkdir /tmp/my-raft-cluster/node{A,B,C}
go run . --raft_bootstrap --raft_id=nodeA --address=localhost:50051 --raft_data_dir /tmp/my-raft-cluster &
go run . --raft_id=nodeB --address=localhost:50052 --raft_data_dir /tmp/my-raft-cluster &
go run . --raft_id=nodeC --address=localhost:50053 --raft_data_dir /tmp/my-raft-cluster &
sleep 2

go install github.com/Jille/raftadmin/cmd/raftadmin@latest
echo -e "\nAdding nodes B and C to the cluster" 
raftadmin localhost:50051 add_voter nodeB localhost:50052 0
raftadmin --leader multi:///localhost:50051,localhost:50052 add_voter nodeC localhost:50053 0

echo -e "\nRequiring a leadership transfer"
raftadmin localhost:50051 leader
raftadmin --leader multi:///localhost:50051,localhost:50052,localhost:50053 leadership_transfer
echo -e "\nNew leader"
raftadmin localhost:50051 leader

echo -e "\nRunning the client"
go run cmd/hammer/hammer.go