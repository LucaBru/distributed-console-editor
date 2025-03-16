package main

import (
	"context"
	pb "editor-service/protos"
	service "editor-service/src"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	leaderHealth "github.com/Jille/raft-grpc-leader-rpc/leaderhealth"
	transport "github.com/Jille/raft-grpc-transport"
	raftAdmin "github.com/Jille/raftadmin"
	"github.com/hashicorp/raft"
	boltDb "github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	myAddr = flag.String("address", "localhost:50051", "TCP host+port for this node")
	raftId = flag.String("raft_id", "", "Node id used by Raft")

	raftDir       = flag.String("raft_data_dir", "data/", "Raft data dir")
	raftBootstrap = flag.Bool("raft_bootstrap", false, "Whether to bootstrap the Raft cluster")
)

func main() {
	flag.Parse()

	if *raftId == "" {
		log.Fatalf("flag --raft_id is required")
	}

	ctx := context.Background()
	_, port, error := net.SplitHostPort(*myAddr)
	if error != nil {
		log.Fatalf("failed to parse local address (%q): %v", *myAddr, error)
	}
	sock, error := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if error != nil {
		log.Fatalf("failed to listen: %v", error)
	}

	wt := &service.Document{}

	raft, transportManager, error := NewRaft(ctx, *raftId, *myAddr, wt)
	if error != nil {
		log.Fatalf("failed to start raft: %v", error)
	}
	server := grpc.NewServer()
	pb.RegisterEditorServer(server, service.NewEditor(raft))
	transportManager.Register(server)
	leaderHealth.Setup(raft, server, []string{"Example"})
	raftAdmin.Register(server, raft)
	reflection.Register(server)
	if err := server.Serve(sock); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func NewRaft(ctx context.Context, myID, myAddress string, fsm raft.FSM) (*raft.Raft, *transport.Manager, error) {
	c := raft.DefaultConfig()
	c.LocalID = raft.ServerID(myID)

	baseDir := filepath.Join(*raftDir, myID)

	ldb, err := boltDb.NewBoltStore(filepath.Join(baseDir, "logs.dat"))
	if err != nil {
		return nil, nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "logs.dat"), err)
	}

	sdb, err := boltDb.NewBoltStore(filepath.Join(baseDir, "stable.dat"))
	if err != nil {
		return nil, nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "stable.dat"), err)
	}

	fileSnapshotStore, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stderr)
	if err != nil {
		return nil, nil, fmt.Errorf(`raft.NewFileSnapshotStore(%q, ...): %v`, baseDir, err)
	}

	transportManager := transport.New(raft.ServerAddress(myAddress), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})

	r, err := raft.NewRaft(c, fsm, ldb, sdb, fileSnapshotStore, transportManager.Transport())
	if err != nil {
		return nil, nil, fmt.Errorf("raft.NewRaft: %v", err)
	}

	if *raftBootstrap {
		cfg := raft.Configuration{
			Servers: []raft.Server{
				{
					Suffrage: raft.Voter,
					ID:       raft.ServerID(myID),
					Address:  raft.ServerAddress(myAddress),
				},
			},
		}
		f := r.BootstrapCluster(cfg)
		if err := f.Error(); err != nil {
			return nil, nil, fmt.Errorf("raft.Raft.BootstrapCluster: %v", err)
		}
	}

	return r, transportManager, nil
}
