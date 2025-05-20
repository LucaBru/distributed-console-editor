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
mkdir /tmp/my-raft-cluster/node{A,B,C,D,E,F,G,H,I,J}

clear

go run . --raft_bootstrap --raft_id=nodeA --address=localhost:50051 --raft_data_dir /tmp/my-raft-cluster &
go run . --raft_id=nodeB --address=localhost:50052 --raft_data_dir /tmp/my-raft-cluster &
go run . --raft_id=nodeC --address=localhost:50053 --raft_data_dir /tmp/my-raft-cluster &
go run . --raft_id=nodeD --address=localhost:50054 --raft_data_dir /tmp/my-raft-cluster &
go run . --raft_id=nodeE --address=localhost:50055 --raft_data_dir /tmp/my-raft-cluster &
go run . --raft_id=nodeF --address=localhost:50056 --raft_data_dir /tmp/my-raft-cluster &
go run . --raft_id=nodeG --address=localhost:50057 --raft_data_dir /tmp/my-raft-cluster &
go run . --raft_id=nodeH --address=localhost:50058 --raft_data_dir /tmp/my-raft-cluster &
go run . --raft_id=nodeI --address=localhost:50059 --raft_data_dir /tmp/my-raft-cluster &
go run . --raft_id=nodeJ --address=localhost:50060 --raft_data_dir /tmp/my-raft-cluster &
sleep 2

go install github.com/Jille/raftadmin/cmd/raftadmin@latest
echo -e "\nAdding nodes B and C to the cluster" 
raftadmin localhost:50051 add_voter nodeB localhost:50052 0
raftadmin --leader multi:///localhost:50051,localhost:50052 add_voter nodeC localhost:50053 0
raftadmin --leader multi:///localhost:50051,localhost:50052,localhost:50053 add_voter nodeD localhost:50054 0
raftadmin --leader multi:///localhost:50051,localhost:50052,localhost:50053,localhost:50054 add_voter nodeE localhost:50055 0
raftadmin --leader multi:///localhost:50051,localhost:50052,localhost:50053,localhost:50054,localhost:50055 add_voter nodeF localhost:50056 0
raftadmin --leader multi:///localhost:50051,localhost:50052,localhost:50053,localhost:50054,localhost:50055,localhost:50056 add_voter nodeG localhost:50057 0
raftadmin --leader multi:///localhost:50051,localhost:50052,localhost:50053,localhost:50054,localhost:50055,localhost:50056,localhost:50057 add_voter nodeH localhost:50058 0
raftadmin --leader multi:///localhost:50051,localhost:50052,localhost:50053,localhost:50054,localhost:50055,localhost:50056,localhost:50057,localhost:50058 add_voter nodeI localhost:50059 0
raftadmin --leader multi:///localhost:50051,localhost:50052,localhost:50053,localhost:50054,localhost:50055,localhost:50056,localhost:50057,localhost:50058,localhost:50059 add_voter nodeJ localhost:50060 0
sleep 2

echo -e "\nCluster is online 🚀🚀"
raftadmin localhost:50051 leader
raftadmin localhost:50051 get_configuration


wait