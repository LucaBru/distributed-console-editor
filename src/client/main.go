package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"client/editor"

	"github.com/nsf/termbox-go"
)

var (
	docId    = flag.String("doc-id", "", "Document ID")
	authorId = flag.String("author", "", "Author ID")
)

func init() {
	flag.Parse()
	logFile, err := os.OpenFile("analysis/logs/"+*authorId+".log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		log.Println("Unable to create Logger file:", err.Error())
		return
	}
	log.SetOutput(logFile)
}

func main() {
	log.Println("start the program")
	err := termbox.Init()
	if err != nil {
		panic("Termbox init failed")
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	screen := editor.NewScreen()
	editor := editor.NewEditor(*docId, *authorId)
	go screen.DrawRemoteEdit(editor.Observe())

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		timer := time.NewTimer(50 * time.Second)
		editTicker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-timer.C:
				return
			case <-editTicker.C:
				{
					log.Print("edit sent")
					editor.Edit(termbox.Event{
						Type: termbox.EventKey,
						Ch:   rune(rand.Intn(26) + 97),
					})
				}
			}
		}
	}()
	wg.Wait()
}
