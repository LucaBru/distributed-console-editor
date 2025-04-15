package editor

import (
	"bufio"
	"editor-service/node/ot"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"client/sync_manager"

	// "log"

	"github.com/nsf/termbox-go"
)

var Writer = initWriter()

func initWriter() *bufio.Writer {
	file, err := os.Create("editor.log")
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return nil
	}
	return bufio.NewWriter(file)
}

type Editor struct {
	buffer                []string // These are the lines of text
	cursor                Cursor
	offsetX               int    // Scroll on X axis
	offsetY               int    // Scroll on Y axis
	filename              string // File where text will be saved
	modified              bool
	statusMsg             string
	backgroundColor       termbox.Attribute
	foregroundColor       termbox.Attribute
	statusBackgroundColor termbox.Attribute
	statusForegroundColor termbox.Attribute
	syncManager           *sync_manager.SyncManager
}

func NewEditor(docId string) (*Editor, <-chan struct{}) {
	authorId := rand.Int()
	syncManager, recvUpdate := sync_manager.NewSyncManager(sync_manager.NewDocumentConfig(docId, "AuthorN"+fmt.Sprint(authorId), 0, "Testing 1", ot.Doc{}))
	editor := &Editor{
		buffer:                []string{""},
		backgroundColor:       termbox.ColorDefault,
		foregroundColor:       termbox.ColorDefault,
		statusBackgroundColor: termbox.ColorBlack,
		statusForegroundColor: termbox.ColorWhite,
		filename:              "untitled.txt",
		cursor:                *newCursor(),
		syncManager:           syncManager,
	}
	// editor.setStatus("Author id: " + fmt.Sprintf("%d", authorId))
	return editor, recvUpdate
}

// Draw editor content to the terminal
func (editor *Editor) Draw() {
	// TODO: get a copy, not the document (pay attention to the use of syncManager (lock needed))
	updatedDoc := editor.syncManager.DocConfig.Document
	lineCounter := 0
	editor.buffer = []string{""}
	fmt.Fprintf(Writer, "ot document %v\n", updatedDoc)
	Writer.Flush()
	for i, digit := range updatedDoc {
		if digit == 0x0A {
			fmt.Fprintf(Writer, "Encountered new line while redrawing\n")
			Writer.Flush()
			editor.buffer = append(editor.buffer, "")
			lineCounter++
		} else {
			editor.buffer[lineCounter] += string(updatedDoc[i])
		}
	}

	// We clear the current text on the screen
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, height := termbox.Size()

	editor.drawText(width, height)
	editor.drawStatus(width, height)

	termbox.SetCursor(editor.cursor.x-editor.offsetX, editor.cursor.y-editor.offsetY)
	termbox.Flush()
}

func (editor *Editor) drawText(width int, height int) {
	for y := 0; y < height-1; y++ {
		lineY := y + editor.offsetY
		if lineY >= len(editor.buffer) {
			// In this case we are trying to display a line that is not in the buffer
			break
		}

		lineContent := editor.buffer[lineY]
		if editor.offsetX < len(lineContent) {
			// We display only the part of the string after the horizontal offset
			displayLine := lineContent[editor.offsetX:]
			for x, char := range []rune(displayLine) {
				if x >= width {
					break
				}
				termbox.SetCell(x, y, char, editor.backgroundColor, editor.foregroundColor)
			}
		}
	}
}

func (editor *Editor) drawStatus(width int, height int) {
	// Now we draw the status line
	statusLine := fmt.Sprintf(" %s - %d lines %s", editor.syncManager.DocConfig.DocId, len(editor.buffer), map[bool]string{true: "[modified]", false: ""}[editor.modified])
	if editor.statusMsg != "" {
		statusLine = editor.statusMsg
	}

	// We fill the status line with spaces
	for x := 0; x < width; x++ {
		termbox.SetCell(x, height-1, ' ', editor.statusBackgroundColor, editor.statusForegroundColor)
	}

	// Draw the status
	for x, char := range []rune(statusLine) {
		if x >= width {
			break
		}
		termbox.SetCell(x, height-1, char, editor.statusBackgroundColor, editor.statusForegroundColor)
	}
}

func (editor *Editor) cursorBounds() (int, int) {
	prev := 0
	next := len(editor.buffer[editor.cursor.y]) - editor.cursor.x
	for _, line := range editor.buffer[:editor.cursor.y] {
		prev += len(line) + 1
	}
	prev += editor.cursor.x

	for _, line := range editor.buffer[editor.cursor.y+1:] {
		next += len(line)
	}
	return prev, next
}

func (editor *Editor) insertRune(char rune) {
	ops := ot.Ops{}
	beforeCursor, afterCursor := editor.cursorBounds()
	ops = append(ops, ot.Op{N: beforeCursor})

	line := []rune(editor.buffer[editor.cursor.y])
	for i := 0; i < editor.cursor.x-len(line); i++ {
		line = append(line, ' ')
		ops = append(ops, ot.Op{N: 0, S: " "})
	}

	ops = append(ops, ot.Op{N: 0, S: string(char)})
	ops = append(ops, ot.Op{N: afterCursor})
	editor.syncManager.ApplyEdit(ops)
	editor.scrollRight()
	editor.modified = true
}

func (editor *Editor) insertNewline() {
	beforeCursor, afterCursor := editor.cursorBounds()
	fmt.Fprintf(Writer, "Inserting new line %d %d %v", beforeCursor, afterCursor, editor.syncManager.DocConfig.Document)
	Writer.Flush()
	editor.syncManager.ApplyEdit(ot.Ops{ot.Op{N: beforeCursor}, ot.Op{N: 0, S: "\n"}, ot.Op{N: afterCursor}})

	editor.cursor.moveDown()
	_, height := termbox.Size()
	if editor.cursor.y > height-2 {
		editor.offsetY++
	}
	editor.cursor.returnToTheBeginOfTheLine()
	Writer.Flush()
	editor.modified = true
}

// DeleteChar deletes the character at the current cursor position
func (editor *Editor) deleteChar() {
	if editor.cursor.x >= len(editor.buffer[editor.cursor.y]) {
		fmt.Fprintln(Writer, "Try to delete an empty char")
		Writer.Flush()
		return
	}

	beforeCursor, afterCursor := editor.cursorBounds()
	editor.setStatus(fmt.Sprintf("delete char in pos %d", beforeCursor))
	editor.syncManager.ApplyEdit(ot.Ops{ot.Op{N: beforeCursor}, ot.Op{N: -1}, ot.Op{N: afterCursor - 1}})
	editor.modified = true
}

// SaveFile saves the current buffer to a file
func (editor *Editor) saveFile() {
	content := strings.Join(editor.buffer, "\n")
	file, file_err := os.Create(editor.filename)
	if file_err != nil {
		editor.setStatus("Error opening file: " + file_err.Error())
		return
	}
	err := os.WriteFile(file.Name(), []byte(content), 0o644)
	if err != nil {
		editor.setStatus("Error saving file: " + err.Error())
	} else {
		editor.modified = false
		editor.setStatus(fmt.Sprintf("Saved %s (%d bytes)", editor.filename, len(content)))
	}
}

func (editor *Editor) scrollRight() {
	_, width := termbox.Size()
	if editor.offsetX > 0 || editor.cursor.x >= width {
		editor.cursor.moveRight()
		editor.offsetX++
		return
	}
	editor.cursor.moveRight()
}

func (editor *Editor) scrollLeft() {
	if editor.offsetX > 0 {
		editor.cursor.moveLeft()
		editor.offsetX--
		return
	}
	editor.cursor.moveLeft()
}

func (editor *Editor) scrollUp() {
	if editor.cursor.y == 0 {
		return
	}
	if editor.offsetY > 0 && editor.cursor.y > 0 {
		editor.cursor.goToTheEndOfPreviousLine(len(editor.buffer[editor.cursor.y-1]))
		editor.offsetY--
		return
	}
	editor.cursor.goToTheEndOfPreviousLine(len(editor.buffer[editor.cursor.y-1]))
}

func (editor *Editor) scrollDown() {
	if editor.cursor.y+1 == len(editor.buffer) {
		return
	}
	_, height := termbox.Size()
	if editor.cursor.y > height {
		editor.cursor.goToTheEndOfNextLine(len(editor.buffer[editor.cursor.y+1]))
		editor.offsetY++
		return
	}
	editor.cursor.goToTheEndOfNextLine(len(editor.buffer[editor.cursor.y+1]))
}

// SetStatus sets a temporary status message
func (e *Editor) setStatus(msg string) {
	e.statusMsg = msg
}

// Callback to handle key events obtained from termbox.
// Returns true if editor should be closed
func (editor *Editor) OnKeyEvent(event termbox.Event) bool {
	switch event.Key {
	case termbox.KeyCtrlQ:
		// We should exit
		return true
	case termbox.KeyCtrlS:
		editor.saveFile()
	case termbox.KeyArrowUp:
		editor.scrollUp()
	case termbox.KeyArrowDown:
		editor.scrollDown()
	case termbox.KeyArrowRight:
		editor.scrollRight()
	case termbox.KeyArrowLeft:
		editor.scrollLeft()
	case termbox.KeyEnter:
		editor.insertNewline()
	case termbox.KeyDelete:
		editor.deleteChar()
	case termbox.KeyBackspace2:
		{
			editor.scrollLeft()
			editor.deleteChar()
		}
	case termbox.KeySpace:
		editor.insertRune(' ')
	default:
		editor.insertRune(event.Ch)
	}

	return false
}
