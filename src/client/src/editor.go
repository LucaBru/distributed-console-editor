package editor

import (
	"fmt"
	"os"
	"strings"

	"github.com/nsf/termbox-go"
)

type Editor struct {
	buffer []string // These are the lines of text
	cursorX int // Position of cursor on X axis
	cursorY int // Position of cursor on axis Y
	offsetX int // Scroll on X axis
	offsetY int // Scroll on Y axis
	filename string // File where text will be saved
	modified bool
	statusMsg string
	backgroundColor termbox.Attribute
	foregroundColor termbox.Attribute
	statusBackgroundColor termbox.Attribute
	statusForegroundColor termbox.Attribute
}

func NewEditor() *Editor {
	return &Editor{
		buffer: []string{""},
		backgroundColor: termbox.ColorDefault,
		foregroundColor: termbox.ColorDefault,
		statusBackgroundColor: termbox.ColorBlack,
		statusForegroundColor: termbox.ColorWhite,
	}
}

// Draw editor content to the terminal
func (editor *Editor) Draw() {
	// We clear the current text on the screen
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, height := termbox.Size()

	editor.drawText(width, height)
	editor.drawStatus(width, height)

	termbox.SetCursor(editor.cursorX - editor.offsetX, editor.cursorY - editor.offsetY)
	termbox.Flush()
}

func (editor *Editor) drawText(width int, height int) {
	for y := 0; y < height - 1; y++ {
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
	statusLine := fmt.Sprintf(" %s - %d lines %s", editor.filename, len(editor.buffer), map[bool]string{true: "[modified]", false: ""}[editor.modified])
	if editor.statusMsg != "" {
		statusLine = editor.statusMsg
	}

	// We fill the status line with spaces
	for x := 0; x < width; x++ {
		termbox.SetCell(x, height - 1, ' ', editor.statusBackgroundColor, editor.statusForegroundColor)
	}

	// Draw the status
	for x, char := range []rune(statusLine) {
		if x >= width {
			break
		}
		termbox.SetCell(x, height - 1, char, editor.statusBackgroundColor, editor.statusForegroundColor)
	}

}

func (editor *Editor) insertRune(char rune) {
	line := []rune(editor.buffer[editor.cursorY])
	width, _ := termbox.Size()
	if editor.cursorX > len(line) {
		// If cursor is over the end of the current line we insert some spaces to fill the void
		for i := len(line); i < editor.cursorX; i++ {
			line = append(line, ' ')
		}
	}

	// Now we can insert the character
	line = append(line[:editor.cursorX], append([]rune{char}, line[editor.cursorX:]...)...)
	if len(line) > width {
		editor.offsetX++
	}
	editor.buffer[editor.cursorY] = string(line)
	editor.cursorX++
	editor.modified = true
}

// InsertNewline inserts a newline at the current cursor position
func (editor *Editor) insertNewline() {
	if editor.cursorY >= len(editor.buffer) {
		editor.buffer = append(editor.buffer, "")
	} else {
		line := editor.buffer[editor.cursorY]
		beforeCursor := ""
		afterCursor := ""
		
		if editor.cursorX < len(line) {
			beforeCursor = line[:editor.cursorX]
			afterCursor = line[editor.cursorX:]
		} else {
			beforeCursor = line
		}
		
		editor.buffer[editor.cursorY] = beforeCursor
		
		// Insert a new line after the current one
		editor.buffer = append(
			editor.buffer[:editor.cursorY+1],
			append([]string{afterCursor}, editor.buffer[editor.cursorY+1:]...)...,
		)
	}
	
	editor.cursorY++
	editor.cursorX = 0
	editor.modified = true
}

// DeleteChar deletes the character at the current cursor position
func (editor *Editor) deleteChar() {
	if editor.cursorX == 0 && editor.cursorY == 0 {
		return
	}
	
	if editor.cursorX > 0 {
		// Delete character before cursor
		line := []rune(editor.buffer[editor.cursorY])
		if editor.cursorX <= len(line) {
			line = append(line[:editor.cursorX-1], line[editor.cursorX:]...)
			editor.buffer[editor.cursorY] = string(line)
			editor.cursorX--
		}
	} else {
		// Backspace at the beginning of a line, join with previous line
		if editor.cursorY > 0 {
			prevLineLen := len(editor.buffer[editor.cursorY-1])
			editor.buffer[editor.cursorY-1] += editor.buffer[editor.cursorY]
			editor.buffer = append(editor.buffer[:editor.cursorY], editor.buffer[editor.cursorY+1:]...)
			editor.cursorY--
			editor.cursorX = prevLineLen
		}
	}
	
	editor.modified = true
}

// SaveFile saves the current buffer to a file
func (editor *Editor) saveFile() {
	content := strings.Join(editor.buffer, "\n")
	err := os.WriteFile(editor.filename, []byte(content), 0644)
	if err != nil {
		editor.setStatus("Error saving file: " + err.Error())
	} else {
		editor.modified = false
		editor.setStatus(fmt.Sprintf("Saved %s (%d bytes)", editor.filename, len(content)))
	}
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
	case termbox.KeyEnter:
		editor.insertNewline()
	case termbox.KeyBackspace, termbox.KeyBackspace2, termbox.KeyDelete:
		editor.deleteChar()
	case termbox.KeySpace:
		editor.insertRune(' ')
	default:
		editor.insertRune(event.Ch)
	}

	return false
}