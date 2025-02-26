package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type tweak struct {
	name              string
	desc              string
	callback          func() error
	selectedByDefault bool
}

type model struct {
	choices  []tweak
	cursor   int
	selected map[int]any
}

func main() {
	// check if root
	if os.Geteuid() != 0 {
		fmt.Println(errFmt("this program requires root privileges"))
		os.Exit(1)
	}

	// check for system updates
	cmd := exec.Command("dnf", "check-upgrade")
	if err := cmd.Start(); err != nil {
		fmt.Println(errFmt(err.Error()))
		os.Exit(2)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok && exiterr.ExitCode() != 0 {
			fmt.Println(errFmt("please update your system"))
			os.Exit(3)
		} else {
			fmt.Println(errFmt(err.Error()))
		}
	}

	// start the program
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(errFmt(err.Error()))
		os.Exit(4)
	}
}

func errFmt(err string) error {
	return fmt.Errorf("\033[0;31merror\033[0m: %v", err)
}

func initialModel() model {
	m := model{
		choices:  []tweak{},
		selected: map[int]any{},
	}

	for i, t := range tweaks() {
		m.choices = append(m.choices, t)

		if t.selectedByDefault {
			m.selected[i] = nil
		}
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "k":
			if m.cursor > 0 {
				m.cursor -= 1
			}
		case "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor += 1
			}
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = tweak{}
			}
		case "r":
			for i := range m.selected {
				err := m.choices[i].callback()
				if err != nil {
					lowerName := strings.ToLower(m.choices[i].name)
					fmt.Printf("%v (%v)", errFmt(err.Error()), lowerName)
					os.Exit(100 + i)
				}
			}
			fmt.Printf("\033[0;32m\nTweaks successfully applied!\033[0m")
			os.Exit(0)
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Select the wanted tweaks:\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.name)
		s += fmt.Sprintf("\033[30m      %s\n\033[0m", choice.desc)
	}

	s += "\n[q] quit   [j] down   [k] up   [space] select   [r] execute selected tweaks\n"

	return s
}
