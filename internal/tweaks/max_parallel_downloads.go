package tweaks

import (
	"bufio"
	"io"
	"os"
	"strings"
)

var maxParallelDownloads = Tweak{
	Name: "Dnf parallel downloads",
	Desc: "Set the number of parallel downloads dnf can do to 10.",
	SelectedByDefault: true,
	SupportedVersions: []int{41},
	Callback: func() error {
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
}
