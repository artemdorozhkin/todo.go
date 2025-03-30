package main

import (
	"fmt"
	"log"
	"slices"

	"github.com/rthornton128/goncurses"
)

const (
	RegularPair     int16 = 1
	HighlightPair   int16 = 2
	CursorInvisible byte  = 0
)

// class UI
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

// end class UI

// class Focus
const (
	FocusTodo = iota
	FocusDone = iota
)

type Focus struct {
	Focus int
}

func (f *Focus) Switch() {
	switch f.Focus {
	case FocusTodo:
		f.Focus = FocusDone
	case FocusDone:
		f.Focus = FocusTodo
	}
}

// end class Focus

func listUp(listCurr *Id) {
	*listCurr = max(*listCurr-1, 0)
}

func listDown(list *[]string, listCurr *Id) {
	if *listCurr+1 < len(*list) {
		*listCurr++
	}
}

func main() {
	stdscr, err := goncurses.Init()
	goncurses.Echo(false)
	goncurses.Cursor(CursorInvisible)

	if err != nil {
		log.Fatal(err)
	}
	defer goncurses.End()

	goncurses.StartColor()
	if err := goncurses.InitPair(RegularPair, goncurses.C_WHITE, goncurses.C_BLACK); err != nil {
		log.Fatal(err)
	}

	if err := goncurses.InitPair(HighlightPair, goncurses.C_BLACK, goncurses.C_WHITE); err != nil {
		log.Fatal(err)
	}

	quit := false

	todos := []string{
		"Write the todo app",
		"Buy a bread",
		"Make a cup of coffee",
	}
	todoCurr := 0

	dones := []string{
		"Start learning go",
		"Have a breakfast",
		"Make a cup of coffee",
	}
	doneCurr := 0

	var ui Ui = Ui{Stdscr: stdscr}
	focus := Focus{}

	for !quit {
		stdscr.Erase()
		ui.Begin(0, 0)
		{
			switch focus.Focus {
			case FocusTodo:
				ui.Label("[TODO] DONE ", RegularPair)
				ui.Label("------------", RegularPair)
				ui.BeginList(todoCurr)
				for index, todo := range todos {
					ui.ListElement(fmt.Sprintf("- [ ] %s", todo), index)
				}
				ui.EndList()
			case FocusDone:
				ui.Label(" TODO [DONE]", RegularPair)
				ui.Label("------------", RegularPair)
				ui.BeginList(doneCurr)
				for index, done := range dones {
					ui.ListElement(fmt.Sprintf("- [x] %s", done), index)
				}
				ui.EndList()
			}
		}
		ui.End()

		stdscr.Refresh()

		key := string(rune(stdscr.GetChar()))

		switch key {
		case "q":
			quit = true
		case "w":
			switch focus.Focus {
			case FocusTodo:
				listUp(&todoCurr)
			case FocusDone:
				listUp(&doneCurr)
			}
		case "s":
			switch focus.Focus {
			case FocusTodo:
				listDown(&todos, &todoCurr)
			case FocusDone:
				listDown(&dones, &doneCurr)
			}
		case "\n":
			switch focus.Focus {
			case FocusTodo:
				if todoCurr < len(todos) {
					dones = append(dones, todos[todoCurr])
					todos = slices.Delete(todos, todoCurr, todoCurr+1)
				}
			case FocusDone:
				if doneCurr < len(dones) {
					todos = append(todos, dones[doneCurr])
					dones = slices.Delete(dones, doneCurr, doneCurr+1)
				}
			}
		case "\t":
			focus.Switch()
		default:
			continue
		}
	}
}
