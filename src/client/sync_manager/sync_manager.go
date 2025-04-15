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

var Writer = initWriter()

func initWriter() *bufio.Writer {
	file, err := os.Create("sync_man.log")
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return nil
	}
	return bufio.NewWriter(file)
}

type SyncManager struct {
	sync.Mutex
	connection     *grpc.ClientConn
	node           editorpb.NodeClient
	DocConfig      DocumentConfig
	notifyOnUpdate chan<- struct{}
	buff           ot.Ops
	wait           ot.Ops
}

func NewSyncManager(docConfig DocumentConfig) (*SyncManager, <-chan struct{}) {
	notifyUpdate := make(chan struct{}, 20)
	syncManager := &SyncManager{
		DocConfig:      docConfig,
		connection:     initConnection(),
		node:           editorpb.NewNodeClient(initConnection()),
		notifyOnUpdate: notifyUpdate,
	}
	go syncManager.startUpdateListener()
	return syncManager, notifyUpdate
}

func (syncManager *SyncManager) ApplyEdit(ops ot.Ops) {
	syncManager.Lock()
	defer syncManager.Unlock()
	fmt.Fprintln(Writer, "Applying edit %v %s", ops, syncManager.DocConfig.Document)
	Writer.Flush()
	var err error
	if err = syncManager.DocConfig.Document.Apply(ops); err != nil {
		fmt.Fprintf(Writer, "Error in edit application: %s", err.Error())
		return
	}
	fmt.Fprintf(Writer, "Edit correctly applied, doc %s\n", string(syncManager.DocConfig.Document))
	switch {
	case syncManager.buff != nil:
		{
			if syncManager.buff, err = ot.Compose(syncManager.buff, ops); err != nil {
				fmt.Fprintf(Writer, "Error in edit composition: %s", err.Error())
			}
		}
	case syncManager.wait != nil:
		{
			syncManager.buff = ops
		}
	default:
		syncManager.wait = ops
		go syncManager.Send(ops, syncManager.DocConfig.version)
	}
	Writer.Flush()
}

func (syncManager *SyncManager) Send(ops ot.Ops, rev int) {
	syncManager.Lock()
	defer syncManager.Unlock()
	_, err := syncManager.node.Edit(context.Background(), &editorpb.EditReq{DocId: syncManager.DocConfig.DocId, Rev: int32(rev), Ops: ops.WireFmt(), UserId: syncManager.DocConfig.authorId, Title: syncManager.DocConfig.title})
	if err != nil {
		fmt.Fprintln(Writer, "Error sending operations:", err)
		return
	}
	go syncManager.Ack()
}

func (syncManager *SyncManager) Ack() {
	syncManager.Lock()
	defer syncManager.Unlock()
	switch {
	case syncManager.buff != nil:
		go syncManager.Send(syncManager.buff, syncManager.DocConfig.version+1)
		syncManager.wait, syncManager.buff = syncManager.buff, nil
	case syncManager.wait != nil:
		syncManager.wait = nil
	default:
		fmt.Fprintln(Writer, "Error sending operations")
		return
	}
	syncManager.DocConfig.version++
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
	syncManager.Lock()
	stream, err := syncManager.node.WatchDocument(context.Background())
	if err != nil {
		fmt.Fprintln(Writer, "Watch error: "+err.Error())
	}
	joinDoc := true
	if syncManager.DocConfig.DocId == "" {
		syncManager.shareDoc(Writer)
		joinDoc = false
	}
	stream.Send(&editorpb.WatchReq{DocId: syncManager.DocConfig.DocId, UserId: syncManager.DocConfig.authorId})
	docSnapshot, err := stream.Recv()
	if err != nil {
		fmt.Fprintln(Writer, "Doc snapshot error: "+err.Error())
	}

	if joinDoc {
		syncManager.DocConfig.title = docSnapshot.Title
		fmt.Fprintf(Writer, "File received: ", string(docSnapshot.Doc))
		syncManager.DocConfig.Document = docSnapshot.Doc
		syncManager.DocConfig.version = int(docSnapshot.Rev)
		fmt.Fprintln(Writer, "Document snapshot recv", docSnapshot.Rev)
		syncManager.notifyOnUpdate <- struct{}{}
		Writer.Flush()
	}
	syncManager.Unlock()

	for {
		updatedData, err := stream.Recv()
		fmt.Fprintln(Writer, "Stream Recv called")
		Writer.Flush()
		if err == io.EOF {
			fmt.Fprintln(Writer, "Stream closed")
			Writer.Flush()
			continue
		}
		if err != nil {
			fmt.Fprintf(Writer, "Error receiving data: %v", err)
			Writer.Flush()
			continue
		}
		syncManager.Lock()
		err = syncManager.DocConfig.Document.Apply(ot.NewOps(updatedData.Ops))
		if err != nil {
			fmt.Fprintln(Writer, "Failed to apply received ops")
			syncManager.Unlock()
			return
		}
		fmt.Fprintf(Writer, "Increasing version due to edit received")
		syncManager.DocConfig.version++
		syncManager.notifyOnUpdate <- struct{}{}
		syncManager.Unlock()
		Writer.Flush()
	}
}

func (syncManager *SyncManager) shareDoc(Writer *bufio.Writer) {
	fmt.Fprintln(Writer, "DocId is empty, creating new document...")
	msg, err := syncManager.node.Share(context.Background(), &editorpb.ShareReq{DocName: syncManager.DocConfig.title, Doc: syncManager.DocConfig.Document, UserId: syncManager.DocConfig.authorId})
	if err != nil {
		fmt.Fprintln(Writer, "Share error "+err.Error())
	}
	syncManager.DocConfig.DocId = msg.DocId
	fmt.Fprintln(Writer, "Document created with ID: "+syncManager.DocConfig.DocId)
	Writer.Flush()
}
