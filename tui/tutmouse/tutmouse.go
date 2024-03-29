package tutmouse

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

// https://github.com/charmbracelet/bubbletea/blob/master/examples/mouse/main.go

func Run() {
	p := tea.NewProgram(model{}, tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	mouseEvent tea.MouseEvent
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s := msg.String(); s == "ctrl+c" || s == "q" || s == "esc" {
			return m, tea.Quit
		}
		return m, tea.Printf("(msg:%s, type:%s)", msg, msg.Type)

	case tea.MouseMsg:
		return m, tea.Printf("(X: %d, Y: %d) %s", msg.X, msg.Y, tea.MouseEvent(msg))
	}

	return m, nil
}

func (m model) View() string {
	s := "Do mouse stuff. When you're done press q to quit.\n"

	return s
}
