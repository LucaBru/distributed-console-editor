package sync_manager

import (
	"bufio"
	"context"
	"editor-service/node/ot"
	"editor-service/protos/editorpb"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/Jille/grpc-multi-resolver"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/health"
	"google.golang.org/grpc/status"
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
	docConfig      DocumentConfig
	notifyOnUpdate chan<- struct{}
	buf            ot.Ops
	wait           ot.Ops
}

func NewSyncManager(docConfig DocumentConfig) (*SyncManager, <-chan struct{}) {
	notifyUpdate := make(chan struct{}, 20)
	syncManager := &SyncManager{
		docConfig:      docConfig,
		connection:     initConnection(),
		node:           editorpb.NewNodeClient(initConnection()),
		notifyOnUpdate: notifyUpdate,
	}
	go syncManager.startUpdateListener()
	return syncManager, notifyUpdate
}

func (m *SyncManager) GetDoc() ot.Doc {
	m.Lock()
	defer m.Unlock()
	return append([]byte{}, m.docConfig.Document...)
}

func (m *SyncManager) GetDocId() string {
	m.Lock()
	defer m.Unlock()
	return m.docConfig.DocId
}

func (syncManager *SyncManager) ApplyEdit(ops ot.Ops) {
	syncManager.Lock()
	defer syncManager.Unlock()
	var err error
	if err = syncManager.docConfig.Document.Apply(ops); err != nil {
		fmt.Fprintf(Writer, "Error in edit application: %s", err.Error())
		return
	}
	switch {
	case syncManager.buf != nil:
		{
			if syncManager.buf, err = ot.Compose(syncManager.buf, ops); err != nil {
				fmt.Fprintf(Writer, "Error in edit composition: %s", err.Error())
			}
		}
	case syncManager.wait != nil:
		{
			syncManager.buf = ops
		}
	default:
		syncManager.wait = ops
		go syncManager.Send(ops, syncManager.docConfig.version)
	}
	Writer.Flush()
}

func (syncManager *SyncManager) Send(ops ot.Ops, rev int) {
	syncManager.Lock()
	defer syncManager.Unlock()
	_, err := syncManager.node.Edit(context.Background(), &editorpb.EditReq{DocId: syncManager.docConfig.DocId, Rev: int32(rev), Ops: ops.WireFmt(), UserId: syncManager.docConfig.authorId, Title: syncManager.docConfig.title})
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
	case syncManager.buf != nil:
		go syncManager.Send(syncManager.buf, syncManager.docConfig.version+1)
		syncManager.wait, syncManager.buf = syncManager.buf, nil
	case syncManager.wait != nil:
		syncManager.wait = nil
	default:
		fmt.Fprintln(Writer, "Error sending operations")
		return
	}
	syncManager.docConfig.version++
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

func (m *SyncManager) startUpdateListener() {
	errorCounter := 0
	m.Lock()
	stream, err := m.node.WatchDocument(context.Background())
	if err != nil {
		fmt.Fprintln(Writer, "Watch error: "+err.Error())
	}
	joinDoc := true
	if m.docConfig.DocId == "" {
		m.shareDoc(Writer)
		joinDoc = false
	}
	stream.Send(&editorpb.WatchReq{DocId: m.docConfig.DocId, UserId: m.docConfig.authorId})
	docSnapshot, err := stream.Recv()
	
	if err != nil {
		fmt.Fprintln(Writer, "Doc snapshot error: "+err.Error())
	}
	if joinDoc {
		m.docConfig.title = docSnapshot.Title
		fmt.Fprintf(Writer, "File received: ", string(docSnapshot.Doc))
		m.docConfig.Document = docSnapshot.Doc
		m.docConfig.version = int(docSnapshot.Rev)
		fmt.Fprintln(Writer, "Document snapshot recv", docSnapshot.Rev)
		m.notifyOnUpdate <- struct{}{}
		Writer.Flush()
	}
	m.Unlock()

	defer Writer.Flush()
	for {
		fmt.Println("Updating...")
		update, err := stream.Recv()
		if err == io.EOF {
			fmt.Fprintln(Writer, "Stream closed")
			Writer.Flush()
			return
		}
		if err != nil {
			log.Println("Error receiving data: ", err)
			code := status.Code(err)
				if code != codes.OK {
					errorCounter++
					fmt.Println("Total errors: ", errorCounter)
				}
			fmt.Fprintf(Writer, "Error receiving data: %v", err)
			Writer.Flush()
			return
		}
		m.Lock()
		ops := ot.NewOps(update.Ops)
		if m.wait != nil {
			if ops, m.wait, err = ot.Transform(ops, m.wait); err != nil {
				fmt.Fprintf(Writer, "Error traforming collaborators updates: %s\n", err.Error())
				return
			}
		}
		if m.buf != nil {
			if ops, m.buf, err = ot.Transform(ops, m.buf); err != nil {
				fmt.Fprintf(Writer, "Error traforming collaborators updates: %s\n", err.Error())
				return
			}
		}
		if err = m.docConfig.Document.Apply(ops); err != nil {
			fmt.Fprintf(Writer, "Error applying collaborators updates: %s\n", err.Error())
			log.Println("Error applying collaborators updates: ", err)
			return
		}
		m.docConfig.version++
		m.notifyOnUpdate <- struct{}{}
		m.Unlock()
	}
}

func (syncManager *SyncManager) shareDoc(Writer *bufio.Writer) {
	fmt.Fprintln(Writer, "DocId is empty, creating new document...")
	msg, err := syncManager.node.Share(context.Background(), &editorpb.ShareReq{DocName: syncManager.docConfig.title, Doc: syncManager.docConfig.Document, UserId: syncManager.docConfig.authorId})
	if err != nil {
		fmt.Fprintln(Writer, "Share error "+err.Error())
	}
	syncManager.docConfig.DocId = msg.DocId
	fmt.Fprintln(Writer, "Document created with ID: "+syncManager.docConfig.DocId)
	Writer.Flush()
}
