// Copyright 2013 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ot

import (
	"editor-service/protos/logpb"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type Doc []byte

type ClientId = uuid.UUID

// Apply applies the operation sequence ops to the document.
// An error is returned if applying ops failed.
func (doc *Doc) Apply(ops Ops) error {
	ret, del, ins := ops.Count()
	i, buf := 0, *doc
	if ret+del != len(buf) {
		return fmt.Errorf("Sum of unchanged and removed char must be equal to the doc length")
	}
	if max := ret + del + ins; max > cap(buf) {
		nbuf := make([]byte, len(buf), max+(max>>2))
		copy(nbuf, buf)
		buf = nbuf
	}
	for _, op := range ops {
		switch {
		case op.N > 0:
			i += op.N
		case op.N < 0:
			copy(buf[i:], buf[i-op.N:])
			buf = buf[:len(buf)+op.N]
		case op.S != "":
			l := len(buf)
			buf = buf[:l+len(op.S)]
			copy(buf[i+len(op.S):], buf[i:l])
			buf = append(buf[:i], op.S...)
			buf = buf[:l+len(op.S)]
			i += len(op.S)
		}
	}
	*doc = buf
	if i != ret+ins {
		return fmt.Errorf("Operation didn't operate on the whole document")
	}
	return nil
}

type LogCmd[T any, C any] struct {
	Ch  chan<- *LogReply[T]
	Cmd C
}

type LogReply[T any] struct {
	Msg T
	Err error
}

type DocId = uuid.UUID

type ShareEntry = LogCmd[*DocId, *logpb.Share]
type DeleteEntry = LogCmd[*struct{}, *logpb.Delete]
type EditEntry = LogCmd[*struct{}, *logpb.Edit]

// SharedDoc represents shared document with revision history.
type SharedDoc struct {
	name    string
	creator string
	doc     Doc
	history []Ops
	editCh  <-chan *EditEntry
}

func NewSharedDoc(name string, doc []byte, ch <-chan *EditEntry) {
	log.Println(fmt.Sprintf("New doc created, initial text: %v", doc))
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
				entry.Ch <- &LogReply[*struct{}]{Err: fmt.Errorf("Operations transformation failed: %w", err)}
			}
		}
		var old []byte
		copy(old, d.doc)
		if err = d.doc.Apply(ops); err != nil {
			entry.Ch <- &LogReply[*struct{}]{Err: fmt.Errorf("Operations application failed: %w", err)}
		}
		fmt.Println(fmt.Sprintf("Shared doc was updated from %s to %s", string(old), string(d.doc)))
		d.history = append(d.history, ops)
		// notify all the collaborators with new ops 	TODO:
		entry.Ch <- &LogReply[*struct{}]{}
	}
}
