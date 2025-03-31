// Binary hammer sends requests to your Raft cluster as fast as it can.
// It sends the written out version of the Dutch numbers up to 2000.
// In the end it asks the Raft cluster what the longest three words were.
package main

import (
	"context"
	"editor-service/node/ot"
	"editor-service/protos/editorpb"
	"fmt"
	"io"
	"sync"
	"time"

	_ "github.com/Jille/grpc-multi-resolver"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/health"
)

type Author struct {
	name   string
	doc    ot.Doc
	docId  string
	rev    int
	title  string
	client editorpb.NodeClient
}

func (a *Author) listenUpdate(wg *sync.WaitGroup, running *sync.WaitGroup) {
	stream, _ := a.client.WatchDocument(context.Background())
	fmt.Printf("%s is watching the doc\n", a.name)
	defer stream.CloseSend()
	stream.Send(&editorpb.WatchReq{DocId: a.docId, UserId: a.name})
	docSnapshot, _ := stream.Recv()
	a.title = docSnapshot.Title
	a.doc = docSnapshot.Doc
	a.rev = int(docSnapshot.Rev)
	running.Done()
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// read done.
			wg.Done()
			return
		}
		if err != nil {
			fmt.Errorf("Failed to receive a note : %v", err)
			wg.Done()
			return
		}
		a.doc.Apply(ot.NewOps(in.Ops))
		if in.Title != "" {
			fmt.Printf("%s has to update his doc title to %s\n", a.name, in.Title)
			a.title = in.Title
		}
		fmt.Printf("%s got message %v, her local doc is updated to: %s\n", a.name, in.Ops, a.doc)
	}
}

func main() {
	// creates alice author
	conn := clientConn()
	alice := &Author{name: "Alice", doc: []byte("Hello"), rev: 0, title: "Thesi", client: editorpb.NewNodeClient(conn)}

	// alice share new doc
	msg, err := alice.client.Share(context.Background(), &editorpb.ShareReq{DocName: alice.title, Doc: alice.doc, UserId: alice.name})
	if err != nil {
		fmt.Errorf("Error from the server %s", err)
	}
	thesisId := msg.DocId
	alice.docId = thesisId
	fmt.Printf("Successfully added new shared doc 'Thesis' with id %s\n", thesisId)

	wg := &sync.WaitGroup{}
	watching := &sync.WaitGroup{}
	watching.Add(2)
	wg.Add(1)
	go alice.listenUpdate(wg, watching)

	conn = clientConn()
	bob := &Author{name: "Bob", doc: []byte{}, docId: thesisId, rev: 0, title: "", client: editorpb.NewNodeClient(conn)}
	wg.Add(1)
	go bob.listenUpdate(wg, watching)

	watching.Wait()

	ops := []*editorpb.Op{&editorpb.Op{N: 5}, &editorpb.Op{N: 0, S: " World!"}}
	bob.doc.Apply(ot.NewOps(ops))
	fmt.Printf("Bob update his local doc to %s\n", string(bob.doc))
	_, err = bob.client.Edit(context.Background(), &editorpb.EditReq{DocId: thesisId, Rev: int32(bob.rev), Ops: ops, UserId: bob.name, Title: "Thesis"})
	if err != nil {
		fmt.Printf("Error from the server %s\n", err)
	} else {
		fmt.Printf("Bob successfully updated shared doc %s\n", msg.DocId)
		// bob should retry or reverse his local doc edits
	}

	wg.Wait()
}

func clientConn() *grpc.ClientConn {
	serviceConfig := `{"healthCheckConfig": {"serviceName": "Example"}, "loadBalancingConfig": [ { "round_robin": {} } ]}`
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(5),
	}
	conn, err := grpc.NewClient("multi:///localhost:50051,localhost:50052,localhost:50053",
		grpc.WithDefaultServiceConfig(serviceConfig), grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)))
	if err != nil {
		fmt.Errorf("dialing failed: %v", err)
	}
	return conn
}
