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
	"time"

	_ "github.com/Jille/grpc-multi-resolver"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/health"
)

func main() {
	var aliceDoc ot.Doc
	aliceDoc = []byte("Hello")
	bobDoc := aliceDoc
	aliceConn := clientConn("Alice")
	alice := editorpb.NewNodeClient(aliceConn)
	msg, err := alice.Share(context.Background(), &editorpb.ShareReq{DocName: "Thesis", Doc: aliceDoc, UserId: "user@mail.com"})
	if err != nil {
		fmt.Errorf("Error from the server %s", err)
	}
	thesisId := msg.DocId
	fmt.Printf("Successfully added new shared doc 'Thesis' with id %s\n", thesisId)
	aliceStream, _ := alice.HandleListener(context.Background())
	aliceStream.Send(&editorpb.ListenerReq{DocId: thesisId, UserId: "Alice"})
	waitAlice := make(chan struct{})
	go func() {
		for {
			in, err := aliceStream.Recv()
			if err == io.EOF {
				// read done.
				close(waitAlice)
				return
			}
			if err != nil {
				fmt.Errorf("Failed to receive a note : %v", err)
			}
			aliceDoc.Apply(ot.NewOps(in.Ops))
			fmt.Printf("Alice got message %v, her local doc is updated to: %s\n", in.Ops, aliceDoc)
		}
	}()

	bobConn := clientConn("Bob")
	bob := editorpb.NewNodeClient(bobConn)
	bobStream, _ := bob.HandleListener(context.Background())
	bobStream.Send(&editorpb.ListenerReq{DocId: thesisId, UserId: "Bob"})
	waitBob := make(chan struct{})
	go func() {
		for {
			in, err := bobStream.Recv()
			if err == io.EOF {
				// read done.
				close(waitBob)
				return
			}
			if err != nil {
				fmt.Errorf("Failed to receive a note : %v", err)
			}
			bobDoc.Apply(ot.NewOps(in.Ops))
			fmt.Printf("Bob got message %v, her local doc is updated to: %s\n", in.Ops, aliceDoc)
		}
	}()
	ops := []*editorpb.Op{&editorpb.Op{N: 5}, &editorpb.Op{N: 0, S: " World!"}}
	bobDoc.Apply(ot.NewOps(ops))
	fmt.Printf("Bob update his local doc to %s\n", string(bobDoc))
	_, err = bob.Edit(context.Background(), &editorpb.EditReq{DocId: thesisId, Rev: 0, Ops: ops, UserId: "Bob"})
	if err != nil {
		fmt.Errorf("Error from the server %s", err)
	}
	fmt.Printf("Bob successfully updated shared doc %s\n", msg.DocId)
	bobStream.CloseSend()
	aliceStream.CloseSend()
	<-waitBob
	<-waitAlice

}

func clientConn(userId string) *grpc.ClientConn {
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
