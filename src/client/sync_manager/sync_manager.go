package sync_manager

import (
	"bufio"
	"context"
	"editor-service/node/ot"
	"editor-service/protos/editorpb"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	_ "github.com/Jille/grpc-multi-resolver"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/health"
)

type SyncManager struct {
	sync.Mutex
	connection     *grpc.ClientConn
	node           editorpb.NodeClient
	DocConfig      DocumentConfig
	notifyOnUpdate chan<- struct{}
}

func NewSyncManager(docConfig DocumentConfig) (*SyncManager, <-chan struct{}) {
	file, err := os.Create("client.log")
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return nil, nil
	}
	writer := bufio.NewWriter(file)

	notifyUpdate := make(chan struct{}, 20)
	syncManager := &SyncManager{
		DocConfig:      docConfig,
		connection:     initConnection(),
		node:           editorpb.NewNodeClient(initConnection()),
		notifyOnUpdate: notifyUpdate,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go syncManager.startUpdateListener(writer, wg)
	wg.Wait()
	return syncManager, notifyUpdate
}

func (syncManager *SyncManager) SendData(operations []*editorpb.Op) {
	if len(operations) == 0 {
		return
	}

	file, err := os.Create("sync.log")
	if err != nil {
		fmt.Println("Error creating log file:", err)
	}
	writer := bufio.NewWriter(file)
	fmt.Fprintln(writer, "Docid: "+syncManager.DocConfig.DocId)
	fmt.Fprintln(writer, "Sending operations:", operations, syncManager.DocConfig.version)
	writer.Flush()

	syncManager.Lock()
	defer syncManager.Unlock()
	errA := syncManager.DocConfig.Document.Apply(ot.NewOps(operations))
	if errA != nil {
		fmt.Fprintln(writer, "Error applying operations:", errA)
	}
	ack, err2 := syncManager.node.Edit(context.Background(), &editorpb.EditReq{DocId: syncManager.DocConfig.DocId, Rev: int32(syncManager.DocConfig.version) /* int32(syncManager.DocConfig.version) */, Ops: operations, UserId: syncManager.DocConfig.authorId, Title: syncManager.DocConfig.title})
	if err2 != nil {
		fmt.Fprintln(writer, "Error sending operations:", err2)
	}
	fmt.Fprintln(writer, "Increasing version due to edit the local file")
	syncManager.DocConfig.version++
	fmt.Fprintln(writer, "Received ack:", ack)
	writer.Flush()
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

func (syncManager *SyncManager) startUpdateListener(writer *bufio.Writer, wg *sync.WaitGroup) {
	syncManager.Lock()
	stream, err := syncManager.node.WatchDocument(context.Background())
	if err != nil {
		fmt.Fprintln(writer, "Watch error: "+err.Error())
	}
	joinDoc := true
	if syncManager.DocConfig.DocId == "" {
		syncManager.shareDoc(writer)
		joinDoc = false
	}
	stream.Send(&editorpb.WatchReq{DocId: syncManager.DocConfig.DocId, UserId: syncManager.DocConfig.authorId})
	docSnapshot, err := stream.Recv()
	if err != nil {
		fmt.Fprintln(writer, "Doc snapshot error: "+err.Error())
	}

	if joinDoc {
		syncManager.DocConfig.title = docSnapshot.Title
		fmt.Fprintf(writer, "File received: ", string(docSnapshot.Doc))
		syncManager.DocConfig.Document = docSnapshot.Doc
		syncManager.DocConfig.version = int(docSnapshot.Rev)
		fmt.Fprintln(writer, "Document snapshot recv", docSnapshot.Rev)
		writer.Flush()
	}
	wg.Done()
	syncManager.Unlock()

	for {
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
		syncManager.Lock()
		err = syncManager.DocConfig.Document.Apply(ot.NewOps(updatedData.Ops))
		if err != nil {
			fmt.Fprintln(writer, "Failed to apply received ops")
			syncManager.Unlock()
			return
		}
		fmt.Fprintf(writer, "Increasing version due to edit received")
		syncManager.DocConfig.version++
		syncManager.notifyOnUpdate <- struct{}{}
		syncManager.Unlock()
		writer.Flush()
	}
}

func (syncManager *SyncManager) shareDoc(writer *bufio.Writer) {
	fmt.Fprintln(writer, "DocId is empty, creating new document...")
	msg, err := syncManager.node.Share(context.Background(), &editorpb.ShareReq{DocName: syncManager.DocConfig.title, Doc: syncManager.DocConfig.Document, UserId: syncManager.DocConfig.authorId})
	if err != nil {
		fmt.Fprintln(writer, "Share error "+err.Error())
	}
	syncManager.DocConfig.DocId = msg.DocId
	fmt.Fprintln(writer, "Document created with ID: "+syncManager.DocConfig.DocId)
	writer.Flush()
}
