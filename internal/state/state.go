package state

import (
	focus "todo.go/internal/ui"

	"fmt"
	"os"
	"strings"
)

func Load(todos *[]string, dones *[]string, filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	for i, line := range strings.Split(string(data), "\n") {
		if len(line) == 0 {
			continue
		}
		foc, it := parseItem(line)
		switch foc {
		case focus.FocusTodo:
			*todos = append(*todos, it)
		case focus.FocusDone:
			*dones = append(*dones, it)
		default:
			fmt.Fprintf(os.Stderr, "%s:%d: Error: ill-formed item line: %s\n", filepath, i, line)
			os.Exit(1)
		}
	}
}

func Save(todos *[]string, dones *[]string, filepath string) {
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("Error when create file: %s", err)
	}
	for _, todo := range *todos {
		file.WriteString(fmt.Sprintf("TODO: %s\n", todo))
	}
	for _, done := range *dones {
		file.WriteString(fmt.Sprintf("DONE: %s\n", done))
	}
	file.Close()
}

func parseItem(line string) (focusState int, text string) {
	const TodoPrefix string = "TODO: "
	const DonePrefix string = "DONE: "
	if strings.HasPrefix(line, TodoPrefix) {
		focusState = focus.FocusTodo
		text = line[len(TodoPrefix):]
		return focusState, text
	}

	if strings.HasPrefix(line, DonePrefix) {
		focusState = focus.FocusDone
		text = line[len(DonePrefix):]
		return focusState, text
	}

	return -1, ""
}
