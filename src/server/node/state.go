package node

import (
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
	docs    map[ot.DocId]chan<- *ot.EditEntry
	shareCh chan *ot.ShareEntry
	delCh   chan *ot.DeleteEntry
	errCh   chan error
}

func NewState() *State {
	shareCh := make(chan *ot.ShareEntry, 50)
	errCh := make(chan error)
	delCh := make(chan *ot.DeleteEntry, 50)
	docs := make(map[ot.DocId](chan<- *ot.EditEntry))
	s := &State{docs: docs, shareCh: shareCh, errCh: errCh, delCh: delCh}
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
				uid, err := uuid.Parse(share.Cmd.DocId)
				if err != nil {
					s.errCh <- fmt.Errorf("Share request has an invalid doc id: %w", err)
					return
				}
				s.docs[uid] = editCh
				share.SendBackCh <- &uid
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
				del.SendBackCh <- nil
			}
		}
	}
}

func (s *State) Apply(l *raft.Log) interface{} {
	cmd := &logpb.Log{}
	err := proto.Unmarshal(l.Data, cmd)
	if err != nil {
		return fmt.Errorf("Failed log un-marshalling: %w", err)
	}
	switch c := cmd.Cmd.(type) {
	case *logpb.Log_Share:
		{
			sync := make(chan *ot.DocId)
			entry := &ot.ShareEntry{SendBackCh: sync, Cmd: c.Share}
			s.shareCh <- entry
			docId := <-sync
			return &editorpb.ShareReply{DocId: docId.String()}
		}
	case *logpb.Log_Delete:
		{
			sync := make(chan *struct{})
			entry := &ot.DeleteEntry{SendBackCh: sync, Cmd: c.Delete}
			s.delCh <- entry
			<-sync
			return &editorpb.DeleteReply{}
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
