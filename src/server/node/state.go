package node

import (
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
	docs       map[DocId]chan<- *rlog.EditLog
	shareDocCh chan *rlog.ShareLog
	delDocCh   chan *rlog.DeleteLog
	editCh     chan *rlog.EditLog
}

type DocId = uuid.UUID

type DocChanReq struct {
	docId   *DocId
	StateCh chan chan<- *rlog.EditLog
}

type ShareLogResult = rlog.LogResult[*DocId]
type DelLogResult = rlog.LogResult[*struct{}]
type EditLogResult = rlog.LogResult[*struct{}]

func NewState() *State {
	shareCh := make(chan *rlog.ShareLog, 50)
	delCh := make(chan *rlog.DeleteLog, 50)
	editCh := make(chan *rlog.EditLog)
	docs := make(map[DocId](chan<- *rlog.EditLog))
	s := &State{docs: docs, shareDocCh: shareCh, delDocCh: delCh, editCh: editCh}
	go s.run()
	return s
}

func (s *State) run() {
	for {
		select {
		case req := <-s.shareDocCh:
			{
				docCh := make(chan *rlog.EditLog, 100)
				ot.NewSharedDoc(req.Cmd.DocName, req.Cmd.Doc, docCh)
				uid, _ := uuid.Parse(req.Cmd.DocId)
				s.docs[uid] = docCh
				req.StateCh <- &ShareLogResult{Msg: &uid}
			}
		case del := <-s.delDocCh:
			{
				uid, err := uuid.Parse(del.Cmd.DocId)
				if err != nil {
					del.StateCh <- &DelLogResult{Err: fmt.Errorf("Share request has an invalid doc id: %w", err)}
					return
				}
				// terminate the doc.run goroutine, at the moment anyone that have the doc id can delete it
				close(s.docs[uid])
				delete(s.docs, uid)
				del.StateCh <- &DelLogResult{Err: nil}
			}
		case req := <-s.editCh:
			{
				uid, err := uuid.Parse(req.Cmd.DocId)
				if err != nil {
					req.StateCh <- &DelLogResult{Err: fmt.Errorf("Share request has an invalid doc id: %w", err)}
					return
				}
				editCh := s.docs[uid]
				if editCh == nil {
					req.StateCh <- &DelLogResult{Err: fmt.Errorf("The document to update was deleted")}
					return
				}
				editCh <- req
				req.StateCh <- &EditLogResult{}
			}
		}
	}
}

func (s *State) Apply(l *raft.Log) interface{} {
	cmd := &logpb.Log{}
	err := proto.Unmarshal(l.Data, cmd)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal log: %w", err)
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
	// add default error
	return nil
}

func (s *State) RegisterListener(docId *DocId) {

}

func (s *State) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (s *State) Restore(snapshot io.ReadCloser) error {
	return nil
}
