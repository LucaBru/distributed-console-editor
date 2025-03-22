// Binary hammer sends requests to your Raft cluster as fast as it can.
// It sends the written out version of the Dutch numbers up to 2000.
// In the end it asks the Raft cluster what the longest three words were.
package main

import (
	"context"
	"editor-service/protos/editorpb"
	"log"
	"time"

	_ "github.com/Jille/grpc-multi-resolver"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/health"
)

func main() {
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
		log.Fatalf("dialing failed: %v", err)
	}
	defer conn.Close()
	client := editorpb.NewNodeClient(conn)
	msg, err := client.Share(context.Background(), &editorpb.ShareReq{DocName: "Hello", Doc: []byte("Hello"), UserId: "user@mail.com"})
	if err != nil {
		log.Fatalln("Error from the server %s", err)
	}
	log.Println("Successfully added new shared doc %s", msg.DocId)

	ops := []*editorpb.Op{&editorpb.Op{N: 5}, &editorpb.Op{N: 0, S: " World!"}}
	_, err = client.Edit(context.Background(), &editorpb.EditReq{DocId: msg.DocId, Rev: 0, Ops: ops, UserId: "user@mail.com"})
	if err != nil {
		log.Fatalln("Error from the server %s", err)
	}
	log.Println("Successfully updated shared doc %s", msg.DocId)

	_, err = client.Delete(context.Background(), &editorpb.DeleteReq{DocId: msg.DocId, UserId: "user@mail.com"})
	if err != nil {
		log.Fatalln("Error from the server %s", err)
	}
	log.Println("Successfully deleted new shared doc %s", msg.DocId)

}
