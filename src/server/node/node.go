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
	log := &logpb.Log{Cmd: &logpb.Log_Share{Share: &logpb.Share{DocName: req.DocName, Doc: req.Doc, DocId: uuid.NewString()}}}
	return Apply[editorpb.ShareReply](n, log)
}

func (n *Node) Delete(ctx context.Context, req *editorpb.DeleteReq) (*editorpb.DeleteReply, error) {
	log := &logpb.Log{Cmd: &logpb.Log_Delete{Delete: &logpb.Delete{DocId: req.DocId, UserId: req.UserId}}}
	return Apply[editorpb.DeleteReply](n, log)
}

func (n *Node) Edit(ctx context.Context, req *editorpb.EditReq) (*editorpb.Ack, error) {
	log := &logpb.Log{Cmd: &logpb.Log_Edit{Edit: &logpb.Edit{DocId: req.DocId, Rev: req.Rev, Ops: req.Ops, UserId: req.UserId}}}
	return Apply[editorpb.Ack](n, log)
}

func Apply[R any](n *Node, log *logpb.Log) (*R, error) {
	b, err := proto.Marshal(log)
	if err != nil {
		return nil, &serror.InternalError{Err: fmt.Errorf("Failed to marshal the request: %w", err)}
	}
	f := n.Raft.Apply(b, time.Second)
	if err := f.Error(); err != nil {
		return nil, rafterrors.MarkRetriable(&serror.InternalError{Err: err})
	}
	iReply := f.Response()
	err, ok := iReply.(error)
	if ok {
		return nil, err
	}
	reply, ok := iReply.(*R)
	if ok {
		return reply, nil
	}
	return nil, &serror.InternalError{Err: errors.New("Failed to convert FSM response")}
}
