package main

import (
	"fmt"
	"os"

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
	selected map[int]interface{}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}

func tweaks() []tweak {
	return []tweak{
		tweak{
			name:              "Dnf parallel downloads",
			desc:              "Set the number of parallel downloads dnf can do to 10.",
			callback:          func() error { return nil },
			selectedByDefault: true,
		},
		tweak{
			name:              "Remove Fedora Flatpak",
			desc:              "Remove the Fedora Flatpak apps and repository.",
			callback:          func() error { return nil },
			selectedByDefault: false,
		},
		tweak{
			name:              "Set flatpak as default in Gnome Software",
			desc:              "Change the order sources appear in Gnome Software so that flatpak is first.",
			callback:          func() error { return nil },
			selectedByDefault: false,
		},
		tweak{
			name:              "Load i2c-dev and i2c-piix4 kernel modules",
			desc:              "Load needed kernel modules for hardware detection in software like OpenRGB.",
			callback:          func() error { return nil },
			selectedByDefault: true,
		},
		tweak{
			name:              "Install systemd-container",
			desc:              "Install the systemd-container dnf package, mainly with GDM Settings in mind.",
			callback:          func() error { return nil },
			selectedByDefault: true,
		},
		tweak{
			name:              "Fix issue between SELinux and Source games",
			desc:              "Install the systemd-container dnf package, mainly with GDM Settings in mind.",
			callback:          func() error { return nil },
			selectedByDefault: true,
		},
		tweak{
			name:              "Fix issue with big games",
			desc:              "Install the systemd-container dnf package, mainly with GDM Settings in mind.",
			callback:          func() error { return nil },
			selectedByDefault: true,
		},
		tweak{
			name:              "Install non-free p7zip with unrar capacities",
			desc:              "Install the systemd-container dnf package, mainly with GDM Settings in mind.",
			callback:          func() error { return nil },
			selectedByDefault: true,
		},
	}
}

func initialModel() model {
	m := model{
		choices:  []tweak{},
		selected: map[int]interface{}{},
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
			fmt.Println("should run all the tweaks")
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
