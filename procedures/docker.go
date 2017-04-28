package procedures

import (
	"fmt"
	"vmkube/model"
	"os"
	"strings"
	"vmkube/utils"
	"os/exec"
	"time"
	"errors"
)

func DockerMachineDefineCloudMachineCommand(machine model.CloudMachine) []string {
	name, driver, uuid, options := machine.Name, machine.Driver, machine.Id, machine.Options
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

func DockerMachineDefineLocalMachineCommand(machine model.LocalMachine, imagePath string) ([]string, int) {
	driver, disksize, cpus, noshare := machine.Driver, machine.DiskSize, machine.Cpus, machine.NoShare
	name, memory, osname := machine.Name, machine.Memory, machine.OSType
	options, engine, swarm, uuid:= machine.Options, machine.Engine, machine.Swarm, machine.Id
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
		if utils.NO_COLORS {
			fmt.Printf("RANCHEROS - Disksize %dGB ignored ....\n", disksize)
		}
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

func DockerMachineInterruptSignal(machine *DockerMachineExecutor, commandPipe chan MachineMessage, operation MachineOperation, state MachineState, message string, err error){
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: []string{},
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: operation,
		Error: err,
		Result: message,
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachineExecutor)  CreateCloudMachine(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd) {
	var name, uuid string
	if machine.NewInfra {
		name, uuid= machine.CMachine.Name, machine.CMachine.Id
	} else {
		name, uuid= machine.CInstance.Name, machine.CInstance.MachineId
	}
	var command []string = DockerMachineDefineCloudMachineCommand(machine.CMachine)
	if machine.Control.Interrupt {
		DockerMachineInterruptSignal(machine, commandPipe,StartMachine,Machine_State_Stopped,"Iterrupted",errors.New("Requested interruption"))
		return
	}
	cmd := executeSyncCommand(command)
	commandChannel <- cmd
	bytes, err := cmd.CombinedOutput()
	state := Machine_State_Running
	if err != nil {
		state = Machine_State_Error
	}
	time.Sleep(3000)
	machineName := name + "-" + uuid
	var message string = ""
	var json string = ""
	var ipAddress string = ""
	cmd1 := executeSyncCommand([]string{"docker-machine", "inspect", machineName})
	if machine.Control.Interrupt {
		DockerMachineInterruptSignal(machine, commandPipe,StartMachine,Machine_State_Stopped,"Iterrupted",errors.New("Requested interruption"))
		return
	}
	commandChannel <- cmd1
	bytes1, err1 := cmd1.CombinedOutput()
	if err1 != nil {
		message += fmt.Sprintf("Inspecting docker machine : %s\n", machineName)
		message += err1.Error() + "\n"
	} else {
		json = fmt.Sprintf("%s\n",bytes1)
	}
	machine.CInstance.InspectJSON = json
	if machine.Control.Interrupt {
		DockerMachineInterruptSignal(machine, commandPipe,StartMachine,Machine_State_Stopped,"Iterrupted",errors.New("Requested interruption"))
		return
	}
	cmd2 := executeSyncCommand([]string{"docker-machine", "ip", machineName})
	commandChannel <- cmd2
	bytes2, err2 := cmd2.CombinedOutput()
	if err2 != nil {
		message += fmt.Sprintf("Getting IPAddress from docker machine : %s\n", machineName)
		message += err2.Error() + "\n"
	} else {
		ipAddress = fmt.Sprintf("%s\n",bytes2)
	}
	if machine.Control.Interrupt {
		DockerMachineInterruptSignal(machine, commandPipe,StartMachine,Machine_State_Stopped,"Iterrupted",errors.New("Requested interruption"))
		return
	}
	cmd3 := executeSyncCommand([]string{"docker-machine", "stop", machineName})
	commandChannel <- cmd1
	message += fmt.Sprintf("Stopping docker machine : %s\n", machineName)
	bytesArray, _ := cmd3.CombinedOutput()
	message += fmt.Sprintf("%s\n",bytesArray)
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
		IPAddress: ipAddress,
		InspectJSON: json,
	}
}

func (machine *DockerMachineExecutor)  CreateMachine(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd) {
	var name, uuid, osname, osver string
	if machine.NewInfra {
		name, uuid, osname, osver = machine.Machine.Name, machine.Machine.Id, machine.Machine.OSType, machine.Machine.OSVersion
	} else {
		name, uuid, osname, osver = machine.Instance.Name, machine.Instance.MachineId, machine.Instance.OSType, machine.Instance.OSVersion
	}
	log, path, err := DownloadISO(osname, osver)
	if err != nil {
		commandPipe <- MachineMessage{
			Complete: true,
			Cmd: []string{},
			Project: machine.Project,
			Infra: machine.Infra,
			Operation: CreateMachine,
			Error: err,
			Result: log,
			State: Machine_State_Error,
			InstanceId: machine.InstanceId,
			IsCloud: machine.IsCloud,
		}
		return
	}
	var command []string
	var diskSize int
	command, diskSize = DockerMachineDefineLocalMachineCommand(machine.Machine, path)
	if machine.Control.Interrupt {
		DockerMachineInterruptSignal(machine, commandPipe,StartMachine,Machine_State_Stopped,"Iterrupted",errors.New("Requested interruption"))
		return
	}
	// Simulate error with 20% probability of un-success ...
	//if rand.Int() % 5 == 0 {
	//	time.Sleep(3*time.Second)
	//	DockerMachineInterruptSignal(machine, commandPipe,StartMachine,Machine_State_Stopped,"Iterrupted",errors.New("Requested interruption"))
	//	return
	//}
	cmd := executeSyncCommand(command)
	commandChannel <- cmd
	bytes, err := cmd.CombinedOutput()
	machineName := name + "-" + uuid
	var message string = log
	time.Sleep(3000)
	var json string = ""
	var ipAddress string = ""
	
	/*
	Fake code
	 */
	//message = "Welldone!!"
	//json = "{\"messge\":\"OK\"}"
	//ipAddress="1.1.1.1"
	//time.Sleep(20*time.Second)
	//bytes := []byte{}
	//var err error
	//if rand.Int() % 5 == 0 {
	//	err = errors.New("Cazzo di errore ....")
	//}
	if machine.Control.Interrupt {
		DockerMachineInterruptSignal(machine, commandPipe,StartMachine,Machine_State_Stopped,"Iterrupted",errors.New("Requested interruption"))
		return
	}
	cmd1 := executeSyncCommand([]string{"docker-machine", "inspect", machineName})
	commandChannel <- cmd1
	bytes1, err1 := cmd1.CombinedOutput()
	if err1 != nil {
		message += fmt.Sprintf("Inspecting docker machine : %s\n", machineName)
		message += err1.Error() + "\n"
	} else {
		json = fmt.Sprintf("%s\n",bytes1)
	}
	machine.Instance.InspectJSON = json
	if machine.Control.Interrupt {
		DockerMachineInterruptSignal(machine, commandPipe,StartMachine,Machine_State_Stopped,"Iterrupted",errors.New("Requested interruption"))
		return
	}
	cmd2 := executeSyncCommand([]string{"docker-machine", "ip", machineName})
	commandChannel <- cmd2
	bytes2, err2 := cmd2.CombinedOutput()
	if err2 != nil {
		message += fmt.Sprintf("Getting IPAddress from docker machine : %s\n", machineName)
		message += err2.Error() + "\n"
	} else {
		ipAddress = fmt.Sprintf("%s\n",bytes2)
	}
	cmd3 := executeSyncCommand([]string{"docker-machine", "stop", machineName})
	if machine.Control.Interrupt {
		DockerMachineInterruptSignal(machine, commandPipe,StartMachine,Machine_State_Stopped,"Iterrupted",errors.New("Requested interruption"))
		return
	}
	commandChannel <- cmd3
	message += fmt.Sprintf("Stopping docker machine : %s\n", machineName)
	bytesArray, _ := cmd3.CombinedOutput()
	message += fmt.Sprintf("%s\n",bytesArray)
	if machine.Control.Interrupt {
		DockerMachineInterruptSignal(machine, commandPipe,StartMachine,Machine_State_Stopped,"Iterrupted",errors.New("Requested interruption"))
		return
	}
	if diskSize > 0 {
		message += fmt.Sprintf("Resizing disk to %sGB", diskSize)
		file := model.HomeFolder() + string(os.PathSeparator) + ".docker" + string(os.PathSeparator) + "machine" + string(os.PathSeparator) + "machines" +
			string(os.PathSeparator) +  machineName + string(os.PathSeparator) + "disk.vmdk"
		cmd4 := executeSyncCommand([]string{"vmware-vdiskmanager", "-x", fmt.Sprintf("%dGB", diskSize), file})
		commandChannel <- cmd4
		bytesArray, _ := cmd4.CombinedOutput()
		message += fmt.Sprintf("%s\n",bytesArray)
	}
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: CreateMachine,
		Error: err,
		Result: string(bytes),
		Supply: message,
		State: Machine_State_Stopped,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
		IPAddress: ipAddress,
		InspectJSON: json,
	}
}

func (machine *DockerMachineExecutor)  RemoveMachine(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd) {
	var name, id string
	if machine.IsCloud {
		if machine.NewInfra {
			name, id= machine.CMachine.Name, machine.CMachine.Id
		} else {
			name, id= machine.CInstance.Name, machine.CInstance.MachineId
		}
	} else {
		if machine.NewInfra {
			name, id = machine.Machine.Name, machine.Machine.Id
		} else {
			name, id = machine.Instance.Name, machine.Instance.MachineId
		}
	}
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "rm")
	command = append( command,  "-f")
	command = append( command,  name + "-" + id)
	cmd := executeSyncCommand(command)
	commandChannel <- cmd
	bytes, err := cmd.CombinedOutput()
	state := Machine_State_Running
	if err != nil {
		state = Machine_State_Stopped
	}
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: DestroyMachine,
		Error: err,
		Result: string(bytes),
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachineExecutor)  StopMachine(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd) {
	var name, id string
	if machine.IsCloud {
		if machine.NewInfra {
			name, id= machine.CMachine.Name, machine.CMachine.Id
		} else {
			name, id= machine.CInstance.Name, machine.CInstance.MachineId
		}
	} else {
		if machine.NewInfra {
			name, id = machine.Machine.Name, machine.Machine.Id
		} else {
			name, id = machine.Instance.Name, machine.Instance.MachineId
		}
	}
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "stop")
	command = append( command,  name + "-" + id)
	cmd := executeSyncCommand(command)
	commandChannel <- cmd
	bytes, err := cmd.CombinedOutput()
	state := Machine_State_Running
	if err != nil {
		state = Machine_State_Stopped
	}
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: StopMachine,
		Error: err,
		Result: string(bytes),
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachineExecutor)  StartMachine(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd) {
	var name, id string
	if machine.IsCloud {
		if machine.NewInfra {
			name, id= machine.CMachine.Name, machine.CMachine.Id
		} else {
			name, id= machine.CInstance.Name, machine.CInstance.MachineId
		}
	} else {
		if machine.NewInfra {
			name, id = machine.Machine.Name, machine.Machine.Id
		} else {
			name, id = machine.Instance.Name, machine.Instance.MachineId
		}
	}
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "start")
	command = append( command,  name + "-" + id)
	cmd := executeSyncCommand(command)
	commandChannel <- cmd
	bytes, err := cmd.CombinedOutput()
	state := Machine_State_Running
	if err != nil {
		state = Machine_State_Stopped
	}
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: StartMachine,
		Error: err,
		Result: string(bytes),
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachineExecutor)  RestartMachine(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd) {
	var name, id string
	if machine.IsCloud {
		if machine.NewInfra {
			name, id= machine.CMachine.Name, machine.CMachine.Id
		} else {
			name, id= machine.CInstance.Name, machine.CInstance.MachineId
		}
	} else {
		if machine.NewInfra {
			name, id = machine.Machine.Name, machine.Machine.Id
		} else {
			name, id = machine.Instance.Name, machine.Instance.MachineId
		}
	}
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "restart")
	command = append( command,  name + "-" + id)
	cmd := executeSyncCommand(command)
	commandChannel <- cmd
	bytes, err := cmd.CombinedOutput()
	state := Machine_State_Running
	if err != nil {
		state = Machine_State_Stopped
	}
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: RestartMachine,
		Error: err,
		Result: string(bytes),
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachineExecutor)  MachineStatus(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd) {
	var name, id string
	if machine.IsCloud {
		if machine.NewInfra {
			name, id= machine.CMachine.Name, machine.CMachine.Id
		} else {
			name, id= machine.CInstance.Name, machine.CInstance.MachineId
		}
	} else {
		if machine.NewInfra {
			name, id = machine.Machine.Name, machine.Machine.Id
		} else {
			name, id = machine.Instance.Name, machine.Instance.MachineId
		}
	}
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "status")
	command = append( command,  name + "-" + id)
	cmd := executeSyncCommand(command)
	commandChannel <- cmd
	bytesArray, err := cmd.CombinedOutput()
	state := Machine_State_None
	state = GetStateFromMachineAnswer(string(bytesArray))
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: StatusMachine,
		Error: err,
		Result: string(bytesArray),
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachineExecutor)  MachineEnv(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd) {
	var name, id string
	if machine.IsCloud {
		if machine.NewInfra {
			name, id= machine.CMachine.Name, machine.CMachine.Id
		} else {
			name, id= machine.CInstance.Name, machine.CInstance.MachineId
		}
	} else {
		if machine.NewInfra {
			name, id = machine.Machine.Name, machine.Machine.Id
		} else {
			name, id = machine.Instance.Name, machine.Instance.MachineId
		}
	}
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "env")
	command = append( command,  name + "-" + id)
	cmd := executeSyncCommand(command)
	commandChannel <- cmd
	bytes, err := cmd.CombinedOutput()
	state := Machine_State_Running
	if err != nil {
		state = Machine_State_Stopped
	}
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: MachineEnvironment,
		Error: err,
		Result: string(bytes),
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachineExecutor)  MachineInspect(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd) {
	var name, id string
	if machine.IsCloud {
		if machine.NewInfra {
			name, id= machine.CMachine.Name, machine.CMachine.Id
		} else {
			name, id= machine.CInstance.Name, machine.CInstance.MachineId
		}
	} else {
		if machine.NewInfra {
			name, id = machine.Machine.Name, machine.Machine.Id
		} else {
			name, id = machine.Instance.Name, machine.Instance.MachineId
		}
	}
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "inspect")
	command = append( command,  name + "-" + id)
	cmd := executeSyncCommand(command)
	commandChannel <- cmd
	bytes, err := cmd.CombinedOutput()
	state := Machine_State_Running
	if err != nil {
		state = Machine_State_Stopped
	}
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: MachineInspect,
		Error: err,
		Result: string(bytes),
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}

func (machine *DockerMachineExecutor)  MachineIPAddress(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd) {
	var name, id string
	if machine.IsCloud {
		if machine.NewInfra {
			name, id= machine.CMachine.Name, machine.CMachine.Id
		} else {
			name, id= machine.CInstance.Name, machine.CInstance.MachineId
		}
	} else {
		if machine.NewInfra {
			name, id = machine.Machine.Name, machine.Machine.Id
		} else {
			name, id = machine.Instance.Name, machine.Instance.MachineId
		}
	}
	var command []string = make([]string, 0)
	command = append( command,  "docker-machine")
	command = append( command,  "ip")
	command = append( command,  name + "-" + id)
	cmd := executeSyncCommand(command)
	commandChannel <- cmd
	bytes, err := cmd.CombinedOutput()
	state := Machine_State_Running
	if err != nil {
		state = Machine_State_Stopped
	}
	commandPipe <- MachineMessage{
		Complete: true,
		Cmd: command,
		Project: machine.Project,
		Infra: machine.Infra,
		Operation: MachineIPAddress,
		Error: err,
		Result: string(bytes),
		State: state,
		InstanceId: machine.InstanceId,
		IsCloud: machine.IsCloud,
	}
}
 func (machine *DockerMachineExecutor) IsThreadSafeCommand() bool {
	 return false
 }

func (machine *DockerMachineExecutor) SetControlStructure(Control *MachineControlStructure) {
	machine.Control = Control
}