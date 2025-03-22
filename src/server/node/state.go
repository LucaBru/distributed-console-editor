package node

import (
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
	docs    map[ot.DocId]chan<- *ot.EditEntry
	shareCh chan *ot.ShareEntry
	delCh   chan *ot.DeleteEntry
	docCh   chan *DocEditReq
	errCh   chan error
}

type DocEditReq struct {
	docId      *ot.DocId
	sendBackCh chan chan<- *ot.EditEntry
}

func NewState() *State {
	shareCh := make(chan *ot.ShareEntry, 50)
	errCh := make(chan error)
	delCh := make(chan *ot.DeleteEntry, 50)
	docCh := make(chan *DocEditReq)
	docs := make(map[ot.DocId](chan<- *ot.EditEntry))
	s := &State{docs: docs, shareCh: shareCh, errCh: errCh, delCh: delCh, docCh: docCh}
	go s.run()
	return s
}

func (s *State) run() {
	for {
		select {
		case share := <-s.shareCh:
			{
				editCh := make(chan *ot.EditEntry, 100)
				ot.NewSharedDoc(share.Cmd.DocName, share.Cmd.Doc, editCh)
				log.Println("New shared doc created")
				uid, err := uuid.Parse(share.Cmd.DocId)
				if err != nil {
					s.errCh <- fmt.Errorf("Share request has an invalid doc id: %w", err)
					return
				}
				s.docs[uid] = editCh
				share.Ch <- &ot.LogReply[*ot.DocId]{Msg: &uid}
			}
		case del := <-s.delCh:
			{
				uid, err := uuid.Parse(del.Cmd.DocId)
				if err != nil {
					s.errCh <- fmt.Errorf("Share request has an invalid doc id: %w", err)
					return
				}
				// terminate the doc.run goroutine, at the moment anyone that have the doc id can delete it
				close(s.docs[uid])
				delete(s.docs, uid)
				del.Ch <- nil
			}
		case req := <-s.docCh:
			{
				ch := s.docs[*req.docId]
				req.sendBackCh <- ch
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
			sync := make(chan *ot.LogReply[*ot.DocId])
			entry := &ot.ShareEntry{Ch: sync, Cmd: c.Share}
			s.shareCh <- entry
			reply := <-sync
			if reply.Err != nil {
				return reply.Err
			}
			return &editorpb.ShareReply{DocId: reply.Msg.String()}
		}
	case *logpb.Log_Delete:
		{
			sync := make(chan *ot.LogReply[*struct{}])
			entry := &ot.DeleteEntry{Ch: sync, Cmd: c.Delete}
			s.delCh <- entry
			<-sync
			return &editorpb.DeleteReply{}
		}
	case *logpb.Log_Edit:
		{
			ch := make(chan chan<- *ot.EditEntry)
			uid, err := uuid.Parse(c.Edit.DocId)
			if err != nil {
				log.Println("Invalid doc id")
				return fmt.Errorf("Share request has an invalid doc id: %w", err)
			}
			editChanReq := &DocEditReq{docId: &uid, sendBackCh: ch}
			s.docCh <- editChanReq
			editCh := <-ch
			sync := make(chan *ot.LogReply[*struct{}])
			entry := &ot.EditEntry{Ch: sync, Cmd: c.Edit}
			// note that can panic, handle delete document better! TODO:
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
