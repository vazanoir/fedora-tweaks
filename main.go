package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
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

	t := tweaks()

	// move selectedByDefault tweaks at the top
	sort.Slice(t, func(i, j int) bool {
		return t[i].selectedByDefault && !t[j].selectedByDefault
	})

	for i, tweak := range t {
		m.choices = append(m.choices, tweak)

		if tweak.selectedByDefault {
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
		case "a":
			for i := range m.choices {
				m.selected[i] = tweak{}
			}
		case "n":
			for i := range m.choices {
				delete(m.selected, i)
			}
		case "enter":
			for i := range m.selected {
				err := m.choices[i].callback()
				if err != nil {
					lowerName := strings.ToLower(m.choices[i].name)
					fmt.Printf("%v (%v)", errFmt(err.Error()), lowerName)
					os.Exit(100 + i)
				}
			}
			fmt.Printf("\n\033[0;32mTweaks successfully applied!\033[0m\n")
			os.Exit(0)
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Select tweaks:\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("\033[1m%s [%s] %s\033[0m\n", cursor, checked, choice.name)
		s += fmt.Sprintf("\033[37m      %s\033[0m\n", choice.desc)
	}

	s += "\n[q] quit   [j] down   [k] up   [space] select   [a] select all   [n] select none   [enter] apply selected tweaks\n"

	return s
}
