package main

import (
	"os"
)

func tweaks() []tweak {
	return []tweak{
		tweak{
			name: "Dnf parallel downloads",
			desc: "Set the number of parallel downloads dnf can do to 10.",
			callback: func() error {
                f, err := os.Open("/etc/dnf/dnf.conf")
                if err != nil {
                    return err
                }
                defer f.Close()
				return nil
			},
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
