package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type tweak struct {
	name              string
	desc              string
	callback          func() error
	selectedByDefault bool
	supportedVersions []int
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

func initialModel() model {
	m := model{
		choices:  []tweak{},
		selected: map[int]any{},
	}

	tweaks := getTweaks()

	// get the version of user's fedora
	version := 0
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		fmt.Println(errFmt(err.Error()))
		os.Exit(8)
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		pair := strings.Split(line, "=")
		if len(pair) != 2 {
			continue
		}

		key := pair[0]
		value := pair[1]

		if key == "VERSION_ID" {
			version, err = strconv.Atoi(value)
			if err != nil {
				fmt.Println(errFmt(err.Error()))
				os.Exit(9)
			}
			break
		}
	}

	// filter unsupported tweaks for the current fedora version
	supportedTweaks := []tweak{}
	for _, tweak := range tweaks {
		if slices.Contains(tweak.supportedVersions, version) {
			supportedTweaks = append(supportedTweaks, tweak)
		}
	}

	// move selectedByDefault tweaks at the top
	sort.Slice(supportedTweaks, func(i, j int) bool {
		return supportedTweaks[i].selectedByDefault && !supportedTweaks[j].selectedByDefault
	})

	// setup the model
	for i, tweak := range supportedTweaks {
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
			fmt.Printf(green("\nTweaks successfully applied!\n"))
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

		s += bold(fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.name))
		s += lightgrey(fmt.Sprintf("      %s\n", choice.desc))
	}

	s += "\n[q] quit   [j] down   [k] up   [space] select   [a] select all   [n] select none   [enter] apply selected tweaks\n"

	return s
}
