package sync_manager

import (
	"context"
	"editor-service/node/ot"
	"editor-service/protos/editorpb"
	"fmt"
	"io"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SyncManager struct {
	connection *grpc.ClientConn
	node       editorpb.NodeClient
	docConfig  DocumentConfig
}

func NewSyncManager(docConfig DocumentConfig) *SyncManager {
	syncManager := &SyncManager{
		docConfig:  docConfig,
		connection: initConnection(),
		node:       editorpb.NewNodeClient(initConnection()),
	}
	go syncManager.startUpdateListener()
	return syncManager
}

func (syncManager *SyncManager) SendData(operations []*editorpb.Op) {
	syncManager.docConfig.document.Apply(ot.NewOps(operations))
	syncManager.node.Edit(context.Background(), &editorpb.EditReq{DocId: syncManager.docConfig.docId, Rev: int32(syncManager.docConfig.version), Ops: operations, UserId: syncManager.docConfig.authorId, Title: syncManager.docConfig.title})
}

func initConnection() *grpc.ClientConn {
	serviceConfig := `{"healthCheckConfig": {"serviceName": "Example"}, "loadBalancingConfig": [ { "round_robin": {} } ]}`
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Millisecond)),
		grpc_retry.WithMax(5),
	}
	conn, err := grpc.NewClient(
		"multi:///localhost:50051,localhost:50052,localhost:50053",
		grpc.WithDefaultServiceConfig(serviceConfig),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
	)
	if err != nil {
		fmt.Errorf("Failed to connect: %v", err)
	}
	return conn
}

func (syncManager *SyncManager) startUpdateListener() {
	stream, _ := syncManager.node.WatchDocument(context.Background())
	stream.Send(&editorpb.WatchReq{DocId: syncManager.docConfig.docId, UserId: syncManager.docConfig.authorId})
	docSnapshot, _ := stream.Recv()
	syncManager.docConfig.title = docSnapshot.Title
	syncManager.docConfig.document = docSnapshot.Doc
	syncManager.docConfig.version = int(docSnapshot.Rev)
	for {
		updatedData, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Errorf("Error receiving data: %v", err)
			return
		}
		syncManager.docConfig.document.Apply(ot.NewOps(updatedData.Ops))
	}
}
