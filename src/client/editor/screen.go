package editor

import (
	"editor-service/node/ot"
	"sync"

	"github.com/nsf/termbox-go"
)

// screen is a passive entity
type Screen struct {
	sync.Mutex
	lines   []string
	offsetX int
	offsetY int
	cursor  Cursor
}

func NewScreen() *Screen {
	return &Screen{
		lines:  []string{""},
		cursor: *newCursor(),
	}
}

func (s *Screen) DrawRemoteEdit(rx <-chan ot.Doc) {
	for doc := range rx {
		s.showDoc(doc)
	}
}

func (s *Screen) showDoc(doc ot.Doc) {
	line := 0
	s.Lock()
	defer s.Unlock()
	s.lines = []string{""}
	for i, digit := range doc {
		if digit == 0x0A {
			line++
			s.lines = append(s.lines, "")
			continue
		}
		s.lines[line] += string(doc[i])
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, height := termbox.Size()

	s.drawText(width, height)
	// s.drawStatus(width, height)

	termbox.SetCursor(s.cursor.x-s.offsetX, s.cursor.y-s.offsetY)
	termbox.Flush()
}

func (s *Screen) drawText(width int, height int) {
	for y := 0; y < height-1; y++ {
		lineY := y + s.offsetY
		if lineY >= len(s.lines) {
			// In this case we are trying to display a line that is not in the lines
			break
		}

		lineContent := s.lines[lineY]
		if s.offsetX < len(lineContent) {
			// We display only the part of the string after the horizontal offset
			displayLine := lineContent[s.offsetX:]
			for x, char := range []rune(displayLine) {
				if x >= width {
					break
				}
				termbox.SetCell(x, y, char, termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}
}

// Callback to handle key events obtained from termbox.
// Returns true if s should be closed
func (s *Screen) OnKeyEvent(event termbox.Event) bool {
	switch event.Key {
	case termbox.KeyCtrlQ:
		// We should exit
		return true
	case termbox.KeyArrowUp:
		s.scrollUp()
	case termbox.KeyArrowDown:
		s.scrollDown()
	case termbox.KeyArrowRight:
		s.scrollRight()
	case termbox.KeyArrowLeft:
		s.scrollLeft()
		// default:
		// 	s.insertRune(event.Ch)
	}

	return false
}

func (s *Screen) scrollRight() {
	_, width := termbox.Size()
	if s.cursor.x+1 == len(s.lines[s.cursor.y]) {
		return
	}
	if s.offsetX > 0 || s.cursor.x >= width {
		s.cursor.moveRight()
		s.offsetX++
		return
	}
	s.cursor.moveRight()
}

func (s *Screen) scrollLeft() {
	if s.cursor.x == 0 && s.cursor.y == 0 {
		return
	}
	if s.cursor.x == 0 {
		s.cursor.goToTheEndOfPreviousLine(len(s.lines[s.cursor.y-1]))
		s.cursor.x++
	}
	if s.offsetX > 0 {
		s.cursor.moveLeft()
		s.offsetX--
		return
	}
	s.cursor.moveLeft()
}

func (s *Screen) scrollUp() {
	if s.cursor.y == 0 {
		return
	}
	if s.offsetY > 0 && s.cursor.y > 0 {
		s.cursor.goToTheEndOfPreviousLine(len(s.lines[s.cursor.y-1]))
		s.offsetY--
		return
	}
	s.cursor.goToTheEndOfPreviousLine(len(s.lines[s.cursor.y-1]))
}

func (s *Screen) scrollDown() {
	if s.cursor.y+1 == len(s.lines) {
		return
	}
	_, height := termbox.Size()
	if s.cursor.y > height {
		s.cursor.goToTheEndOfNextLine(len(s.lines[s.cursor.y+1]))
		s.offsetY++
		return
	}
	s.cursor.goToTheEndOfNextLine(len(s.lines[s.cursor.y+1]))
}

func (s *Screen) CursorBounds() (int, int) {
	prev := 0
	next := len(s.lines[s.cursor.y]) - s.cursor.x
	for _, line := range s.lines[:s.cursor.y] {
		prev += len(line) + 1
	}
	prev += s.cursor.x

	for _, line := range s.lines[s.cursor.y+1:] {
		next += len(line) + 1
	}
	return prev, next
}
