package tweaks

import (
	"os/exec"
	"strings"
)

var installSystemdContainer = Tweak{
	Name: "Install systemd-container, dependency for apps like GDM Settings",
	Desc: "This tweak exists because nothing tells you that GDM Settings need that package installed.",
	SelectedByDefault: true,
	SupportedVersions: []int{41},
	Callback: func() error {
		stdOut, err := exec.Command("dnf", "list", "--installed").Output()
		if err != nil {
			return err
		}
		if strings.Contains(string(stdOut), "systemd-container") {
			return nil
		}

		stdOut, err = exec.Command("dnf", "install", "-y", "systemd-container").Output()
		if err != nil {
			return err
		}
		return nil
	},
}
