package tweaks

import (
	"os/exec"
	"strings"
)

var removeFedoraRemote = Tweak{
	Name: "Remove Fedora flatpak remote",
	Desc: "Remove the Fedora Flatpak apps and repository.",
	SelectedByDefault: true,
	SupportedVersions: []int{41},
	Callback: func() error {
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
}
