package main

import (
	"client/editor"
	"flag"

	"github.com/nsf/termbox-go"
)

var (
	docId = flag.String("doc-id", "", "Document ID")
)

func main() {
	flag.Parse()

	err := termbox.Init()
	if err != nil {
		panic("Termbox init failed")
	}
	defer termbox.Close()

	// This input mode recognize escape characters
	termbox.SetInputMode(termbox.InputEsc)

	vEditor := editor.NewEditor(*docId)

	shouldExit := false

	for !shouldExit {
		vEditor.Draw()

		switch event := termbox.PollEvent(); event.Type {
		case termbox.EventKey:
			shouldExit = vEditor.OnKeyEvent(event)
		}
	}
}
