package main

import (
	"log"
	"client/src"

	"github.com/nsf/termbox-go"
)

var(
	logger = log.Default()
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic("Termbox init failed")
	}
	defer termbox.Close()
	
	// This input mode recognize escape characters
	termbox.SetInputMode(termbox.InputEsc)

	vEditor := editor.NewEditor()

	shouldExit := false

	for !shouldExit {
		vEditor.Draw()

		switch event := termbox.PollEvent(); event.Type {
		case termbox.EventKey:
			shouldExit = vEditor.OnKeyEvent(event)
		}
	}
}