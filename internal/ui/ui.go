package ui

import "github.com/rthornton128/goncurses"

const (
	RegularPair     int16 = 1
	HighlightPair   int16 = 2
	CursorInvisible byte  = 0
)

type Id = int

type Ui struct {
	Stdscr   *goncurses.Window
	ListCurr *Id
	Row      int
	Col      int
}

func (u *Ui) Begin(row int, col int) {
	u.Row = row
	u.Col = col
}

func (u *Ui) Label(text string, pair int16) {
	u.Stdscr.Move(u.Row, u.Col)
	u.Stdscr.AttrOn(goncurses.ColorPair(pair))
	u.Stdscr.Print(text)
	u.Stdscr.AttrOff(goncurses.ColorPair(pair))
	u.Row++
}

func (u *Ui) BeginList(id Id) {
	if u.ListCurr != nil {
		panic("Nested lists are not allowed!")
	}

	u.ListCurr = &id
}

func (u *Ui) ListElement(label string, id Id) bool {
	if u.ListCurr == nil {
		panic("Not allowed to create list elements outside off lists")
	}
	idCurr := *u.ListCurr
	var pair int16 = RegularPair
	if idCurr == id {
		pair = HighlightPair
	}
	u.Label(label, pair)

	return false
}

func (u *Ui) EndList() {
	u.ListCurr = nil
}

func (u *Ui) End() {
}
