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

type DocId = uuid.UUID

type State struct {
	docs    map[DocId]chan<- *logpb.Edit
	shareCh chan *Share
}

type LogCmd[T any] struct {
	sendBackCh chan<- *T
}

type Share struct {
	LogCmd[struct{}]
	docId   DocId
	docName string
	doc     []byte
	userId  string
}

func NewState() *State {
	shareCh := make(chan *Share, 50)
	delCh := make(chan *logpb.Delete, 50)
	s := &State{docs: make(map[DocId](chan<- *logpb.Edit)), shareCh: shareCh}
	go s.run(delCh)
	return s
}

func (s *State) run(deleteCh <-chan *logpb.Delete) {
	for {
		select {
		case share := <-s.shareCh:
			{
				editCh := make(chan *logpb.Edit, 100)
				ot.NewSharedDoc(share.docName, share.doc, editCh)
				s.docs[share.docId] = editCh
				share.sendBackCh <- &struct{}{}
			}
		case d := <-deleteCh:
			{
				uid, err := uuid.Parse(d.DocId)
				if err != nil {
					log.Panicln("Error while casting an doc uuid to string")
				}
				delete(s.docs, uid)
			}
		}
	}
}

func (s *State) Apply(l *raft.Log) interface{} {
	cmd := &logpb.Log{}
	err := proto.Unmarshal(l.Data, cmd)
	if err != nil {
		return fmt.Errorf("Failed log conversion: %w", err)
	}
	switch c := cmd.Cmd.(type) {
	case *logpb.Log_Share:
		{
			ch := make(chan *struct{})
			uid, err := uuid.Parse(c.Share.DocId)
			if err != nil {
				return fmt.Errorf("Share request has an invalid doc id: %w", err)
			}
			share := &Share{docId: uid, docName: c.Share.DocName, doc: c.Share.Doc, userId: c.Share.UserId, LogCmd: LogCmd[struct{}]{sendBackCh: ch}}
			s.shareCh <- share
			<-ch
			return &editorpb.ShareReply{DocId: uid.String()}
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
