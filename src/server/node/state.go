package node

import (
	rlog "editor-service/node/log"
	"editor-service/node/ot"
	"editor-service/protos/editorpb"
	"editor-service/protos/logpb"
	"fmt"
	"io"
	"log"

	"github.com/google/uuid"
	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
)

type State struct {
	docs       map[DocId]chan<- *rlog.EditLog
	shareDocCh chan *rlog.ShareLog
	delDocCh   chan *rlog.DeleteLog
	docChanCh  chan *DocChanReq
}

type DocId = uuid.UUID

type DocChanReq struct {
	docId   *DocId
	StateCh chan chan<- *rlog.EditLog
}

type ShareLogResult = rlog.LogResult[*DocId]
type DelLogResult = rlog.LogResult[*struct{}]

func NewState() *State {
	shareCh := make(chan *rlog.ShareLog, 50)
	delCh := make(chan *rlog.DeleteLog, 50)
	docCh := make(chan *DocChanReq)
	docs := make(map[DocId](chan<- *rlog.EditLog))
	s := &State{docs: docs, shareDocCh: shareCh, delDocCh: delCh, docChanCh: docCh}
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
		case req := <-s.docChanCh:
			{
				ch := s.docs[*req.docId]
				// send nil if the doc isn't available
				req.StateCh <- ch
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
			sync := make(chan *rlog.LogResult[*DocId])
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
			sync := make(chan *rlog.LogResult[*struct{}])
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
			ch := make(chan chan<- *rlog.EditLog)
			uid, err := uuid.Parse(c.Edit.DocId)
			if err != nil {
				log.Println("Invalid doc id")
				return fmt.Errorf("Share request has an invalid doc id: %w", err)
			}
			editChanReq := &DocChanReq{docId: &uid, StateCh: ch}
			s.docChanCh <- editChanReq
			editCh := <-ch
			if editCh == nil {
				return fmt.Errorf("The document to update was deleted")
			}
			sync := make(chan *rlog.LogResult[*struct{}])
			entry := &rlog.EditLog{StateCh: sync, Cmd: c.Edit}
			editCh <- entry
			reply := <-sync
			if reply.Err != nil {
				return fmt.Errorf("Failed to modify the document: %w", err)
			}
			return &editorpb.Ack{}
		}

	}
	// add default error
	return nil
}

func (s *State) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (s *State) Restore(snapshot io.ReadCloser) error {
	return nil
}
