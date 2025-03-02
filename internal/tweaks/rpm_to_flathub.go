package tweaks

import (
	"fmt"
	"os/exec"
	"strings"
)

var packageList = map[string]string{
	"firefox":             "org.mozilla.firefox",
	"gnome-calendar":      "org.gnome.Calendar",
	"gnome-text-editor":   "org.gnome.TextEditor",
	"gnome-contacts":      "org.gnome.Contacts",
	"snapshot":            "org.gnome.Snapshot",
	"gnome-weather":       "org.gnome.Weather",
	"gnome-clocks":        "org.gnome.clocks",
	"gnome-maps":          "org.gnome.Maps",
	"mediawriter":         "org.fedoraproject.MediaWriter",
	"libreoffice-writer":  "org.libreoffice.LibreOffice",
	"libreoffice-calc":    "org.libreoffice.LibreOffice",
	"libreoffice-impress": "org.libreoffice.LibreOffice",
	"totem":               "org.gnome.Totem",
	"gnome-calculator":    "org.gnome.Calculator",
	"simple-scan":         "org.gnome.SimpleScan",
	"gnome-boxes":         "org.gnome.Boxes",
	"rhythmbox":           "org.gnome.Rhythmbox3",
	"baobab":              "org.gnome.baobab",
	"gnome-connections":   "org.gnome.Connections",
	"evince":              "org.gnome.Evince",
	"loupe":               "org.gnome.Loupe",
	"gnome-characters":    "org.gnome.Characters",
	"gnome-logs":          "org.gnome.Logs",
	"gnome-font-viewer":   "org.gnome.font-viewer",
}

var rpmToFlathub = Tweak{
	Name: "Swap all default rpm apps for Flathub's flatpaks",
	Desc: "Install the flatpak version from Flathub of all default apps and remove the rpm ones.",
	SelectedByDefault: true,
	SupportedVersions: []int{41},
	Callback: func() error {

		for dnfPkg, flatpakPkg := range packageList {
			stdOut, err := exec.Command("dnf", "list", "--installed").Output()
			if err != nil {
				return fmt.Errorf("listing installed packages: %v", err)
			}

			if !strings.Contains(string(stdOut), dnfPkg) {
				continue
			}

			_, err = exec.Command("dnf", "remove", "-y", dnfPkg).Output()
			if err != nil {
				return fmt.Errorf("removing %v: %v", dnfPkg, err)
			}

			_, err = exec.Command("flatpak", "install", "-y", "flathub", flatpakPkg).Output()
			if err != nil {
				return fmt.Errorf("installing %v: %v", flatpakPkg, err)
			}
		}

		_, err := exec.Command("dnf", "autoremove", "-y").Output()
		if err != nil {
			return fmt.Errorf("autoremoving: %v", err)
		}

		return nil
	},
}
