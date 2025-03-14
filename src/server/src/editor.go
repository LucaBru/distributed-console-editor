package service

import (
	"context"
	pb "editor-service/protos"
	"io"
	"log"
	"sync"

	"github.com/hashicorp/raft"
	"github.com/mb0/lab/ot"
)

type Editor struct {
	pb.UnimplementedEditorServer
	Doc  Document
	raft *raft.Raft
}

func NewEditor(r *raft.Raft) *Editor {
	return &Editor{raft: r}
}

type Document struct {
	Doc ot.Doc
	mu  sync.Mutex
}

func (d *Document) Apply(l *raft.Log) interface{} {
	return nil
}

func (d *Document) Snapshot() (raft.FSMSnapshot, error) {
	return &snapshot{}, nil
}

func (d *Document) Restore(r io.ReadCloser) error {
	return nil
}

type snapshot struct {
	doc ot.Doc
}

func (s *snapshot) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (s *snapshot) Release() {}

func (e *Editor) FetchUpdates(_ context.Context, req *pb.FetchUpdatesReq) (*pb.FetchUpdatesReply, error) {
	log.Println("Server FetchUpdates gRPC call")
	return nil, nil
}

func (e *Editor) PushOps(_ context.Context, req *pb.Ops) (*pb.PushOpsReply, error) {
	log.Println("Server PushOps gRPC call")
	return nil, nil
}
