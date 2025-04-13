package sync_manager

import (
	"bufio"
	"context"
	"editor-service/node/ot"
	"editor-service/protos/editorpb"
	"fmt"
	"io"
	"os"
	"time"

	_ "github.com/Jille/grpc-multi-resolver"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/health"
)

type SyncManager struct {
	connection *grpc.ClientConn
	node       editorpb.NodeClient
	DocConfig  DocumentConfig
}

func NewSyncManager(docConfig DocumentConfig) *SyncManager {
	file, err := os.Create("client.log")
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return nil
	}
	writer := bufio.NewWriter(file)

	syncManager := &SyncManager{
		DocConfig:  docConfig,
		connection: initConnection(),
		node:       editorpb.NewNodeClient(initConnection()),
	}
	go syncManager.startUpdateListener(writer)
	return syncManager
}

func (syncManager *SyncManager) SendData(operations []*editorpb.Op) {
	if len(operations) != 0 {
		file, err := os.Create("sync.log")
		if err != nil {
			fmt.Println("Error creating log file:", err)
		}
		writer := bufio.NewWriter(file)
		fmt.Fprintln(writer, "Docid: "+syncManager.DocConfig.DocId)
		fmt.Fprintln(writer, "Sending operations:", operations)
		writer.Flush()
		errA := syncManager.DocConfig.Document.Apply(ot.NewOps(operations))
		if errA != nil {
			fmt.Fprintln(writer, "Error applying operations:", errA)
		}
		ack, err2 := syncManager.node.Edit(context.Background(), &editorpb.EditReq{DocId: syncManager.DocConfig.DocId, Rev: 0 /* int32(syncManager.DocConfig.version) */, Ops: operations, UserId: syncManager.DocConfig.authorId, Title: syncManager.DocConfig.title})
		if err2 != nil {
			fmt.Fprintln(writer, "Error sending operations:", err2)
		}
		fmt.Fprintln(writer, "Received ack:", ack)
		writer.Flush()
	}
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

func (syncManager *SyncManager) startUpdateListener(writer *bufio.Writer) {
	if syncManager.DocConfig.DocId == "" {
		fmt.Fprintln(writer, "DocId is empty, creating new document...")
		msg, err := syncManager.node.Share(context.Background(), &editorpb.ShareReq{DocName: syncManager.DocConfig.title, Doc: syncManager.DocConfig.Document, UserId: syncManager.DocConfig.authorId})
		if err != nil {
			fmt.Fprintln(writer, "Share error "+err.Error())
		}
		syncManager.DocConfig.DocId = msg.DocId
		fmt.Fprintln(writer, "Document created with ID: "+syncManager.DocConfig.DocId)
		writer.Flush()
	}
	stream, err := syncManager.node.WatchDocument(context.Background())
	if err != nil {
		fmt.Fprintln(writer, "Watch error: "+err.Error())
	}
	stream.Send(&editorpb.WatchReq{DocId: syncManager.DocConfig.DocId, UserId: syncManager.DocConfig.authorId})
	docSnapshot, err := stream.Recv()
	if err != nil {
		fmt.Fprintln(writer, "Doc snapshot error: "+err.Error())
	}
	writer.Flush()
	syncManager.DocConfig.title = docSnapshot.Title
	syncManager.DocConfig.Document = docSnapshot.Doc
	syncManager.DocConfig.version = int(docSnapshot.Rev)
	fmt.Fprintln(writer, "Document snapshot received with rev: ", docSnapshot.Rev)
	writer.Flush()
	for true {
		time.Sleep(2 * time.Second)
		fmt.Fprintln(writer, "Waiting for updates...")
		writer.Flush()

		updatedData, err := stream.Recv()
		fmt.Fprintln(writer, "Stream Recv called")
		writer.Flush()
		if err == io.EOF {
			fmt.Fprintln(writer, "Stream closed")
			writer.Flush()
			continue
		}
		if err != nil {
			fmt.Fprintf(writer, "Error receiving data: %v", err)
			writer.Flush()
			continue
		}
		syncManager.DocConfig.Document.Apply(ot.NewOps(updatedData.Ops))
		fmt.Fprintf(writer, "Received update: %s\n", string(updatedData.Doc))
		writer.Flush()
	}
}
