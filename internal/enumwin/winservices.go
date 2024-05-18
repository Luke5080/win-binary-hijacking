package enumwin

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type WeakServ struct {
	// Struct to hold weak system service attributes
	// including their name, their start mode (manual, auto) and
	// their path to the service's binary
	Name      string
	StartMode string
	BinPath   string
	StartName string
}

func makeWeakServ(serv string, c chan *WeakServ, wg *sync.WaitGroup) {
	// Function to obtain more details about services that we can modify
	// Grab information such as their Start mode andthe path name to the
	// service binary. Create a weakServ struct for each weak
	// service and write to the channel
	defer wg.Done()

	cmdFormat := fmt.Sprintf("Get-CimInstance -Class Win32_Service | Select-Object -Property Name, StartMode, PathName, StartName | Where {$_.Name -eq %q} | ForEach-Object { $_.StartMode + '_' + $_.PathName + '_' + $_.StartName }", serv)

	cmd := exec.Command("powershell", "-command", cmdFormat)

	out, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
	}

	err = cmd.Start()

	if err != nil {
		fmt.Println(err)
	}

	scanner := bufio.NewScanner(out)

	for scanner.Scan() {
		line := scanner.Text()
		lineReFactor := strings.Split(line, "_")

		serv := &WeakServ{serv, lineReFactor[0], lineReFactor[1], lineReFactor[2]}

		c <- serv
	}
}

func checkServPerms(serv string, c chan string, wg *sync.WaitGroup) {

	// Create regex pattern to search in service SDDL.
	// To enable privilege escalation from an account
	// who is in not in a  user group, the service needs
	// to have write perms as an AUTHENTICATED USER.
	// Regex pattern should look like the following:
	pattern := `\(A;;[^;]*W[^;]*\;;\;AU\)`

	compileRegex, err := regexp.Compile(pattern)

	if err != nil {
		fmt.Println("Cannot compile regex", err)
		return
	}

	defer wg.Done()

	// Get SDDL of service and grab the output
	out, err := exec.Command("sc", "sdshow", serv).Output()

	if err != nil {
		fmt.Printf("Error %v: with service: %s\n", err, serv)
	}

	// Match it against the regex pattern. If the pattern
	// matches, append service name to the channel
	if compileRegex.MatchString(string(out)) {
		fmt.Printf("Can Modify Service: %s\n", serv)
		c <- serv
	}

}

func getServ() []string {
	// Function to grab all services on the user system
	// Done by running a PS command and append the results
	// to a slice of strings
	fmt.Println("Finding all services on system..")

	allServ := []string{}

	cmd := exec.Command("powershell", "-command", "Get-CimInstance -Class Win32_Service | Select-Object -ExpandProperty Name")

	// Pipe stdout to the out var
	out, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println("Could not execute command")
	}

	defer out.Close()

	err = cmd.Start()

	if err != nil {
		fmt.Println("Could not run command")
	}

	// Create scanner to read command output stored in `out`
	scanner := bufio.NewScanner(out)

	for scanner.Scan() {
		line := scanner.Text()

		line = strings.Replace(line, " ", "", -1)

		// Cleanup empty lines here and dont append
		// them to the systemServ slice
		if len(line) == 0 {

			continue

		} else { // If the service is not junk append it to the systemServ slice and return it
			allServ = append(allServ, line)
		}
	}
	fmt.Printf("Found %d services on target machine\n", len(allServ))
	return allServ
}

func EnumServ() chan *WeakServ {
	// getLocalSystemServ returns slice of strings and is held in localSystemServ
	localServ := getServ()

	wg := new(sync.WaitGroup)

	// Make channel with a capacity of length of the localSystemServ slice
	canModify := make(chan string, len(localServ))

	// Pass each service in checkServPerms which only appends
	// weak services to the canModify channel
	for _, serv := range localServ {
		wg.Add(1)
		fmt.Printf("Checking Service permissions: %s\n", serv)
		go checkServPerms(serv, canModify, wg)
	}

	wg.Wait()
	close(canModify)

	serviceStruct := make(chan *WeakServ, len(canModify))

	wgServ := new(sync.WaitGroup)

	for val := range canModify {
		wgServ.Add(1)
		fmt.Printf("Grabbing information for service: %s\n", val)
		go makeWeakServ(val, serviceStruct, wgServ)
	}

	wgServ.Wait()
	close(serviceStruct)

	return serviceStruct
}
