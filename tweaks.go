package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
)

func tweaks() []tweak {
	return []tweak{
		tweak{
			name: "Dnf parallel downloads",
			desc: "Set the number of parallel downloads dnf can do to 10.",
			callback: func() error {
				f, err := os.OpenFile("/etc/dnf/dnf.conf", os.O_APPEND|os.O_RDWR, 0644)
				if err != nil {
					return err
				}
				defer f.Close()

				reader := bufio.NewReader(f)
				for {
					line, err := reader.ReadString('\n')
					if err == io.EOF {
						break
					}
					if err != nil {
						return err
					}

					if strings.Contains(line, "max_parallel_downloads") {
						return nil
					}
				}

				_, err = f.WriteString("max_parallel_downloads=10\n")
				if err != nil {
					return err
				}

				return nil
			},
			selectedByDefault: true,
		},
		tweak{
			name: "Remove Fedora flatpak remote",
			desc: "Remove the Fedora Flatpak apps and repository.",
			callback: func() error {
				stdOut, err := exec.Command("flatpak", "remotes").Output()
				if err != nil {
					return err
				}

				if !strings.Contains(string(stdOut), "fedora") {
					return nil
				}

				_, err = exec.Command("flatpak", "remote-delete", "fedora").Output()
				if err != nil {
					return err
				}

				return nil
			},
			selectedByDefault: true,
		},
		tweak{
			name: "Add Flathub flatpak remote",
			desc: "Add the most popular remote for flatpaks.",
			callback: func() error {
				_, err := exec.Command("flatpak", "remote-add", "--if-not-exists", "flathub", "https://dl.flathub.org/repo/flathub.flatpakrepo").Output()
				if err != nil {
					return err
				}

				return nil
			},
			selectedByDefault: true,
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
			selectedByDefault: false,
		},
		tweak{
			name:              "Install systemd-container",
			desc:              "Install the systemd-container dnf package, mainly with GDM Settings in mind.",
			callback:          func() error { return nil },
			selectedByDefault: false,
		},
		tweak{
			name:              "Fix issue between SELinux and Source games",
			desc:              "Install the systemd-container dnf package, mainly with GDM Settings in mind.",
			callback:          func() error { return nil },
			selectedByDefault: false,
		},
		tweak{
			name:              "Fix issue with big games",
			desc:              "Install the systemd-container dnf package, mainly with GDM Settings in mind.",
			callback:          func() error { return nil },
			selectedByDefault: false,
		},
		tweak{
			name:              "Install non-free p7zip with unrar capacities",
			desc:              "Install the systemd-container dnf package, mainly with GDM Settings in mind.",
			callback:          func() error { return nil },
			selectedByDefault: false,
		},
	}
}
