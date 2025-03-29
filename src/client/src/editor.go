package editor

import (
	"fmt"
	"os"
	"strings"

	"github.com/nsf/termbox-go"
)

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
}

func NewEditor() *Editor {
	return &Editor{
		buffer:                []string{""},
		backgroundColor:       termbox.ColorDefault,
		foregroundColor:       termbox.ColorDefault,
		statusBackgroundColor: termbox.ColorBlack,
		statusForegroundColor: termbox.ColorWhite,
		filename: "untitled.txt",
	}
}

// Draw editor content to the terminal
func (editor *Editor) Draw() {
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
	statusLine := fmt.Sprintf(" %s - %d lines %s", editor.filename, len(editor.buffer), map[bool]string{true: "[modified]", false: ""}[editor.modified])
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

func (editor *Editor) insertRune(char rune) {
	line := []rune(editor.buffer[editor.cursor.y])
	width, _ := termbox.Size()
	if editor.cursor.x > len(line) {
		// If cursor is over the end of the current line we insert some spaces to fill the void
		for i := len(line); i < editor.cursor.x; i++ {
			line = append(line, ' ')
		}
	}

	// Now we can insert the character
	line = append(line[:editor.cursor.x], append([]rune{char}, line[editor.cursor.x:]...)...)
	if len(line) > width {
		// In this case since we are writing over the available space we scroll horizontally
		editor.offsetX++
	}
	editor.buffer[editor.cursor.y] = string(line)
	editor.cursor.moveRight()
	editor.modified = true
}

// InsertNewline inserts a newline at the current cursor position
func (editor *Editor) insertNewline() {
	_, height := termbox.Size()
	if editor.cursor.y >= len(editor.buffer) {
		editor.buffer = append(editor.buffer, "")
	} else {
		line := editor.buffer[editor.cursor.y]
		beforeCursor := ""
		afterCursor := ""

		if editor.cursor.x < len(line) {
			beforeCursor = line[:editor.cursor.x]
			afterCursor = line[editor.cursor.x:]
		} else {
			beforeCursor = line
		}

		editor.buffer[editor.cursor.y] = beforeCursor

		// Insert a new line after the current one
		editor.buffer = append(
			editor.buffer[:editor.cursor.y+1],
			append([]string{afterCursor}, editor.buffer[editor.cursor.y+1:]...)...,
		)
	}

	editor.cursor.moveDown()
	if editor.cursor.y > height-2 {
		editor.offsetY++
	}
	editor.cursor.returnToTheBeginOfTheLine()
	editor.modified = true
}

// DeleteChar deletes the character at the current cursor position
func (editor *Editor) deleteChar() {
	if editor.cursor.x == 0 && editor.cursor.y == 0 {
		return
	}

	if editor.cursor.x > 0 {
		// Delete character before cursor
		line := []rune(editor.buffer[editor.cursor.y])
		if editor.cursor.x <= len(line) {
			line = append(line[:editor.cursor.x-1], line[editor.cursor.x:]...)
			editor.buffer[editor.cursor.y] = string(line)
			editor.cursor.moveLeft()
		}
	} else {
		// Backspace at the beginning of a line, join with previous line
		if editor.cursor.y > 0 {
			prevLineLen := len(editor.buffer[editor.cursor.y-1])
			editor.buffer[editor.cursor.y-1] += editor.buffer[editor.cursor.y]
			editor.buffer = append(editor.buffer[:editor.cursor.y], editor.buffer[editor.cursor.y+1:]...)
			editor.cursor.goToTheEndOfPreviousLine(prevLineLen)
		}
	}

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
	err := os.WriteFile(file.Name(), []byte(content), 0644)
	if err != nil {
		editor.setStatus("Error saving file: " + err.Error())
	} else {
		editor.modified = false
		editor.setStatus(fmt.Sprintf("Saved %s (%d bytes)", editor.filename, len(content)))
	}
}

func (editor *Editor) scrollRight() {
	if editor.cursor.x < len(editor.buffer[editor.cursor.y]) {
		editor.cursor.moveRight()
		editor.offsetX++
	}
}

func (editor *Editor) scrollLeft() {
	if editor.offsetX > 0 && editor.cursor.x > 0 {
		editor.cursor.moveLeft()
		editor.offsetX--
	}
}

func (editor *Editor) scrollUp() {
	if editor.offsetY > 0 && editor.cursor.y > 0 {
		editor.cursor.goToTheEndOfPreviousLine(len(editor.buffer[editor.cursor.y - 1]))
		editor.offsetY--
	}
}

func (editor *Editor) scrollDown() {
	if editor.cursor.y < len(editor.buffer)-1 {
		editor.cursor.goToTheEndOfNextLine(len(editor.buffer[editor.cursor.y+1]))
		editor.offsetY++
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
	case termbox.KeyBackspace, termbox.KeyBackspace2, termbox.KeyDelete:
		editor.deleteChar()
	case termbox.KeySpace:
		editor.insertRune(' ')
	default:
		editor.insertRune(event.Ch)
	}

	return false
}
