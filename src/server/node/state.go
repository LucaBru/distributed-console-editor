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
		return serror.NewInvalidReqError(err)
	}
	s.Lock()
	defer s.Unlock()
	doc := s.docs[uid]
	if doc == nil {
		return serror.NewInvalidReqError(fmt.Errorf("doc with id %s not found", uid))
	}
	doc.Delete()
	delete(s.docs, uid)
	return nil
}

func (s *State) editDoc(l *rlogpb.Edit) error {
	uid, err := uuid.Parse(l.DocId)
	if err != nil {
		return serror.NewInvalidReqError(err)
	}
	s.Lock()
	defer s.Unlock()
	doc := s.docs[uid]
	if doc == nil {
		return serror.NewInvalidReqError(fmt.Errorf("doc with id %s not found", uid))
	}
	err = doc.Edit(int(l.Rev), ot.NewOps(l.Ops), l.UserId, l.Title)
	if err != nil {
		return serror.NewInternalError(err)
	}
	return nil
}

func (s *State) Apply(l *raft.Log) interface{} {
	cmd := &rlogpb.Log{}
	err := proto.Unmarshal(l.Data, cmd)
	if err != nil {
		return serror.InternalError(fmt.Errorf("Log unmarshal failed: %w", err))
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
		return nil, nil, "", 0, serror.NewInvalidReqError(err)
	}
	recvUpdate, docContent, title, rev := s.docs[uid].AddListener(req.UserId)
	return recvUpdate, docContent, title, rev, nil
}

func (s *State) UnsubListener(req *editorpb.WatchReq) error {
	uid, err := uuid.Parse(req.DocId)
	s.RLock()
	defer s.RUnlock()
	if err != nil || s.docs[uid] == nil {
		return serror.NewInvalidReqError(err)
	}
	return s.docs[uid].DeleteListener(req.UserId)
}

func (s *State) Snapshot() (raft.FSMSnapshot, error) {
	s.RLock()
	docsCopy := make(map[uuid.UUID]*ot.SharedDoc, len(s.docs))
	for docId, doc := range s.docs {
		docsCopy[docId] = doc.Clone()
	}
	s.RUnlock()

	return &StateSnapshot{
		docs: docsCopy,
	}, nil
}

func (s *State) Restore(snapshot io.ReadCloser) error {
	b, err := io.ReadAll(snapshot)
	if err != nil {
		return serror.NewInternalError(err)
	}
	docsMsg := &editorpb.Docs{}
	err = proto.Unmarshal(b, docsMsg)
	if err != nil {
		return serror.NewInternalError(err)
	}
	docs := make(map[uuid.UUID]*ot.SharedDoc)
	for docId, doc := range docsMsg.Docs {
		uid, _ := uuid.Parse(docId)
		docs[uid] = ot.NewSharedDoc(doc.Title, doc.Content, doc.Author)
	}
	s.docs = docs
	return nil
}

type StateSnapshot struct {
	docs map[uuid.UUID]*ot.SharedDoc
}

func (s *StateSnapshot) Persist(sink raft.SnapshotSink) error {
	docs := make(map[string]*editorpb.Doc, len(s.docs))
	for docId, doc := range s.docs {
		docs[docId.String()] = doc.ToProtoMsg()
	}
	state := &editorpb.Docs{
		Docs: docs,
	}
	b, err := proto.Marshal(state)
	if err != nil {
		return serror.NewInternalError(err)
	}

	_, err = sink.Write(b)
	if err != nil {
		sink.Cancel()
		return serror.NewInternalError(err)
	}

	return sink.Close()
}

func (S *StateSnapshot) Release() {}
