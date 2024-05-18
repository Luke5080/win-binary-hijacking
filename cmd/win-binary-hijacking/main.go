package main

import (
	"fmt"

	"github.com/Luke5080/win-binary-hijacking/internal/enumwin"
)

func main() {
	fmt.Println("HIJACKED!")

	chanBack := enumwin.EnumServ()

	fmt.Println("Can Modify the following services:")

	for val := range chanBack {
		fmt.Println((*val).Name)

		if (*val).StartName == "LocalSystem" {
			fmt.Println((*val).Name + "LOCALSYSTEM")
		}
	}
}
