package ot

import (
	rlog "editor-service/node/log"
	"fmt"
	"sync"
)

// SharedDoc represents shared document with revision history.
type SharedDoc struct {
	sync.RWMutex
	name      string
	creator   string
	doc       Doc
	history   []Ops
	Listeners map[string]chan<- Ops
}

func NewSharedDoc(name string, doc []byte, ch chan *rlog.EditLog) *SharedDoc {
	d := &SharedDoc{
		doc:       doc,
		Listeners: make(map[string]chan<- Ops),
	}
	return d
}

func (d *SharedDoc) AddListener(listenerId string, ch chan<- Ops) {
	d.Lock()
	defer d.Unlock()
	d.Listeners[listenerId] = ch
}

func (d *SharedDoc) DeleteListener(listenerId string) {
	d.Lock()
	defer d.Unlock()
	delete(d.Listeners, listenerId)
}

func (d *SharedDoc) Edit(log *rlog.EditLog) (Ops, error) {
	rev := int(log.Cmd.Rev)
	ops := NewOps(log.Cmd.Ops)

	if rev < 0 || len(d.history) < rev {
		return nil, fmt.Errorf("Revision not in history")
	}

	var err error
	for _, other := range d.history[rev:] {
		if ops, _, err = Transform(ops, other); err != nil {
			return nil, fmt.Errorf("Operations transformation failed: %w", err)
		}
	}

	var old []byte
	copy(old, d.doc)
	d.Lock()
	defer d.Unlock()
	if err = d.doc.Apply(ops); err != nil {
		return nil, fmt.Errorf("Operations application failed: %w", err)
	}
	fmt.Printf(fmt.Sprintf("Shared doc was updated from '%s' to '%s\n'", string(old), string(d.doc)))
	d.history = append(d.history, ops)

	// notify all the collaborators with new ops
	for userId, collaborator := range d.Listeners {
		if userId != log.Cmd.UserId {
			collaborator <- ops
		}
	}

	return ops, nil
}
