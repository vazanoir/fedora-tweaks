package tweaks

import (
	"bufio"
	"io"
	"os"
	"strings"
)

var vmMaxMapCount = Tweak{
	Name: "Increase vm.max_map_count to 16G",
	Desc: "Some applications and games (like Red Dead Redemption 2 or Star Citizen) crash because\n      of this value being too low. This tweak increase it to 16G, don't use this tweak if you\n      have less than that amount in RAM.",
	Callback: func() error {
		f, err := os.OpenFile("/etc/sysctl.conf", os.O_APPEND|os.O_RDWR, 0644)
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

			if strings.Contains(line, "vm.max_map_count") {
				return nil
			}
		}

		_, err = f.WriteString("vm.max_map_count = 16777216\n")
		if err != nil {
			return err
		}

		return nil
	},
	SelectedByDefault: false,
	SupportedVersions: []int{41},
}
