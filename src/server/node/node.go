package node

import (
	"context"
	serror "editor-service/errors"
	"editor-service/node/ot"
	"editor-service/protos/editorpb"
	"editor-service/protos/logpb"
	"fmt"
	"io"
	"time"

	"github.com/Jille/raft-grpc-leader-rpc/rafterrors"
	"github.com/google/uuid"
	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
)

type Node struct {
	editorpb.UnimplementedNodeServer
	Raft  *raft.Raft
	state *State
}

func NewNode(r *raft.Raft, s State) *Node {
	return &Node{
		Raft:  r,
		state: &s,
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
		return nil, &serror.InternalError{Err: fmt.Errorf("Log marshal failed: %w", err)}
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
	return nil, &serror.InternalError{Err: fmt.Errorf("FSM response into %T conversion failed", reply)}
}

func (n *Node) HandleListener(stream editorpb.Node_HandleListenerServer) error {
	reqCh := make(chan *editorpb.ListenerReq)
	errCh := make(chan error)
	go func(stream *editorpb.Node_HandleListenerServer, reqCh chan<- *editorpb.ListenerReq, errCh chan<- error) {
		req, err := (*stream).Recv()
		if err == io.EOF {
			errCh <- nil
			return
		}
		if err != nil {
			errCh <- err
			return
		}
		reqCh <- req
	}(&stream, reqCh, errCh)

	listenUpdates := make(chan ot.Ops, 20)
	var req *editorpb.ListenerReq
	select {
	case req = <-reqCh:
		{
			err := n.state.SubListener(req, listenUpdates)
			if err != nil {
				return err
			}
		}
	case err := <-errCh:
		return err
	}

	for {
		select {
		case req := <-reqCh:
			return n.state.UnsubListener(req)

		case err := <-errCh:
			{
				n.state.UnsubListener(req)
				return err
			}
		case ops, ok := <-listenUpdates:
			stream.Send(&editorpb.Update{Ops: ops.WireFmt()})
			if !ok {
				return fmt.Errorf("Document you were listening to has been deleted")
			}
		}
	}
}
