package tweaks

import (
	"os"
	"os/exec"
)

var preferedPackagingFormat = Tweak{
	Name: "Set flatpak as the prefered packaging format in Gnome Software",
	Desc: "Change the order packaging formats appear in Gnome Software so that flatpak shows first.",
	Callback: func() error {
		_, err := exec.Command(
			"runuser",
			"--user",
			os.Getenv("SUDO_USER"),
			"--",
			"sh",
			"-c",
			"gsettings set org.gnome.software packaging-format-preference \"['flatpak:flathub', 'flatpak:flathub-beta', 'rpm']\"",
		).Output()
		if err != nil {
			return err
		}

		return nil
	},
	SelectedByDefault: true,
	SupportedVersions: []int{41},
}
