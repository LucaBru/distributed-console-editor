package ot

import (
	"fmt"
	"sync"
)

type SharedDoc struct {
	sync.RWMutex
	title     string
	author    string
	doc       Doc
	history   []Ops
	listeners map[string]chan<- Update
}

func NewSharedDoc(title string, doc []byte, author string) *SharedDoc {
	d := &SharedDoc{
		doc:       doc,
		title:     title,
		author:    author,
		listeners: make(map[string]chan<- Update),
	}
	return d
}

func (d *SharedDoc) AddListener(listenerId string) (<-chan Update, []byte, string, int) {
	sendUpdate := make(chan Update)
	d.Lock()
	defer d.Unlock()
	d.listeners[listenerId] = sendUpdate
	return sendUpdate, d.doc, d.title, len(d.history)
}

func (d *SharedDoc) DeleteListener(listenerId string) error {
	d.Lock()
	defer d.Unlock()
	if d.listeners[listenerId] == nil {
		return fmt.Errorf("Listener id not found")
	}
	close(d.listeners[listenerId])
	delete(d.listeners, listenerId)
	return nil
}

func (d *SharedDoc) Delete() {
	d.Lock()
	defer d.Unlock()
	for _, ch := range d.listeners {
		close(ch)
	}
}

func (d *SharedDoc) Edit(rev int, ops Ops, authorId string, title string) error {
	if rev < 0 || len(d.history) < rev {
		return fmt.Errorf("Revision not in history")
	}

	var err error
	for _, other := range d.history[rev:] {
		if ops, _, err = Transform(ops, other); err != nil {
			return fmt.Errorf("Operations transformation failed: %w", err)
		}
	}

	var old []byte
	copy(old, d.doc)
	d.Lock()
	defer d.Unlock()
	if err = d.doc.Apply(ops); err != nil {
		return fmt.Errorf("Operations application failed: %w", err)
	}
	fmt.Printf(fmt.Sprintf("Shared doc was updated from '%s' to '%s\n'", string(old), string(d.doc)))
	d.history = append(d.history, ops)

	// notify all the collaborators with new ops
	for id, ch := range d.listeners {
		if id != authorId {
			ch <- Update{Ops: ops, Title: d.title}
		}
	}

	if title != "" && title != d.title {
		d.title = title
	}

	return nil
}

type Update struct {
	Ops   Ops
	Title string
}
