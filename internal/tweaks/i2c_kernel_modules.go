package tweaks

import (
	"errors"
	"os"
)

var i2cKernelModules = Tweak{
	Name: "Load i2c-dev and i2c-piix4 kernel modules",
	Desc: "Load needed kernel modules for hardware detection in software like OpenRGB.",
	SelectedByDefault: true,
	SupportedVersions: []int{41},
	Callback: func() error {
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
}
