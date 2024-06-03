# Hijacked!

### Project Description
Hijacked is an experimental script which I created to exploit weak/misconfigured services on Windows which can be used to perform privilege escalation. The script finds services which can be modified by authenticated users, and grabs as much information as possible for the weak service. Once found, the user can replace the service binary with a reverse shell which will connect back to the your host system, or you can replace the service's binary with your own payload which must be installed on the system before using the script. This is my first personal project written in Go, so I would appreciate as much constructive criticism as possible in order to better my code quality/logic.  

### How To Use
The script assumes you already have access to your victim Windows system, and will begin enumerating the services available, similiar to scripts such as **WinPeas**. In order to get this script onto your target system, you can clone this repository locally to your host system, and set up a python http server which will allow you to download the script from the target system using utilities such as `curl`. The main purpose of this script is to perform **privilege escalation**, with the main goal being finding weak services which start at the highest privileges (LocalSystem) and changing the binary path to a reverse shell, which will allow for a System32 shell to be executed from the service, which will connect back to your host system. If service(s) are found which can be modified by users in the Authenicated Users group, the script will display them to the user and also highlight important information, for example if the service starts as LocalSystem (important for privilege escalation) and if the current user has some control rights on the service (if the user can stop/start the service). If the script finds services which can be modified, the user will have the option to either replace the service's binary with a reverse shell which is available in this repository, or the user can use their own custom binary (e.g. A reverse shell payload from msfvenom). The script then attempts to replace the service binary, and will restart the service if the user has permissions to do so, or will set the start mode of the script to AUTO_START to allow for it to be started on the next reboot.  

Once the script is installed on your target system, you can follow the following steps:  
1. Change into the repository folder:  
`cd win-binary-hijacking`  
2. Run the main executable:  
`bin\win-binary-hijacking.exe`  
3. The script will attempt to find all services on the system which can be modified by users in the Authenticated Users group. If no services are found which can be modified, the script will output the message:  
`No services found which you can modify`  
and will exit. If the script does find services which can be modified, it will displayed in a menu format, with subsequent information found about the service. For example, if the `wuauserv` service can be modified by an authenticated user, and the user has permission to stop the service, the script will output the following:  
`1. wuauserv CAN STOP STARTS AS LocalSystem`  
The script will the following options for the user:  
```
What binary would you like to replace the service binary with?  
1. Reverse Shell
2. Custom Binary
```  
The reverse shell executable which comes with this repository is available in the `internal/malbinaries/revshell.exe`. **Even though this reverse shell is available, I still recommend using a reverse shell payload built with msfvenom, as it will provide a more sophisticated and robust reverse shell which can be used with process migration if the shell exits. The reverse shell available is a simple TCP reverse shell which will connect back to your host system, and offers no other options.** Once you have chosen your option, enter the menu number of the option (1 or 2). If you choose option one, you will be prompted to enter your host IP:  
`What is the host IP for the reverse shell?`  
And the port number:  
`What is the port number for the host?`  
If you choose option 2, the script will prompt for the path of the executable which will replace the service binary (it must be already installed on the target Windows system):  
`Please enter FULL path of your custom payload:`  
4. The script will then attempt to replace the choosen service's binary path with the new malicious binary. If succesful, it will check whether you have appropriate permissions to start/stop the service, so that the new malicious binary can be executed. If you do have these permissions, the script will prompt whether you would to restart the service:  
`Can Start/Stop Service...`  
`Restart Service? [y/n]`  
Choosing `n` will exit the script, where as `y` will stop then start the service. If the script can start/stop the service sucesfully, the script will exit and the service will have started with your new malicious binary. If you do not have succifient permissions to start/stop the service, the script will modify the service for the start type to be set as `AUTO_START`, which will require for you to restart the target system for the service to be started.  

### Personal Notes:  
The source code for the reverse shell used in this script can be found [here](https://github.com/Luke5080/go-reverse-shell)  
