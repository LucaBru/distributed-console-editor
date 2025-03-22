package node

import (
	"context"
	serror "editor-service/errors"
	"editor-service/protos/editorpb"
	"editor-service/protos/logpb"
	"errors"
	"fmt"
	"time"

	"github.com/Jille/raft-grpc-leader-rpc/rafterrors"
	"github.com/google/uuid"
	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
)

type Node struct {
	editorpb.UnimplementedNodeServer
	Raft *raft.Raft
}

func NewNode(r *raft.Raft) *Node {
	return &Node{
		Raft: r,
	}
}

func (n *Node) Share(ctx context.Context, req *editorpb.ShareReq) (*editorpb.ShareReply, error) {
	entry := &logpb.Log{Cmd: &logpb.Log_Share{Share: &logpb.Share{DocName: req.DocName, Doc: req.Doc, DocId: uuid.NewString()}}}
	b, err := proto.Marshal(entry)
	if err != nil {
		return nil, &serror.InternalError{Err: fmt.Errorf("Failed to marshal the request: %w", err)}
	}
	f := n.Raft.Apply(b, time.Second)
	if err := f.Error(); err != nil {
		return nil, rafterrors.MarkRetriable(&serror.InternalError{Err: err})
	}
	reply, ok := f.Response().(*editorpb.ShareReply)
	if !ok {
		return nil, &serror.InternalError{Err: errors.New("Failed to convert the response")}
	}
	return reply, nil
}

func (n *Node) Delete(ctx context.Context, req *editorpb.DeleteReq) (*editorpb.DeleteReply, error) {
	entry := &logpb.Log{Cmd: &logpb.Log_Delete{Delete: &logpb.Delete{DocId: req.DocId, UserId: req.UserId}}}
	b, err := proto.Marshal(entry)
	if err != nil {
		return nil, &serror.InternalError{Err: fmt.Errorf("Failed to marshal the request: %w", err)}
	}
	f := n.Raft.Apply(b, time.Second)
	if err := f.Error(); err != nil {
		return nil, rafterrors.MarkRetriable(&serror.InternalError{Err: err})
	}
	return &editorpb.DeleteReply{}, nil
}

func (n *Node) Edit(ctx context.Context, req *editorpb.EditReq) (*editorpb.Ack, error) {
	entry := &logpb.Log{Cmd: &logpb.Log_Edit{Edit: &logpb.Edit{DocId: req.DocId, Rev: req.Rev, Ops: req.Ops, UserId: req.UserId}}}
	b, err := proto.Marshal(entry)
	if err != nil {
		return nil, &serror.InternalError{Err: fmt.Errorf("Failed to marshal the request: %w", err)}
	}
	f := n.Raft.Apply(b, time.Second)
	if err := f.Error(); err != nil {
		return nil, rafterrors.MarkRetriable(&serror.InternalError{Err: err})
	}
	err, ok := f.Response().(error)
	if ok {
		return nil, err
	}
	return &editorpb.Ack{}, nil
}
