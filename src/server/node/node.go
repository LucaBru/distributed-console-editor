package node

import (
	"context"
	"fmt"
	"io"
	"log"
	logMsg "log"
	"sync"
	"time"

	serror "editor-service/errors"
	"editor-service/protos/editorpb"
	"editor-service/protos/rlogpb"

	"github.com/Jille/raft-grpc-leader-rpc/rafterrors"
	"github.com/google/uuid"
	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
)

type Node struct {
	editorpb.UnimplementedNodeServer
	Raft                 *raft.Raft
	state                *State
	editReceivedTime     map[int32]time.Time
	editReceivedTimeLock sync.Mutex
	replicationLock      sync.Mutex
}

func NewNode(r *raft.Raft, s *State) *Node {
	return &Node{
		Raft:                 r,
		state:                s,
		editReceivedTime:     make(map[int32]time.Time),
		editReceivedTimeLock: sync.Mutex{},
		replicationLock:      sync.Mutex{},
	}
}

func (n *Node) Share(ctx context.Context, req *editorpb.ShareReq) (*editorpb.ShareReply, error) {
	docId := uuid.NewString()
	log := &rlogpb.Log{Cmd: &rlogpb.Log_Share{Share: &rlogpb.Share{DocName: req.DocName, Doc: req.Doc, DocId: docId}}}
	err := n.replicateLog(log)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Author %s share new doc %s\n", req.UserId, docId)
	return &editorpb.ShareReply{DocId: docId}, nil
}

func (n *Node) Delete(ctx context.Context, req *editorpb.DeleteReq) (*editorpb.DeleteReply, error) {
	log := &rlogpb.Log{Cmd: &rlogpb.Log_Delete{Delete: &rlogpb.Delete{DocId: req.DocId, UserId: req.UserId}}}
	err := n.replicateLog(log)
	if err != nil {
		return nil, err
	}
	return &editorpb.DeleteReply{}, nil
}

func (n *Node) Edit(ctx context.Context, req *editorpb.EditReq) (*editorpb.Ack, error) {
	n.editReceivedTimeLock.Lock()
	n.editReceivedTime[req.Rev] = time.Now()
	n.editReceivedTimeLock.Unlock()

	log := &rlogpb.Log{Cmd: &rlogpb.Log_Edit{Edit: &rlogpb.Edit{DocId: req.DocId, Rev: req.Rev, Ops: req.Ops, UserId: req.UserId, Title: req.Title}}}
	n.replicationLock.Lock()
	err := n.replicateLog(log)
	n.replicationLock.Unlock()
	if err != nil {
		return nil, err
	}
	return &editorpb.Ack{}, nil
}

func (n *Node) replicateLog(log *rlogpb.Log) error {
	b, err := proto.Marshal(log)
	if err != nil {
		return serror.NewInternalError(err)
	}

	f := n.Raft.Apply(b, time.Second)
	if err := f.Error(); err != nil {
		return rafterrors.MarkRetriable(serror.NewInternalError(err))
	}
	reply := f.Response()
	err, ok := reply.(error)
	if ok {
		logMsg.Print("failed to apply the log to the state: ", err)
		return err
	}

	return nil
}

func (n *Node) WatchDocument(stream editorpb.Node_WatchDocumentServer) error {
	req, err := stream.Recv()
	log.Print("watch document request from ", req.UserId)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	recvUpdate, doc, title, rev, err := n.state.SubListener(req)
	if err != nil {
		log.Fatal("failed to subscribe listener: ", err)
		return err
	}
	stream.Send(&editorpb.Update{Doc: doc, Rev: int32(rev), Title: title})

	for {
		update, ok := <-recvUpdate
		err := stream.Send(&editorpb.Update{Ops: update.Ops.WireFmt(), Title: update.Title})
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}
}
