package main

import (
	"todo.go/internal/state"
	list "todo.go/internal/todo"
	"todo.go/internal/ui"

	"fmt"
	"log"
	"os"

	"github.com/rthornton128/goncurses"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "Usage: todo-go <file-path>\n")
		fmt.Fprint(os.Stderr, "Error: file path is not provided\n")
		os.Exit(1)
	}

	filepath := os.Args[1]

	todos := []string{}
	todoCurr := 0

	dones := []string{}
	doneCurr := 0

	state.Load(&todos, &dones, filepath)

	stdscr, err := goncurses.Init()
	goncurses.Echo(false)
	goncurses.Cursor(ui.CursorInvisible)

	if err != nil {
		log.Fatal(err)
	}
	defer goncurses.End()

	goncurses.StartColor()
	if err := goncurses.InitPair(ui.RegularPair, goncurses.C_WHITE, goncurses.C_BLACK); err != nil {
		log.Fatal(err)
	}

	if err := goncurses.InitPair(ui.HighlightPair, goncurses.C_BLACK, goncurses.C_WHITE); err != nil {
		log.Fatal(err)
	}

	quit := false

	var window ui.Ui = ui.Ui{Stdscr: stdscr}
	status := ui.Status{}

	for !quit {
		stdscr.Erase()
		window.Begin(0, 0)
		{
			switch status.Focus {
			case ui.FocusTodo:
				window.Label("[TODO] DONE ", ui.RegularPair)
				window.Label("------------", ui.RegularPair)
				window.BeginList(todoCurr)
				for index, todo := range todos {
					window.ListElement(fmt.Sprintf("- [ ] %s", todo), index)
				}
				window.EndList()
			case ui.FocusDone:
				window.Label(" TODO [DONE]", ui.RegularPair)
				window.Label("------------", ui.RegularPair)
				window.BeginList(doneCurr)
				for index, done := range dones {
					window.ListElement(fmt.Sprintf("- [x] %s", done), index)
				}
				window.EndList()
			}
		}
		window.End()

		stdscr.Refresh()

		key := string(rune(stdscr.GetChar()))

		switch key {
		case "q":
			quit = true
		case "w":
			switch status.Focus {
			case ui.FocusTodo:
				list.Up(&todoCurr)
			case ui.FocusDone:
				list.Up(&doneCurr)
			}
		case "s":
			switch status.Focus {
			case ui.FocusTodo:
				list.Down(&todos, &todoCurr)
			case ui.FocusDone:
				list.Down(&dones, &doneCurr)
			}
		case "\n":
			switch status.Focus {
			case ui.FocusTodo:
				list.Transfer(&dones, &todos, &todoCurr)
			case ui.FocusDone:
				list.Transfer(&todos, &dones, &doneCurr)
			}
		case "\t":
			status.Switch()
		default:
			continue
		}
	}
	state.Save(&todos, &dones, filepath)
}
