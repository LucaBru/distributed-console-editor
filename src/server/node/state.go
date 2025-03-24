package node

import (
	serror "editor-service/errors"
	rlog "editor-service/node/log"
	"editor-service/node/ot"
	"editor-service/protos/editorpb"
	"editor-service/protos/logpb"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
)

type State struct {
	docs            map[DocId]*ot.SharedDoc
	shareDocCh      chan<- *rlog.ShareLog
	delDocCh        chan<- *rlog.DeleteLog
	editCh          chan<- *rlog.EditLog
	subListenerCh   chan<- *ListenerReq
	unsubListenerCh chan<- *ListenerReq
}

type DocId = uuid.UUID

type ShareLogResult = rlog.LogResult[*DocId]
type DelLogResult = rlog.LogResult[*struct{}]
type EditLogResult = rlog.LogResult[*struct{}]

type ListenerReq struct {
	docId     DocId
	userId    string
	stateCh   chan error
	updatesCh chan<- ot.Ops
}

func NewState() *State {
	shareCh := make(chan *rlog.ShareLog, 50)
	delCh := make(chan *rlog.DeleteLog, 50)
	editCh := make(chan *rlog.EditLog, 200)
	subCh := make(chan *ListenerReq, 50)
	unsubCh := make(chan *ListenerReq, 50)
	docs := make(map[DocId]*ot.SharedDoc)
	s := &State{docs: docs, shareDocCh: shareCh, delDocCh: delCh, editCh: editCh, subListenerCh: subCh, unsubListenerCh: unsubCh}
	go s.run(shareCh, delCh, editCh, subCh, unsubCh)
	return s
}

func (s *State) run(shareCh <-chan *rlog.ShareLog, delCh <-chan *rlog.DeleteLog, editCh <-chan *rlog.EditLog, subCh <-chan *ListenerReq, unsubCh <-chan *ListenerReq) {
	for {
		select {
		case req := <-shareCh:
			{
				uid, _ := uuid.Parse(req.Cmd.DocId)
				docCh := make(chan *rlog.EditLog, 100)
				s.docs[uid] = ot.NewSharedDoc(req.Cmd.DocName, req.Cmd.Doc, docCh)
				req.StateCh <- &ShareLogResult{Msg: &uid}
			}
		case del := <-delCh:
			{
				uid, err := uuid.Parse(del.Cmd.DocId)
				if err != nil {
					del.StateCh <- &DelLogResult{Err: &serror.DocIdError{Err: err}}
					break
				}
				doc := s.docs[uid]
				close(doc.EditCh)
				for _, listener := range doc.Listeners {
					close(listener)
				}
				delete(s.docs, uid)
				del.StateCh <- &DelLogResult{Err: nil}
			}
		case req := <-editCh:
			{
				uid, err := uuid.Parse(req.Cmd.DocId)
				if err != nil {
					req.StateCh <- &DelLogResult{Err: &serror.DocIdError{Err: err}}
					break
				}
				editCh := s.docs[uid].EditCh
				if editCh == nil {
					req.StateCh <- &DelLogResult{Err: &serror.SharedDocNotFound{}}
					break
				}
				editCh <- req
				req.StateCh <- &EditLogResult{}
			}
		case req := <-subCh:
			{
				doc := s.docs[req.docId]
				if doc == nil {
					req.stateCh <- &serror.SharedDocNotFound{}
					break
				}
				doc.Listeners[req.userId] = req.updatesCh
				req.stateCh <- nil
			}
		case req := <-unsubCh:
			{
				doc := s.docs[req.docId]
				if doc == nil {
					req.stateCh <- &serror.SharedDocNotFound{}
					break
				}
				delete(doc.Listeners, req.userId)
				req.stateCh <- nil
			}
		}
	}
}

func (s *State) Apply(l *raft.Log) interface{} {
	cmd := &logpb.Log{}
	err := proto.Unmarshal(l.Data, cmd)
	if err != nil {
		return fmt.Errorf("Log unmarshal failed: %w", err)
	}
	switch c := cmd.Cmd.(type) {
	case *logpb.Log_Share:
		{
			sync := make(chan *ShareLogResult)
			entry := &rlog.ShareLog{StateCh: sync, Cmd: c.Share}
			s.shareDocCh <- entry
			res := <-sync
			if res.Err != nil {
				return res.Err
			}
			return &editorpb.ShareReply{DocId: res.Msg.String()}
		}
	case *logpb.Log_Delete:
		{
			sync := make(chan *DelLogResult)
			entry := &rlog.DeleteLog{StateCh: sync, Cmd: c.Delete}
			s.delDocCh <- entry
			res := <-sync
			if res.Err != nil {
				return res.Err
			}
			return &editorpb.DeleteReply{}
		}
	case *logpb.Log_Edit:
		{
			sync := make(chan *EditLogResult)
			entry := &rlog.EditLog{StateCh: sync, Cmd: c.Edit}
			s.editCh <- entry
			res := <-sync
			if res.Err != nil {
				return fmt.Errorf("Failed to modify the document: %w", err)
			}
			return &editorpb.Ack{}
		}
	}
	return nil
}

func (s *State) SubListener(req *editorpb.ListenerReq, updatesCh chan<- ot.Ops) error {
	uid, err := uuid.Parse(req.DocId)
	if err != nil {
		return &serror.DocIdError{Err: err}
	}
	sync := make(chan error)
	subReq := &ListenerReq{
		docId: uid, userId: req.UserId, stateCh: sync, updatesCh: updatesCh,
	}
	s.subListenerCh <- subReq
	return <-sync
}

func (s *State) UnsubListener(req *editorpb.ListenerReq) error {
	uid, err := uuid.Parse(req.DocId)
	if err != nil {
		return &serror.DocIdError{Err: err}
	}
	sync := make(chan error)
	unsubReq := &ListenerReq{
		docId: uid, userId: req.UserId, stateCh: sync,
	}
	s.unsubListenerCh <- unsubReq
	return <-sync
}

func (s *State) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (s *State) Restore(snapshot io.ReadCloser) error {
	return nil
}
