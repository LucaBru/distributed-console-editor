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

clear

go run . --raft_bootstrap --raft_id=nodeA --address=127.0.0.1:50051 --raft_data_dir /tmp/my-raft-cluster &
go run . --raft_id=nodeB --address=127.0.0.1:50052 --raft_data_dir /tmp/my-raft-cluster &
go run . --raft_id=nodeC --address=127.0.0.1:50053 --raft_data_dir /tmp/my-raft-cluster &
sleep 2

go install github.com/Jille/raftadmin/cmd/raftadmin@latest
echo -e "\nAdding nodes B and C to the cluster" 
raftadmin 127.0.0.1:50051 add_voter nodeB 127.0.0.1:50052 0
raftadmin --leader multi:///127.0.0.1:50051,127.0.0.1:50052 add_voter nodeC 127.0.0.1:50053 0
sleep 2

echo -e "\nCluster is online 🚀🚀"
raftadmin 127.0.0.1:50051 leader
raftadmin 127.0.0.1:50051 get_configuration


wait