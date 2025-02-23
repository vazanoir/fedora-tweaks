package main

import (
	"bufio"
	"errors"
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
			name: "Set flatpak as the prefered packaging format in Gnome Software",
			desc: "Change the order packaging formats appear in Gnome Software so that flatpak shows first.",
			callback: func() error {
				_, err := exec.Command("runuser", "--user", os.Getenv("SUDO_USER"), "--", "sh", "-c", "gsettings set org.gnome.software packaging-format-preference \"['flatpak:flathub', 'flatpak:flathub-beta', 'rpm']\"").Output()
				if err != nil {
					return err
				}

				return nil
			},
			selectedByDefault: true,
		},
		tweak{
			name: "Load i2c-dev and i2c-piix4 kernel modules",
			desc: "Load needed kernel modules for hardware detection in software like OpenRGB.",
			callback: func() error {
				filePath := "/etc/modules-load.d/i2c.conf"
				_, err := os.Stat(filePath)
				if err == nil {
					return nil
				}

				if !errors.Is(err, os.ErrNotExist) {
					return err
				}

				err = os.WriteFile(filePath, []byte("i2c-dev\ni2c-piix4\n"), 0644)
				if err != nil {
					return err
				}

				return nil
			},
			selectedByDefault: true,
		},
		tweak{
			name: "Install systemd-container, dependency for apps like GDM Settings",
			desc: "This tweak exists because nothing tells you that GDM Settings need that package installed.",
			callback: func() error {
				stdOut, err := exec.Command("dnf", "list", "--installed").Output()
				if err != nil {
					return err
				}
				if strings.Contains(string(stdOut), "systemd-container") {
					return nil
				}

				stdOut, err = exec.Command("dnf", "install", "systemd-container").Output()
				if err != nil {
					return err
				}
				return nil
			},
			selectedByDefault: true,
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
