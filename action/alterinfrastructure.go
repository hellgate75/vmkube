package action

import (
	"errors"
	"vmkube/model"
	"vmkube/utils"
	"fmt"
	"strings"
	"vmkube/procedures"
	"vmkube/tasks"
	"vmkube/vmio"
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

func StartInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool, instanceState procedures.MachineState) error {
	var  exclusionIdList []string = make([]string, 0)
	if isCloud {
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{cloudInstance.Id})
	} else  {
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{instance.Id})
	}
	startCouples, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.StartMachine, exclusionIdList)
	if err != nil {
		return err
	}
	if len(startCouples) == 0 {
		return  errors.New("No Instance available for start procedure...")
	}
	NumThreads := 1
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))
	errorsList := ExecuteInfrastructureActions(infrastructure, startCouples, NumThreads, func(task tasks.ScheduleTask) {})
	if len(errorsList) > 0 {
		return  errorsList[0]
	}
	return nil
}

func RestartInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool, instanceState procedures.MachineState) error {
	var  exclusionIdList []string = make([]string, 0)
	if isCloud {
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{cloudInstance.Id})
	} else  {
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{instance.Id})
	}
	restartCouples, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.RestartMachine, exclusionIdList)
	if err != nil {
		return err
	}
	if len(restartCouples) == 0 {
		return  errors.New("No Instance available for restart procedure...")
	}
	NumThreads := 1
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))
	errorsList := ExecuteInfrastructureActions(infrastructure, restartCouples, NumThreads, func(task tasks.ScheduleTask) {})
	if len(errorsList) > 0 {
		return  errorsList[0]
	}
	return nil
}

func StopInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool, instanceState procedures.MachineState) error {
	var  exclusionIdList []string = make([]string, 0)
	if isCloud {
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{cloudInstance.Id})
	} else  {
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{instance.Id})
	}
	stopCouples, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.StopMachine, exclusionIdList)
	if err != nil {
		return err
	}
	if len(stopCouples) == 0 {
		return  errors.New("No Instance available for stop procedure...")
	}
	NumThreads := 1
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))
	errorsList := ExecuteInfrastructureActions(infrastructure, stopCouples, NumThreads, func(task tasks.ScheduleTask) {})
	if len(errorsList) > 0 {
		return  errorsList[0]
	}
	return nil
}

func DisableInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement Disable Instance ??
	return errors.New("Command Disable Instance not implemented!!")
}

func EnableInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement Enable Instance ??
	return errors.New("Command Enable Instance not implemented!!")
}

func RecreateInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool, descriptor model.ProjectsDescriptor) error {
	//TODO: Implement Re-Create Instance
	return errors.New("Command Re-Create Instance not implemented!!")
}

func DestroyInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool, descriptor model.ProjectsDescriptor) error {
	//TODO: Test Destroy Instance

	utils.PrintlnWarning(fmt.Sprintf("Load Project : '%s' for inspection ...", descriptor.Name))

	project, err := vmio.LoadProject(descriptor.Id)

	if err != nil {
		return  err
	}

	var  exclusionIdList []string = make([]string, 0)
	var machineId, instanceId string
	if isCloud {
		machineId = cloudInstance.MachineId
		instanceId = cloudInstance.Id
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{cloudInstance.Id})
	} else  {
		machineId = instance.MachineId
		instanceId = instance.Id
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{instance.Id})
	}
	stopCouples, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.DestroyMachine, exclusionIdList)
	if err != nil {
		return err
	}
	utils.PrintlnWarning(fmt.Sprintf("Removing Machine Id: '%s' from Project : '%s' [Id: %s]...", machineId, project.Name, project.Id))
	if len(stopCouples) == 0 {
		return  errors.New("No Instance available for destroy procedure...")
	}
	err = RemoveProjectMachineById(&project, machineId)
	if err != nil {
		return  err
	}
	utils.PrintlnWarning(fmt.Sprintf("Removing Instance Id: %s from Infrastructure : '%s' [Id: %s]...", instanceId, infrastructure.Name, infrastructure.Id))
	err = RemoveInfrastructureInstanceById(&infrastructure, instanceId)
	if err != nil {
		return  err
	}

	NumThreads := 1
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))
	errorsList := ExecuteInfrastructureActions(infrastructure, stopCouples, NumThreads, func(task tasks.ScheduleTask) {})

	if len(errorsList) > 0 {
		return  errorsList[0]
	}

	utils.PrintlnWarning(fmt.Sprintf("Overwriting Project : %s [Id: %s]", project.Name, project.Id))

	err = vmio.SaveProject(project)
	if err != nil {
		return  err
	}

	utils.PrintlnWarning(fmt.Sprintf("Overwriting Infrastructure : %s [Id: %s]", infrastructure.Name, infrastructure.Id))

	err = vmio.SaveInfrastructure(infrastructure)
	if err != nil {
		return  err
	}
	return nil
}

func AutoFixInfrastructureInstances(infrastructure model.Infrastructure) error {
	return nil
}
