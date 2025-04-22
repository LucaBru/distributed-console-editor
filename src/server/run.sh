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

go run . --raft_bootstrap --raft_id=nodeA --address=https://server1.enricozangrando.com --raft_data_dir /tmp/my-raft-cluster &
sleep 2 &&
go run . --raft_id=nodeB --https://server2.enricozangrando.com --raft_data_dir /tmp/my-raft-cluster &
sleep 2 &&
go run . --raft_id=nodeC --address=https://server3.enricozangrando.com --raft_data_dir /tmp/my-raft-cluster &
sleep 2

go install github.com/Jille/raftadmin/cmd/raftadmin@latest
echo -e "\nAdding nodes B and C to the cluster" 
raftadmin https://server1.enricozangrando.com add_voter nodeB 127.0.0.1:50052 0
raftadmin --leader multi:///https://server1.enricozangrando.com,127.0.0.1:50052 add_voter nodeC https://server3.enricozangrando.com 0
sleep 2

echo -e "\nCluster is online 🚀🚀"
raftadmin https://server1.enricozangrando.com leader
raftadmin https://server1.enricozangrando.com get_configuration


wait