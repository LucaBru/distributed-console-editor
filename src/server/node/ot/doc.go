// Copyright 2013 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ot

import (
	"fmt"
)

type Doc []byte

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
