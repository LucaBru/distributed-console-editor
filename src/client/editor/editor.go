package editor

import (
	"context"
	"editor-service/node/ot"
	"editor-service/protos/editorpb"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	_ "github.com/Jille/grpc-multi-resolver"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/nsf/termbox-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/health"
)

type Editor struct {
	sync.Mutex
	doc   *ot.Doc // the document
	docId string
	Rev   int    // last acknowledged revision
	Wait  ot.Ops // pending ops or nil
	Buf   ot.Ops // buffered ops or nil
	// Send is called when a new revision can be sent to the server.
	Send      func(int, ot.Ops)
	observers []chan<- ot.Doc
}

func NewEditor(docId string, authorId string) *Editor {
	cEditor := Editor{
		doc:   &ot.Doc{},
		docId: docId,
	}
	node := editorpb.NewNodeClient(initConn())
	send := func(rev int, ops ot.Ops) {
		stopWatch := make(chan time.Time, 1)
		stopWatch <- time.Now()
		log.Print("send edit to cluster")
		req := &editorpb.EditReq{DocId: docId, Rev: int32(rev), Ops: ops.WireFmt(), UserId: authorId, Title: "random"}
		_, err := node.Edit(context.Background(), req)
		log.Print("cluster answers to doc edit")
		if err != nil {
			log.Fatal("failed to edit remote document: ", err)
			return
		}
		log.Print("successful edit")
		cEditor.Ack(stopWatch)
	}
	cEditor.Send = send
	wait := make(chan struct{})
	go cEditor.Recv(node, authorId, wait)
	<-wait
	log.Print("returned")
	return &cEditor
}

func initConn() *grpc.ClientConn {
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

func (c *Editor) Observe() <-chan ot.Doc {
	c.Lock()
	defer c.Unlock()
	ch := make(chan ot.Doc)
	c.observers = append(c.observers, ch)
	return ch
}

func (c *Editor) notifyObservers() {
	c.Lock()
	defer c.Unlock()
	for _, obs := range c.observers {
		obs <- *c.doc
	}
}

func (c *Editor) apply(ops ot.Ops) error {
	log.Print("apply edit to local doc")
	var err error
	if err = c.doc.Apply(ops); err != nil {
		return err
	}
	go c.notifyObservers()
	switch {
	case c.Buf != nil:
		if c.Buf, err = ot.Compose(c.Buf, ops); err != nil {
			return err
		}
	case c.Wait != nil:
		c.Buf = ops
	default:
		c.Wait = ops
		rev := c.Rev
		go c.Send(rev, ops)
	}
	return nil
}

// Ack acknowledges a pending server update and sends buffered updates if any.
// An error is returned if no update is pending.
func (c *Editor) Ack(stopWatch <-chan time.Time) {
	log.Print("ack")
	c.Lock()
	defer c.Unlock()
	switch {
	case c.Buf != nil:
		go c.Send(c.Rev+1, c.Buf)
		c.Wait, c.Buf = c.Buf, nil
	case c.Wait != nil:
		c.Wait = nil
	default:
		log.Fatal("no pending operation")
	}
	c.Rev++
	log.Print("successful ack")
	log.Print("edited remote document in ", time.Since(<-stopWatch))
}

// Recv receives server updates originating from other participants.
// An error is returned if the server update could not be applied.
func (c *Editor) Recv(node editorpb.NodeClient, author string, wait chan<- struct{}) {
	var err error
	stream, err := node.WatchDocument(context.Background())
	if err != nil {
		log.Fatal("failed to watch remote document: ", err)
	}
	stream.Send(&editorpb.WatchReq{DocId: c.docId, UserId: author})
	snap, err := stream.Recv()
	log.Print("snapshot received")
	if err != nil {
		log.Fatal("failed to get document snapshot: ", err)
	}
	c.doc = (*ot.Doc)(&snap.Doc)
	c.Rev = int(snap.Rev)
	close(wait)

	for {
		edit, err := stream.Recv()
		log.Println("received remote update")
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("failed to receive a remote edit: ", err)
		}
		ops := ot.NewOps(edit.Ops)

		c.Lock()
		if c.Wait != nil {
			if ops, c.Wait, err = ot.Transform(ops, c.Wait); err != nil {
				log.Fatal(err)
			}
		}
		if c.Buf != nil {
			if ops, c.Buf, err = ot.Transform(ops, c.Buf); err != nil {
				log.Fatal(err)
			}
		}
		if err = c.doc.Apply(ops); err != nil {
			log.Fatal(err)
		}
		c.Rev++
		c.Unlock()
		for _, obs := range c.observers {
			obs <- *c.doc
		}
	}
}

func (e *Editor) Edit(event termbox.Event) {
	e.Lock()
	defer e.Unlock()
	switch event.Key {
	default:
		{
			// insert rune always at pos 0 for the moment
			ops := ot.Ops{ot.Op{S: string(event.Ch)}, ot.Op{N: int(len(*e.doc))}}
			err := e.apply(ops)
			if err != nil {
				log.Fatal("local edit failed: ", err)
			}
			log.Print("successful local edit")
		}
	}
}
