package node

import (
	serror "editor-service/errors"
	"editor-service/node/ot"
	"editor-service/protos/editorpb"
	"editor-service/protos/rlogpb"
	"fmt"
	"io"
	"sync"

	"github.com/google/uuid"
	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
)

type State struct {
	sync.RWMutex
	docs map[DocId]*ot.SharedDoc
}

type DocId = uuid.UUID

func NewState() *State {
	return &State{docs: make(map[DocId]*ot.SharedDoc)}
}

func (s *State) shareDoc(l *rlogpb.Share) {
	uid, _ := uuid.Parse(l.DocId)
	s.Lock()
	defer s.Unlock()
	s.docs[uid] = ot.NewSharedDoc(l.DocName, l.Doc, l.UserId)
}

func (s *State) deleteDoc(l *rlogpb.Delete) error {
	uid, err := uuid.Parse(l.DocId)
	if err != nil {
		return &serror.DocIdError{}
	}
	s.Lock()
	defer s.Unlock()
	doc := s.docs[uid]
	if doc == nil {
		return &serror.SharedDocNotFound{}
	}
	doc.Delete()
	delete(s.docs, uid)
	return nil
}

func (s *State) editDoc(l *rlogpb.Edit) error {
	uid, err := uuid.Parse(l.DocId)
	if err != nil {
		return &serror.DocIdError{}
	}
	s.Lock()
	defer s.Unlock()
	doc := s.docs[uid]
	if doc == nil {
		return &serror.SharedDocNotFound{}
	}
	err = doc.Edit(int(l.Rev), ot.NewOps(l.Ops), l.UserId, l.Title)
	if err != nil {
		return &serror.InternalError{Err: err}
	}
	return nil
}

func (s *State) Apply(l *raft.Log) interface{} {
	cmd := &rlogpb.Log{}
	err := proto.Unmarshal(l.Data, cmd)
	if err != nil {
		return serror.InternalError{Err: fmt.Errorf("Log unmarshal failed: %w", err)}
	}
	switch l := cmd.Cmd.(type) {
	case *rlogpb.Log_Share:
		{
			s.shareDoc(l.Share)
			return nil
		}
	case *rlogpb.Log_Delete:
		return s.deleteDoc(l.Delete)
	case *rlogpb.Log_Edit:
		return s.editDoc(l.Edit)
	default:
		return nil
	}
}

func (s *State) SubListener(req *editorpb.WatchReq) (<-chan ot.Update, []byte, string, int, error) {
	uid, err := uuid.Parse(req.DocId)
	s.RLock()
	defer s.RUnlock()
	if err != nil || s.docs[uid] == nil {
		return nil, nil, "", 0, &serror.DocIdError{}
	}
	recvUpdate, docContent, title, rev := s.docs[uid].AddListener(req.UserId)
	return recvUpdate, docContent, title, rev, nil
}

func (s *State) UnsubListener(req *editorpb.WatchReq) error {
	uid, err := uuid.Parse(req.DocId)
	s.RLock()
	defer s.RUnlock()
	if err != nil || s.docs[uid] == nil {
		return &serror.DocIdError{}
	}
	return s.docs[uid].DeleteListener(req.UserId)
}

func (s *State) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (s *State) Restore(snapshot io.ReadCloser) error {
	return nil
}
