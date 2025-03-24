package ot

import (
	rlog "editor-service/node/log"
	"fmt"
)

// SharedDoc represents shared document with revision history.
type SharedDoc struct {
	name      string
	creator   string
	doc       Doc
	history   []Ops
	EditCh    chan *rlog.EditLog
	Listeners map[string]chan<- Ops
}

func NewSharedDoc(name string, doc []byte, ch chan *rlog.EditLog) *SharedDoc {
	d := &SharedDoc{
		doc:       doc,
		EditCh:    ch,
		Listeners: make(map[string]chan<- Ops),
	}
	go d.run()
	return d
}

func (d *SharedDoc) run() {
	for log := range d.EditCh {
		rev := int(log.Cmd.Rev)
		ops := NewOps(log.Cmd.Ops)
		var err error
		if rev < 0 || len(d.history) < rev {
			fmt.Errorf("Revision not in history")
		}
		for _, other := range d.history[rev:] {
			if ops, _, err = Transform(ops, other); err != nil {
				log.StateCh <- &rlog.LogResult[*struct{}]{Err: fmt.Errorf("Operations transformation failed: %w", err)}
			}
		}
		var old []byte
		copy(old, d.doc)
		if err = d.doc.Apply(ops); err != nil {
			log.StateCh <- &rlog.LogResult[*struct{}]{Err: fmt.Errorf("Operations application failed: %w", err)}
		}
		fmt.Printf(fmt.Sprintf("Shared doc was updated from '%s' to '%s\n'", string(old), string(d.doc)))
		d.history = append(d.history, ops)

		// notify all the collaborators with new ops
		for userId, collaborator := range d.Listeners {
			if userId != log.Cmd.UserId {
				collaborator <- ops
			}
		}

		log.StateCh <- &rlog.LogResult[*struct{}]{}
	}
}
