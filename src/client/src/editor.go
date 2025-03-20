package editor

import(
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
}

func NewEditor() *Editor {
	return &Editor{
		buffer: []string{""},
	}
}

// Draw editor content to the terminal
func (editor *Editor) Draw() {
	// We clear the current text on the screen
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, height := termbox.Size()

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
				termbox.SetCell(x, y, char, termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}

	termbox.SetCursor(editor.cursorX - editor.offsetX, editor.cursorY - editor.offsetY)
	termbox.Flush()
}

func (editor *Editor) InsertRune(char rune) {
	line := []rune(editor.buffer[editor.cursorY])
	if editor.cursorX > len(line) {
		// If cursor is over the end of the current line we insert some spaces to fill the void
		for i := len(line); i < editor.cursorX; i++ {
			line = append(line, ' ')
		}
	}

	// Now we can insert the character
	line = append(line[:editor.cursorX], append([]rune{char}, line[editor.cursorX:]...)...)
	editor.buffer[editor.cursorY] = string(line)
	editor.cursorX++
	editor.modified = true
}

// Callback to handle key events obtained from termbox.
// Returns true if editor should be closed
func (editor *Editor) OnKeyEvent(event termbox.Event) bool {
	switch event.Key {
	case termbox.KeyCtrlQ:
		// We should exit
		return true
	case termbox.KeySpace:
		editor.InsertRune(' ')
	default:
		editor.InsertRune(event.Ch)
	}

	return false
}