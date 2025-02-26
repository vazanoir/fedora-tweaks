package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func tweaks() []tweak {
	return []tweak{
		tweak{
			name: "Dnf parallel downloads",
			desc: "Set the number of parallel downloads dnf can do to 10.",
			callback: func() error {
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
			selectedByDefault: true,
		},
		tweak{
			name: "Remove Fedora flatpak remote",
			desc: "Remove the Fedora Flatpak apps and repository.",
			callback: func() error {
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
			selectedByDefault: true,
		},
		tweak{
			name: "Add Flathub flatpak remote",
			desc: "Add the most popular remote for flatpaks.",
			callback: func() error {
				_, err := exec.Command("flatpak", "remote-add", "--if-not-exists", "flathub", "https://dl.flathub.org/repo/flathub.flatpakrepo").Output()
				if err != nil {
					return err
				}

				return nil
			},
			selectedByDefault: true,
		},
		tweak{
			name: "Set flatpak as the prefered packaging format in Gnome Software",
			desc: "Change the order packaging formats appear in Gnome Software so that flatpak shows first.",
			callback: func() error {
				_, err := exec.Command("runuser", "--user", os.Getenv("SUDO_USER"), "--", "sh", "-c", "gsettings set org.gnome.software packaging-format-preference \"['flatpak:flathub', 'flatpak:flathub-beta', 'rpm']\"").Output()
				if err != nil {
					return err
				}

				return nil
			},
			selectedByDefault: true,
		},
		tweak{
			name: "Load i2c-dev and i2c-piix4 kernel modules",
			desc: "Load needed kernel modules for hardware detection in software like OpenRGB.",
			callback: func() error {
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
			selectedByDefault: true,
		},
		tweak{
			name: "Install systemd-container, dependency for apps like GDM Settings",
			desc: "This tweak exists because nothing tells you that GDM Settings need that package installed.",
			callback: func() error {
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
			selectedByDefault: true,
		},
		tweak{
			name: "Fix issue between SELinux and Source games",
			desc: "Some Source games weren't made with the best security practices, and have some sound assets\n      diffusion block by SELinux, this tweak lower the security for a better experience.",
			callback: func() error {
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
			selectedByDefault: false,
		},
		tweak{
			name: "Increase vm.max_map_count to 16 GB",
			desc: "Some applications and games (like Red Dead Redemption 2 or Star Citizen) crash because\n      of this value being too low. This tweak increase it to 16 GB, don't use this tweak if you\n      have less than that amount in RAM.",
			callback: func() error {
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
			selectedByDefault: false,
		},
		tweak{
			name: "Install non-free p7zip with unrar capacities",
			desc: "Fedora removed the rar capabilities of the shipped p7zip package, this tweak install\n      an older version that has that capability.",
			callback: func() error {
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
						if !strings.Contains(line, "p7zip p7zip-plugins") {
							foundExcludeLine = true
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
			selectedByDefault: false,
		},
	}
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
