package editor

type Cursor struct {
	x int
	y int
}

func newCursor() *Cursor {
	return &Cursor{
		x: 0,
		y: 0,
	}
}

func (cursor *Cursor) moveUp() {
	if cursor.y > 0 {
		cursor.y --
	}
}

func (cursor *Cursor) moveDown() {
	cursor.y++
}

func (cursor *Cursor) moveRight() {
	cursor.x++
}

func (cursor *Cursor) moveLeft() {
	if (cursor.x > 0) {
		cursor.x--
	}
}

func (cursor *Cursor) returnToTheBeginOfTheLine() {
	cursor.x = 0
}

func (cursor *Cursor) goToTheEndOfPreviousLine(lineLength int) {
	if cursor.y > 0 {
		cursor.y--
		cursor.x = lineLength
	}
}

func (cursor *Cursor) goToTheEndOfNextLine(lineLength int) {
	cursor.y++
	cursor.x = lineLength
}