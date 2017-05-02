package operations

import (
	"errors"
	"vmkube/model"
	"vmkube/utils"
	"fmt"
	"strings"
	"vmkube/procedures"
)

func DescribeInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool, instanceState procedures.MachineState) error {
	//TODO: Implement Instance Status
	if isCloud {
		utils.PrintlnImportant(fmt.Sprintf("Cloud Instance : %s [Id: %s]", cloudInstance.Name, cloudInstance.Id))
		utils.PrintlnImportant(fmt.Sprintf("Status : %s", instanceState.String()))
		utils.PrintlnInfo(fmt.Sprintf("Driver : %s", cloudInstance.Driver))
		utils.PrintlnInfo(fmt.Sprintf("Hostname : %s", cloudInstance.Hostname))
		utils.PrintlnInfo(fmt.Sprintf("IPAdddress : %s", cloudInstance.IPAddress))
		utils.PrintlnInfo(fmt.Sprintf("Project Machine Reference : %s", cloudInstance.MachineId))
		utils.PrintlnInfo(fmt.Sprintf("Roles : %s", strings.Join(cloudInstance.Roles, ", ")))
		utils.PrintlnImportant("Options : ")
		for _,couple := range cloudInstance.Options {
			if len(couple) > 1 {
				utils.PrintlnInfo(fmt.Sprintf("  %s=%s", couple[0], couple[1]))
			} else  {
				utils.PrintlnInfo(fmt.Sprintf("  %s", couple[0]))
			}
		}
		if len(cloudInstance.Options) == 0 {
			utils.PrintlnImportant("  Options not defined")
		}
		utils.PrintlnInfo(fmt.Sprintf("Instance Details : %s", cloudInstance.InspectJSON))
	} else  {
		utils.PrintlnImportant(fmt.Sprintf("Local Instance : %s [Id: %s]", instance.Name, instance.Id))
		utils.PrintlnImportant(fmt.Sprintf("Status : %s", instanceState.String()))
		utils.PrintlnInfo(fmt.Sprintf("Driver : %s", instance.Driver))
		utils.PrintlnInfo(fmt.Sprintf("Hostname : %s", instance.Hostname))
		utils.PrintlnInfo(fmt.Sprintf("IPAdddress : %s", instance.IPAddress))
		utils.PrintlnInfo(fmt.Sprintf("OS : %s v.%s", instance.OSType,instance.OSVersion))
		utils.PrintlnInfo(fmt.Sprintf("Project Machine Reference : %s", instance.MachineId))
		if instance.Cpus > 0 {
			utils.PrintlnInfo(fmt.Sprintf("Cpu : %d", instance.Cpus))
		} else  {
			utils.PrintlnInfo("Cpu : default")
		}
		var  counter int  = 0
		for _,disk := range instance.Disks {
			counter++
			if instance.Memory > 0 {
				utils.PrintlnInfo(fmt.Sprintf("Disk #%d : %d GB", counter, disk.Size))
			} else  {
				utils.PrintlnInfo(fmt.Sprintf("Disk #%d : default", counter))
			}
		}
		if instance.Memory > 0 {
			utils.PrintlnInfo(fmt.Sprintf("Memory : %dMB", instance.Memory))
		} else  {
			utils.PrintlnInfo("Memory : default")
		}
		utils.PrintlnInfo(fmt.Sprintf("Prevent Share Home : %t", instance.NoShare))
		utils.PrintlnInfo(fmt.Sprintf("Roles : %s", strings.Join(instance.Roles, ", ")))
		utils.PrintlnImportant("Docker Engine Options : ")
		utils.PrintlnInfo(fmt.Sprintf("  Storage Driver : %s", instance.Engine.StorageDriver))
		utils.PrintlnInfo(fmt.Sprintf("  Install URL : %s", instance.Engine.InstallURL))
		utils.PrintlnInfo(fmt.Sprintf("  Insecure Registry Options : %s", strings.Join(instance.Engine.InsecureRegistry, ", ")))
		utils.PrintlnInfo(fmt.Sprintf("  Registry Mirror Options : %s", strings.Join(instance.Engine.RegistryMirror, ", ")))
		utils.PrintlnInfo(fmt.Sprintf("  Environment Registry Options : %s", strings.Join(instance.Engine.Environment, ", ")))
		utils.PrintlnInfo(fmt.Sprintf("  Engine Options : %s", strings.Join(instance.Engine.Options, ", ")))
		utils.PrintlnInfo(fmt.Sprintf("  Label : %s", strings.Join(instance.Engine.Labels, ", ")))
		utils.PrintlnImportant("Swarm Options : ")
		utils.PrintlnInfo(fmt.Sprintf("  Enabled : %t", instance.Swarm.Enabled))
		utils.PrintlnInfo(fmt.Sprintf("  Master Node : %t", instance.Swarm.IsMaster))
		utils.PrintlnInfo(fmt.Sprintf("  Swarm Startegy : %s", instance.Swarm.Strategy))
		utils.PrintlnInfo(fmt.Sprintf("  Discovery Token : %s", instance.Swarm.DiscoveryToken))
		utils.PrintlnInfo(fmt.Sprintf("  Swarm Host : %s", instance.Swarm.Host))
		utils.PrintlnInfo(fmt.Sprintf("  Swarm Docker Image : %s", instance.Swarm.Image))
		utils.PrintlnInfo(fmt.Sprintf("  Join Options : %s", strings.Join(instance.Swarm.JoinOpts, ", ")))
		utils.PrintlnInfo(fmt.Sprintf("  TLS San Options : %s", strings.Join(instance.Swarm.TLSSan, ", ")))
		utils.PrintlnInfo(fmt.Sprintf("  Use Address : %t", instance.Swarm.UseAddress))
		utils.PrintlnInfo(fmt.Sprintf("  Experimental : %t", instance.Swarm.UseAddress))
		utils.PrintlnImportant("Options : ")
		for _,couple := range instance.Options {
			if len(couple) > 1 {
				utils.PrintlnInfo(fmt.Sprintf("  %s=%s", couple[0], couple[1]))
			} else  {
				utils.PrintlnInfo(fmt.Sprintf("  %s", couple[0]))
			}
		}
		if len(instance.Options) == 0 {
			utils.PrintlnImportant("  Options not defined")
		}
		utils.PrintlnInfo(fmt.Sprintf("Instance Details : %s", instance.InspectJSON))
	}
	return nil
}

func StartInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement Start Instance
	return errors.New("Command Start Instance not implemented!!")
}

func RestartInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement ReStart Instance
	return errors.New("Command Restart Instance not implemented!!")
}

func StopInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement Stop Instance
	return errors.New("Command Stop Instance not implemented!!")
}

func DisableInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement Disable Instance ??
	return errors.New("Command Disable Instance not implemented!!")
}

func EnableInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement Enable Instance ??
	return errors.New("Command Enable Instance not implemented!!")
}

func RecreateInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement Re-Create Instance
	return errors.New("Command Re-Create Instance not implemented!!")
}

func DestroyInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement Destroy Instance
	return errors.New("Command Destroy Instance not implemented!!")
}

func AutoFixInfrastructureInstances(infrastructure model.Infrastructure) error {
	return nil
}
