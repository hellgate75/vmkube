package procedures

import (
	"fmt"
	"vmkube/model"
	"log"
	"bufio"
	"os"
	"strings"
	"vmkube/utils"
	"os/exec"
	"errors"
)

func DownloadISO(machineType string, version string) (string, bool) {
	machineAction, error := model.GetMachineAction(machineType)
	if error == nil {
		if ! machineAction.Check(version) {
			fmt.Println("Machine",machineType,"Version",version,"not present, downloading from internet...")
			downloaded := machineAction.Download(version)
			fmt.Println("Machine",machineType,"Version",version,"dowanloaded:",downloaded)
			return machineAction.Path(version), downloaded
		} else {
			fmt.Println("Machine",machineType,"Version",version,"already dowanloaded...")
			return machineAction.Path(version), true
		}
	} else {
		log.Fatal("Machine",machineType,"not found!!","-","error:", error)
		return  "", false
	}
}

func RequestConfirmation(reason string) bool {
	text := ""
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stdout, "%s.Confirm [y/n] :", reason)
	text, _ = reader.ReadString('\n')
	if text != "Y" && text != "y" && text != "N" && text != "n" {
		fmt.Fprintln(os.Stdout, "Current text is not allowed :", text)
		return  RequestConfirmation(reason)
	}
	return false
}

func CreateCloudServer(server model.ProjectCloudServer) ([]byte, error) {
	name, driver, hostname, uuid, options, roles := server.Name, server.Driver, server.Hostname, server.Id, server.Options, server.Roles
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "create")
	command = append( command,  "-d")
	command = append( command,  driver)
	// Custom Driver Options
	for _,option := range options {
		command = append( command,  "--"+strings.ToLower(driver)+"-"+option[0]+"")
		command = append( command,  option[1])
	}
	command = append( command,  name + "-" + uuid)
	
	fmt.Println(os.Stdout,"Running Create for hostname: " + hostname + " - Roles "+strings.Join(roles, ",")+" - command : '" + strings.Join(command, " ") + "'")
	return  executeCommand(command)
}


func CreateServer(server model.ProjectServer) ([]byte, error) {
	driver, disksize, cpus, hostname, noshare := server.Driver, server.DiskSize, server.Cpus, server.Hostname, server.NoShare
	name, memory, osname, osver, roles := server.Name, server.Memory, server.OSType, server.OSVersion, server.Roles
	options, engine, swarm, uuid:= server.Options, server.Engine, server.Swarm, server.Id
	path, success := DownloadISO(osname, osver)
	if !success {
		return  []byte{}, errors.New("Unable to discover os : " + osname + ":" + osver)
	}
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "create")
	command = append( command,  "-d")
	command = append( command,  driver)
	if cpus > 0 {
		cpus_str := utils.IntToString(cpus)
		command = append( command,  "--"+strings.ToLower(driver)+"-cpu-count")
		command = append( command,  cpus_str)
	}
	if disksize > 0 {
		disksize_str := utils.IntToString(disksize)
		command = append( command,  "--"+strings.ToLower(driver)+"-disk-size")
		command = append( command,  disksize_str)
	}
	if memory > 0 {
		memory_str := utils.IntToString(memory)
		command = append( command,  "--"+strings.ToLower(driver)+"-memory")
		command = append( command,  memory_str)
	}
	if "virtualbox" == strings.ToLower(driver) && noshare {
		command = append( command,  "virtualbox-no-share")
	}
	// Custom Driver Options
	for _,option := range options {
		command = append( command,  "--"+strings.ToLower(driver)+"-"+option[0]+"")
		command = append( command,  option[1])
	}

	// Docker Engine Specific commands (RancherOS doesn't allow Engine customizations)
	if osname != "rancheros" {
		for _,option := range engine.Environment {
			command = append( command,  "--engine-env")
			command = append( command,  option)
		}

		for _,option := range engine.InsecureRegistry {
			command = append( command,  "--engine-insecure-registry")
			command = append( command,  option)
		}

		for _,option := range engine.RegistryMirror {
			command = append( command,  "--engine-registry-mirror")
			command = append( command,  option)
		}

		for _,option := range engine.Labels {
			command = append( command,  "--engine-label")
			command = append( command,  option)
		}

		for _,option := range engine.Options {
			command = append( command,  "--engine-opt")
			command = append( command,  option)
		}

		if "" != strings.TrimSpace(engine.InstallURL) {
			command = append( command,  "--engine-install-url")
			command = append( command,  strings.TrimSpace(engine.InstallURL))
		}

		if "" != strings.TrimSpace(engine.StorageDriver) {
			command = append( command,  "--engine-storage-driver")
			command = append( command,  strings.TrimSpace(engine.StorageDriver))
		}
	}

	// Docker Swarm Specific commands  (RancherOS Swarm feature do not tested yet)
	if swarm.Enabled && osname != "rancheros" {
		command = append( command,  "--swarm")
		if swarm.IsMaster {
			command = append( command,  "--swarm-master")
		}
		if "" != swarm.DiscoveryToken {
			command = append( command,  "--swarm-discovery")
			command = append( command,  swarm.DiscoveryToken)
		}
		if swarm.UseExperimental {
			command = append( command,  "--swarm-experimental")
		}
		if swarm.UseAddress {
			command = append( command,  "--swarm-addr")
		}
		if "" != swarm.Host {
			command = append( command,  "--swarm-host")
			command = append( command,  swarm.Host)
		}
		if "" != swarm.Image {
			command = append( command,  "--swarm-image")
			command = append( command,  swarm.Image)
		}
		if "" != swarm.Strategy {
			command = append( command,  "--swarm-strategy")
			command = append( command,  swarm.Strategy)
		}
		for _,option := range swarm.JoinOpts {
			command = append( command,  "--swarm-opt")
			command = append( command,  option)
		}
		for _,option := range swarm.TLSSan {
			if strings.Index(option, " ") > 0 {
				optionCouple := strings.Split(option, " ")
				command = append( command,  optionCouple[0])
				command = append( command,  strings.Join(optionCouple[1:], " "))
			} else  {
				command = append( command,  option)
			}
		}
	}
	command = append( command,  "--"+strings.ToLower(driver)+"-boot2docker-url")
	command = append( command,  "file://" + path)
	command = append( command,  name + "-" + uuid)

	fmt.Println(os.Stdout,"Running Create for hostname: " + hostname + " - Roles "+strings.Join(roles, ",")+" - command : '" + strings.Join(command, " ") + "'")
	return  executeCommand(command)
}

func RemoveServer(name string, id string) ([]byte, error) {
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "rm")
	command = append( command,  "-f")
	command = append( command,  name + "-" + id)
	fmt.Println(os.Stdout,"Running Delete command : '" + strings.Join(command, " ") + "'")
	return  executeCommand(command)
}

func ServerStatus(name string, id string) ([]byte, error) {
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "status")
	command = append( command,  name + "-" + id)
	fmt.Println(os.Stdout,"Running Status command : '" + strings.Join(command, " ") + "'")
	return  executeCommand(command)
}

func ServerEnv(name string, id string) ([]byte, error) {
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "env")
	command = append( command,  name + "-" + id)
	fmt.Println(os.Stdout,"Running Environment command : '" + strings.Join(command, " ") + "'")
	return  executeCommand(command)
}

func ServerInspect(name string, id string) ([]byte, error) {
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "inspect")
	command = append( command,  name + "-" + id)
	fmt.Println(os.Stdout,"Running Inspect command : '" + strings.Join(command, " ") + "'")
	return  executeCommand(command)
}

func ServerIPAddr(name string, id string) ([]byte, error) {
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "ip")
	command = append( command,  name + "-" + id)
	fmt.Println(os.Stdout,"Running Inspect command : '" + strings.Join(command, " ") + "'")
	return  executeCommand(command)
}

func executeCommand(command []string) ([]byte, error) {
	//cmd := exec.Command(command[0], command[1:]...)
	cmd := exec.Command(command[0], command[1:]...)
	stdoutStderr, err := cmd.CombinedOutput()
	//if err != nil {
	//	log.Fatal(err)
	//}
	return  stdoutStderr, err
}