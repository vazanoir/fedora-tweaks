package main

import (
	"fmt"
)

func errFmt(err string) error {
	return fmt.Errorf("%v: %v", red("error"), err)
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
