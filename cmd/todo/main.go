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
	// doneCurr := 0

	var ui Ui = Ui{Stdscr: stdscr}

	for !quit {
		stdscr.Erase()
		ui.Begin(0, 0)
		{
			ui.Label("TODO:", RegularPair)
			ui.BeginList(todoCurr)
			for index, todo := range todos {
				ui.ListElement(fmt.Sprintf("- [ ] %s", todo), index)
			}
			ui.EndList()

			ui.Label("----------------------------------------", RegularPair)

			ui.Label("DONE:", RegularPair)
			ui.BeginList(0)
			for index, done := range dones {
				ui.ListElement(fmt.Sprintf("- [x] %s", done), index+1)
			}
			ui.EndList()
		}
		ui.End()

		stdscr.Refresh()

		key := string(rune(stdscr.GetChar()))

		switch key {
		case "q":
			quit = true
		case "w":
			todoCurr = max(todoCurr-1, 0)
		case "s":
			if todoCurr+1 < len(todos) {
				todoCurr++
			}
		case "\n":
			if todoCurr < len(todos) {
				dones = append(dones, todos[todoCurr])
				todos = slices.Delete(todos, todoCurr, todoCurr+1)
			}
		default:
			continue
		}
	}
}
