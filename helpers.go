package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

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

func bold(str string) string {
	return fmt.Sprintf("\033[1m%v\033[0m", str)
}

func red(str string) string {
	return fmt.Sprintf("\033[0;31m%v\033[0m", str)
}

func green(str string) string {
	return fmt.Sprintf("\033[0;32m%v\033[0m", str)
}

func lightgrey(str string) string {
	return fmt.Sprintf("\033[37m%v\033[0m", str)
}
