package ot

import (
	rlog "editor-service/node/log"
	"fmt"
)

// SharedDoc represents shared document with revision history.
type SharedDoc struct {
	name    string
	creator string
	doc     Doc
	history []Ops
	editCh  <-chan *rlog.EditLog
}

func NewSharedDoc(name string, doc []byte, ch <-chan *rlog.EditLog) {
	d := &SharedDoc{
		doc:    doc,
		editCh: ch,
	}
	go d.run()
}

func (d *SharedDoc) run() {
	for entry := range d.editCh {
		rev := int(entry.Cmd.Rev)
		ops := NewOps(entry.Cmd.Ops)
		var err error
		if rev < 0 || len(d.history) < rev {
			fmt.Errorf("Revision not in history")
		}
		for _, other := range d.history[rev:] {
			if ops, _, err = Transform(ops, other); err != nil {
				entry.sendState <- &rlog.LogResult[*struct{}]{Err: fmt.Errorf("Operations transformation failed: %w", err)}
			}
		}
		var old []byte
		copy(old, d.doc)
		if err = d.doc.Apply(ops); err != nil {
			entry.sendState <- &rlog.LogResult[*struct{}]{Err: fmt.Errorf("Operations application failed: %w", err)}
		}
		fmt.Println(fmt.Sprintf("Shared doc was updated from '%s' to '%s'", string(old), string(d.doc)))
		d.history = append(d.history, ops)
		// notify all the collaborators with new ops 	TODO:
		entry.sendState <- &rlog.LogResult[*struct{}]{}
	}
}
