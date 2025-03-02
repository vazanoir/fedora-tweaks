package tweaks

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var nonFreeP7zip = Tweak{
	Name: "Install non-free p7zip with unrar capacities",
	Desc: "Fedora removed the rar capabilities of the shipped p7zip package, this tweak install\n      an older version that has that capability.",
	SelectedByDefault: false,
	SupportedVersions: []int{41},
	Callback: func() error {
		// removing already installed p7zip
		stdOut, err := exec.Command("dnf", "list", "--installed").Output()
		if err != nil {
			return fmt.Errorf("listing packaging: %v", err)
		}
		if strings.Contains(string(stdOut), "p7zip") {
			stdOut, err = exec.Command("dnf", "remove", "-y", "p7zip", "p7zip-plugins").Output()
			if err != nil {
				return fmt.Errorf("removing current p7zip: %v", err)
			}
		}

		// p7zip
		u, err := url.Parse("https://github.com/ttys3/fedora-rpm-p7zip/releases/download/16.02/p7zip-16.02-24.fc37.x86_64.rpm")
		if err != nil {
			return fmt.Errorf("parsing url: %v", err)
		}

		fpath, err := downloadFromGithub(u)
		if err != nil {
			return fmt.Errorf("downloading: %v", err)
		}

		stdOut, err = exec.Command("dnf", "install", "-y", fpath).Output()
		if err != nil {
			return fmt.Errorf("installing: %v", err)
		}

		// p7zip-plugins
		u, err = url.Parse("https://github.com/ttys3/fedora-rpm-p7zip/releases/download/16.02/p7zip-plugins-16.02-24.fc37.x86_64.rpm")
		if err != nil {
			return fmt.Errorf("parsing url: %v", err)
		}

		fpath, err = downloadFromGithub(u)
		if err != nil {
			return fmt.Errorf("downloading: %v", err)
		}

		stdOut, err = exec.Command("dnf", "install", "-y", fpath).Output()
		if err != nil {
			return fmt.Errorf("installing: %v", err)
		}

		// add exception to dnf.conf
		f, err := os.OpenFile("/etc/dnf/dnf.conf", os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		reader := bufio.NewReader(f)
		fileContent := ""
		foundExcludeLine := false
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			if strings.Contains(line, "exclude=") {
				foundExcludeLine = true

				if !strings.Contains(line, "p7zip p7zip-plugins") {
					line = line[:len(line)-1] + " p7zip p7zip-plugins\n"
				}
			}

			fileContent += line
		}

		if !foundExcludeLine {
			fileContent += "exclude=p7zip p7zip-plugins\n"
		}

		err = f.Truncate(0)
		if err != nil {
			return err
		}

		_, err = f.Seek(0, 0)
		if err != nil {
			return err
		}

		_, err = f.WriteString(fileContent)
		if err != nil {
			return err
		}

		return nil
	},
}

func downloadFromGithub(u *url.URL) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 && resp.StatusCode <= 399 {
		redirectUrl, err := resp.Location()
		if err != nil {
			return "", err
		}

		req.URL = redirectUrl
		resp, err = client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
	}

	fpath := fmt.Sprintf("/tmp/%v", filepath.Base(u.Path))
	out, err := os.Create(fpath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return fpath, nil
}
