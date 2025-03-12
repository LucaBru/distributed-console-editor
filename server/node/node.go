package node

import (
	"dconsole-editor/raft"
	"github.com/mb0/lab/ot"
)

type Node struct {
	Doc ot.Doc
	raft.NodeModule
}
