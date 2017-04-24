package procedures

import (
	"fmt"
	"vmkube/model"
	"os"
	"strings"
	"vmkube/utils"
	"errors"
	"os/exec"
	"time"
)

func DefineCloudServerCommand(server model.ProjectCloudServer) []string {
	name, driver, uuid, options := server.Name, server.Driver, server.Id, server.Options
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
	return command
}

func DefineLocalServerCommand(server model.ProjectServer, imagePath string) ([]string, int) {
	driver, disksize, cpus, noshare := server.Driver, server.DiskSize, server.Cpus, server.NoShare
	name, memory, osname := server.Name, server.Memory, server.OSType
	options, engine, swarm, uuid:= server.Options, server.Engine, server.Swarm, server.Id
	DiskResize := 0
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
	if osname != "rancheros" ||  "virtualbox" != strings.ToLower(driver) {
		if disksize > 0 {
			disksize_str := utils.IntToString(disksize)
			command = append( command,  "--"+strings.ToLower(driver)+"-disk-size")
			command = append( command,  disksize_str)
		}
	} else {
		fmt.Printf("RANCHEROS - Disksize %dGB ignored ....\n", disksize)
		DiskResize = disksize
	}
	if memory > 0 {
		memory_str := utils.IntToString(memory)
		if "vmwarefusion" == strings.ToLower(driver) {
			command = append( command,  "--"+strings.ToLower(driver)+"-memory-size")
		} else {
			command = append( command,  "--"+strings.ToLower(driver)+"-memory")
		}
		command = append( command,  memory_str)
	}
	if "virtualbox" == strings.ToLower(driver) && noshare {
		command = append( command,  "--virtualbox-no-share")
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
	command = append( command,  "file://" + strings.Replace(imagePath, string(os.PathSeparator), "/", len(imagePath)))
	// Custom Driver Options
	for _,option := range options {
		command = append( command,  "--"+strings.ToLower(driver)+"-"+option[0]+"")
		command = append( command,  option[1])
	}
	command = append( command,  name + "-" + uuid)
	return command, DiskResize
}

func (machine *DockerMachine) CreateCloudServer(server model.ProjectCloudServer, commandPipe chan MachineMessage) {
	name, uuid:= server.Name, server.Id
	var command []string = DefineCloudServerCommand(server)
	bytes, err := executeSyncCommand(command)
	state := Machine_State_Running
	if err != nil {
		state = Machine_State_Error
	}
	time.Sleep(3000)
	machineName := name + "-" + uuid
	var message string = ""
	cmd := exec.Command("docker-machine", "stop", machineName)
	message += fmt.Sprintf("Stopping docker machine : %s\n", machineName)
	bytes, err = cmd.CombinedOutput()
	message += fmt.Sprintf("%s\n",bytes)
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: CreateCloud,
		Error: err,
		Result: string(bytes),
		Supply: message,
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachine)  CreateServer(server model.ProjectServer, commandPipe chan MachineMessage) {
	name, uuid, osname, osver:= server.Name, server.Id, server.OSType, server.OSVersion
	path, success := DownloadISO(osname, osver)
	if !success {
		commandPipe <- MachineMessage{
			Complete: true,
			Cmd: []string{},
			Project: machine.Project,
			Infra: machine.Infra,
			Operation: CreateServer,
			Error: errors.New(fmt.Sprintf("Unable to download or recover iso image for OS %s v.%s", osname, osver)),
			Result: "",
			State: Machine_State_Error,
			InstanceId: machine.InstanceId,
			IsCloud: machine.IsCloud,
		}
		return
	}
	var command []string
	var diskSize int
	command, diskSize = DefineLocalServerCommand(server, path)
	bytes, err := executeSyncCommand(command)
	var message string = ""
	time.Sleep(3000)
	machineName := name + "-" + uuid
	cmd := exec.Command("docker-machine", "stop", machineName)
	message += fmt.Sprintf("Stopping docker machine : %s\n", machineName)
	bytesArray, _ := cmd.CombinedOutput()
	message += fmt.Sprintf("%s\n",bytesArray)
	if diskSize > 0 {
		message += fmt.Sprintf("Resizing disk to %sGB", diskSize)
		file := model.HomeFolder() + string(os.PathSeparator) + ".docker" + string(os.PathSeparator) + "machine" + string(os.PathSeparator) + "machines" +
			string(os.PathSeparator) +  name + "-" + uuid + string(os.PathSeparator) + "disk.vmdk"
		cmd := exec.Command("vmware-vdiskmanager", "-x", fmt.Sprintf("%dGB", diskSize), file)
		bytesArray, _ := cmd.CombinedOutput()
		message += fmt.Sprintf("%s\n",bytesArray)
	}
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: CreateServer,
		Error: err,
		Result: string(bytes),
		Supply: message,
		State: Machine_State_Stopped,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachine)  RemoveServer(name string, id string, commandPipe chan MachineMessage) {
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "rm")
	command = append( command,  "-f")
	command = append( command,  name + "-" + id)
	fmt.Printf("Running Delete command :  '%s'\n", strings.Join(command, " "))
	bytes, err := executeSyncCommand(command)
	state := Machine_State_Running
	if err != nil {
		state = Machine_State_Stopped
	}
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: DestroyServer,
		Error: err,
		Result: string(bytes),
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachine)  StopServer(name string, id string, commandPipe chan MachineMessage) {
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "stop")
	command = append( command,  name + "-" + id)
	fmt.Printf("Running Stop command :  '%s'\n", strings.Join(command, " "))
	bytes, err := executeSyncCommand(command)
	state := Machine_State_Running
	if err != nil {
		state = Machine_State_Stopped
	}
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: StopServer,
		Error: err,
		Result: string(bytes),
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachine)  StartServer(name string, id string, commandPipe chan MachineMessage) {
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "start")
	command = append( command,  name + "-" + id)
	fmt.Printf("Running Start command :  '%s'\n", strings.Join(command, " "))
	bytes, err := executeSyncCommand(command)
	state := Machine_State_Running
	if err != nil {
		state = Machine_State_Stopped
	}
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: StartServer,
		Error: err,
		Result: string(bytes),
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachine)  RestartServer(name string, id string, commandPipe chan MachineMessage) {
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "restart")
	command = append( command,  name + "-" + id)
	fmt.Printf("Running Restart command :  '%s'\n", strings.Join(command, " "))
	bytes, err := executeSyncCommand(command)
	state := Machine_State_Running
	if err != nil {
		state = Machine_State_Stopped
	}
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: RestartServer,
		Error: err,
		Result: string(bytes),
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachine)  ServerStatus(name string, id string, commandPipe chan MachineMessage) {
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "status")
	command = append( command,  name + "-" + id)
	fmt.Printf("Running Status command :  '%s'\n", strings.Join(command, " "))
	bytesArray, err :=  executeSyncCommand(command)
	state := Machine_State_None
	state = GetStateFromMachineAnswer(string(bytesArray))
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: StatusServer,
		Error: err,
		Result: string(bytesArray),
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachine)  ServerEnv(name string, id string, commandPipe chan MachineMessage) {
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "env")
	command = append( command,  name + "-" + id)
	fmt.Printf("Running Environment command :  '%s'\n", strings.Join(command, " "))
	bytes, err := executeSyncCommand(command)
	state := Machine_State_Running
	if err != nil {
		state = Machine_State_Stopped
	}
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: ServerEnvironment,
		Error: err,
		Result: string(bytes),
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachine)  ServerInspect(name string, id string, commandPipe chan MachineMessage) {
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "inspect")
	command = append( command,  name + "-" + id)
	fmt.Printf("Running Inspect command :  '%s'\n", strings.Join(command, " "))
	bytes, err := executeSyncCommand(command)
	state := Machine_State_Running
	if err != nil {
		state = Machine_State_Stopped
	}
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: ServerInspect,
		Error: err,
		Result: string(bytes),
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachine)  ServerIPAddress(name string, id string, commandPipe chan MachineMessage) {
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "ip")
	command = append( command,  name + "-" + id)
	fmt.Printf("Running IP Address command :  '%s'\n", strings.Join(command, " "))
	bytes, err := executeSyncCommand(command)
	state := Machine_State_Running
	if err != nil {
		state = Machine_State_Stopped
	}
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: ServerIPAddress,
		Error: err,
		Result: string(bytes),
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}
