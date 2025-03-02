package tweaks

import "os/exec"

var addFlathubRemote = Tweak{
	Name: "Add Flathub flatpak remote",
	Desc: "Add the most popular remote for flatpaks.",
	SelectedByDefault: true,
	SupportedVersions: []int{41},
	Callback: func() error {
		_, err := exec.Command(
			"flatpak",
			"remote-add",
			"--if-not-exists",
			"flathub",
			"https://dl.flathub.org/repo/flathub.flatpakrepo",
		).Output()
		if err != nil {
			return err
		}

		return nil
	},
}
