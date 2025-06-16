package ui

type Status struct {
	Focus int
}

func (s *Status) Switch() {
	switch s.Focus {
	case FocusTodo:
		s.Focus = FocusDone
	case FocusDone:
		s.Focus = FocusTodo
	}
}
