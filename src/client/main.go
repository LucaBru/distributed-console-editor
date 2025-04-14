package main

import (
	"flag"

	"client/editor"

	"github.com/nsf/termbox-go"
)

var docId = flag.String("doc-id", "", "Document ID")

func main() {
	flag.Parse()

	err := termbox.Init()
	if err != nil {
		panic("Termbox init failed")
	}
	defer termbox.Close()

	// This input mode recognize escape characters
	termbox.SetInputMode(termbox.InputEsc)

	vEditor, recvUpdate := editor.NewEditor(*docId)

	shouldExit := false

	recvKeyDigit := make(chan *termbox.Event, 20)

	go func() {
		for !shouldExit {
			e := termbox.PollEvent()
			recvKeyDigit <- &e
		}
	}()

	for !shouldExit {
		vEditor.Draw()
		select {
		case event := <-recvKeyDigit:
			{
				switch event.Type {
				case termbox.EventKey:
					shouldExit = vEditor.OnKeyEvent(*event)
				}
			}
		case <-recvUpdate:
		}
	}
}
