package tweaks

import (
	"os/exec"
	"strings"
)

var selinuxAllowExecheap = Tweak{
	Name: "Fix issue between SELinux and Source games",
	Desc: "Some Source games weren't made with the best security practices, and have some sound assets\n      diffusion block by SELinux, this tweak lower the security for a better experience.",
	Callback: func() error {
		stdOut, err := exec.Command("getsebool", "allow_execheap").Output()
		if err != nil {
			return err
		}
		if strings.Contains(string(stdOut), "on") {
			return nil
		}

		stdOut, err = exec.Command("setsebool", "-P", "allow_execheap", "1").Output()
		if err != nil {
			return err
		}

		return nil
	},
	SelectedByDefault: false,
	SupportedVersions: []int{41},
}
