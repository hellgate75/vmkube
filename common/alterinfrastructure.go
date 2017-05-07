package common

import (
	"errors"
	"fmt"
	"strings"
	"vmkube/model"
	"vmkube/procedures"
	"vmkube/tasks"
	"vmkube/utils"
	"vmkube/vmio"
)

func DescribeInstallation(installation *model.Installation, instanceName string, padding string) error {
	utils.PrintlnImportant(fmt.Sprintf("%sInstallation : %s", padding, installation.Id))
	utils.PrintlnImportant(fmt.Sprintf("%sInstance %s [Id: %s]", padding, instanceName, installation.InstanceId))
	utils.PrintlnImportant(fmt.Sprintf("%sOn Cloud : %s", padding, BoolToString(installation.IsCloud)))
	utils.PrintlnImportant(fmt.Sprintf("%sSuccess : %s", padding, BoolToString(installation.Success)))
	utils.PrintlnImportant(fmt.Sprintf("%sErrors : %s", padding, BoolToString(installation.Errors)))
	utils.PrintlnInfo(fmt.Sprintf("%sType : %s", padding, model.InstanceInstallationToString(installation.Type)))
	utils.PrintlnInfo(fmt.Sprintf("%sEnvironment : %s", padding, model.InstanceEnvironmentToString(installation.Environment)))
	utils.PrintlnInfo(fmt.Sprintf("%sRole : %s", padding, model.InstanceRoleToString(installation.Role)))
	fmt.Printf("%sLast Execution : %d-%02d-%02d %02d:%02d:%02d\n", padding,
		installation.LastExecution.Year(), installation.LastExecution.Month(), installation.LastExecution.Day(),
		installation.LastExecution.Hour(), installation.LastExecution.Minute(), installation.LastExecution.Second())
	utils.PrintlnInfo(fmt.Sprintf("%sLast Message : %s", padding, installation.LastMessage))

	utils.PrintlnImportant(fmt.Sprintf("%sInstallation Logs : ", padding))
	var logsInfo InfrastructureLogsInfo = InfrastructureLogsInfo{
		Format: "",
		Logs:   installation.Logs,
	}
	err := logsInfo.ReadLogFiles()
	if err == nil {
		for _, line := range logsInfo.Logs.LogLines {
			utils.PrintlnInfo(fmt.Sprintf("%s  %s", line, padding))
		}
		if len(logsInfo.Logs.LogLines) == 0 {
			utils.PrintlnImportant(fmt.Sprintf("%s  No logs available", padding))
		}
	} else {
		utils.PrintlnImportant(fmt.Sprintf("%s  Unable to read logs", padding))
	}
	return nil
}

func DescribeInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool, instanceState procedures.MachineState) error {
	if isCloud {
		utils.PrintlnImportant(fmt.Sprintf("Cloud Instance : %s [Id: %s]", cloudInstance.Name, cloudInstance.Id))
		utils.PrintlnImportant(fmt.Sprintf("Status : %s", instanceState.String()))
		utils.PrintlnImportant(fmt.Sprintf("Disabled : %s", BoolToString(cloudInstance.Disabled)))
		utils.PrintlnInfo(fmt.Sprintf("Driver : %s", cloudInstance.Driver))
		utils.PrintlnInfo(fmt.Sprintf("Hostname : %s", cloudInstance.Hostname))
		utils.PrintlnInfo(fmt.Sprintf("IPAdddress : %s", cloudInstance.IPAddress))
		utils.PrintlnInfo(fmt.Sprintf("Project Machine Reference : %s", cloudInstance.MachineId))
		utils.PrintlnInfo(fmt.Sprintf("Roles : %s", strings.Join(cloudInstance.Roles, ", ")))
		utils.PrintlnImportant("Options : ")
		for _, couple := range cloudInstance.Options {
			if len(couple) > 1 {
				utils.PrintlnInfo(fmt.Sprintf("  %s=%s", couple[0], couple[1]))
			} else {
				utils.PrintlnInfo(fmt.Sprintf("  %s", couple[0]))
			}
		}
		if len(cloudInstance.Options) == 0 {
			utils.PrintlnImportant("  Options not defined")
		}
		utils.PrintlnInfo(fmt.Sprintf("Instance Details : %s", cloudInstance.InspectJSON))
	} else {
		utils.PrintlnImportant(fmt.Sprintf("Local Instance : %s [Id: %s]", instance.Name, instance.Id))
		utils.PrintlnImportant(fmt.Sprintf("Status : %s", instanceState.String()))
		utils.PrintlnImportant(fmt.Sprintf("Disabled : %s", BoolToString(instance.Disabled)))
		utils.PrintlnInfo(fmt.Sprintf("Driver : %s", instance.Driver))
		utils.PrintlnInfo(fmt.Sprintf("Hostname : %s", instance.Hostname))
		utils.PrintlnInfo(fmt.Sprintf("IPAdddress : %s", instance.IPAddress))
		utils.PrintlnInfo(fmt.Sprintf("OS : %s v.%s", instance.OSType, instance.OSVersion))
		utils.PrintlnInfo(fmt.Sprintf("Project Machine Reference : %s", instance.MachineId))
		if instance.Cpus > 0 {
			utils.PrintlnInfo(fmt.Sprintf("Cpu : %d", instance.Cpus))
		} else {
			utils.PrintlnInfo("Cpu : default")
		}
		var counter int = 0
		for _, disk := range instance.Disks {
			counter++
			if instance.Memory > 0 {
				utils.PrintlnInfo(fmt.Sprintf("Disk #%d : %d GB", counter, disk.Size))
			} else {
				utils.PrintlnInfo(fmt.Sprintf("Disk #%d : default", counter))
			}
		}
		if instance.Memory > 0 {
			utils.PrintlnInfo(fmt.Sprintf("Memory : %dMB", instance.Memory))
		} else {
			utils.PrintlnInfo("Memory : default")
		}
		utils.PrintlnInfo(fmt.Sprintf("Prevent Share Home : %s", BoolToString(instance.NoShare)))
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
		utils.PrintlnInfo(fmt.Sprintf("  Enabled : %s", BoolToString(instance.Swarm.Enabled)))
		utils.PrintlnInfo(fmt.Sprintf("  Master Node : %s", BoolToString(instance.Swarm.IsMaster)))
		utils.PrintlnInfo(fmt.Sprintf("  Swarm Startegy : %s", instance.Swarm.Strategy))
		utils.PrintlnInfo(fmt.Sprintf("  Discovery Token : %s", instance.Swarm.DiscoveryToken))
		utils.PrintlnInfo(fmt.Sprintf("  Swarm Host : %s", instance.Swarm.Host))
		utils.PrintlnInfo(fmt.Sprintf("  Swarm Docker Image : %s", instance.Swarm.Image))
		utils.PrintlnInfo(fmt.Sprintf("  Join Options : %s", strings.Join(instance.Swarm.JoinOpts, ", ")))
		utils.PrintlnInfo(fmt.Sprintf("  TLS San Options : %s", strings.Join(instance.Swarm.TLSSan, ", ")))
		utils.PrintlnInfo(fmt.Sprintf("  Use Address : %s", BoolToString(instance.Swarm.UseAddress)))
		utils.PrintlnInfo(fmt.Sprintf("  Experimental : %s", BoolToString(instance.Swarm.UseAddress)))
		utils.PrintlnImportant("Options : ")
		for _, couple := range instance.Options {
			if len(couple) > 1 {
				utils.PrintlnInfo(fmt.Sprintf("  %s=%s", couple[0], couple[1]))
			} else {
				utils.PrintlnInfo(fmt.Sprintf("  %s", couple[0]))
			}
		}
		if len(instance.Options) == 0 {
			utils.PrintlnImportant("  Options not defined")
		}
		utils.PrintlnInfo(fmt.Sprintf("Instance Details : \n%s\n", instance.InspectJSON))
		utils.PrintlnImportant("Instance Logs : ")
		var logsInfo InfrastructureLogsInfo
		var instanceName string = ""
		if isCloud {
			logsInfo = InfrastructureLogsInfo{
				Format: "",
				Logs:   cloudInstance.Logs,
			}
			instanceName = cloudInstance.Name
		} else {
			logsInfo = InfrastructureLogsInfo{
				Format: "",
				Logs:   instance.Logs,
			}
			instanceName = instance.Name
		}
		err := logsInfo.ReadLogFiles()
		if err == nil {
			for _, line := range logsInfo.Logs.LogLines {
				utils.PrintlnInfo(fmt.Sprintf("  %s", line))
			}
			if len(logsInfo.Logs.LogLines) == 0 {
				utils.PrintlnImportant("  No logs available")
			}
		} else {
			utils.PrintlnImportant("  Unable to read logs")
		}
		utils.PrintlnImportant("Installations : ")
		var installations []*model.Installation = ExtractInstallations(&infrastructure, instance, cloudInstance, isCloud)
		for _, installation := range installations {
			DescribeInstallation(installation, instanceName, "  ")
		}
		if len(installations) == 0 {
			utils.PrintlnImportant("No installation for instance")
		}
	}
	return nil
}

func StartInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool, instanceState procedures.MachineState) error {
	var exclusionIdList []string = make([]string, 0)
	if isCloud {
		if cloudInstance.Disabled {
			if cloudInstance.Disabled {
				return errors.New(fmt.Sprintf("Cloud Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s] is disabled!!", cloudInstance.Name, cloudInstance.Id, infrastructure.Name, infrastructure.Id))
			}
			if instanceState == procedures.Machine_State_Running {
				return errors.New(fmt.Sprintf("Cloud Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s] already started!!", cloudInstance.Name, cloudInstance.Id, infrastructure.Name, infrastructure.Id))
			}
		}
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{cloudInstance.Id})
	} else {
		if instance.Disabled {
			return errors.New(fmt.Sprintf("Local Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s] is disabled!!", instance.Name, instance.Id, infrastructure.Name, infrastructure.Id))
		}
		if instanceState == procedures.Machine_State_Running {
			return errors.New(fmt.Sprintf("Local Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s] already started!!", instance.Name, instance.Id, infrastructure.Name, infrastructure.Id))
		}
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{instance.Id})
	}
	startCouples, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.StartMachine, exclusionIdList)
	if err != nil {
		return err
	}
	if len(startCouples) == 0 {
		return errors.New("No Instance available for start procedure...")
	}
	NumThreads := 1
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))
	errorsList := ExecuteInfrastructureActions(infrastructure, startCouples, NumThreads, func(task tasks.SchedulerTask) {})
	if len(errorsList) > 0 {
		return errorsList[0]
	}
	return nil
}

func RestartInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool, instanceState procedures.MachineState) error {
	var exclusionIdList []string = make([]string, 0)
	if isCloud {
		if cloudInstance.Disabled {
			return errors.New(fmt.Sprintf("Cloud Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s] is disabled!!", cloudInstance.Name, cloudInstance.Id, infrastructure.Name, infrastructure.Id))
		}

		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{cloudInstance.Id})
	} else {
		if instance.Disabled {
			return errors.New(fmt.Sprintf("Local Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s] is disabled!!", instance.Name, instance.Id, infrastructure.Name, infrastructure.Id))
		}
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{instance.Id})
	}
	restartCouples, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.RestartMachine, exclusionIdList)
	if err != nil {
		return err
	}
	if len(restartCouples) == 0 {
		return errors.New("No Instance available for restart procedure...")
	}
	NumThreads := 1
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))
	errorsList := ExecuteInfrastructureActions(infrastructure, restartCouples, NumThreads, func(task tasks.SchedulerTask) {})
	if len(errorsList) > 0 {
		return errorsList[0]
	}
	return nil
}

func StopInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool, instanceState procedures.MachineState) error {
	var exclusionIdList []string = make([]string, 0)
	if isCloud {
		if cloudInstance.Disabled {
			return errors.New(fmt.Sprintf("Cloud Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s] is disabled!!", cloudInstance.Name, cloudInstance.Id, infrastructure.Name, infrastructure.Id))
		}
		if instanceState == procedures.Machine_State_Stopped {
			return errors.New(fmt.Sprintf("Cloud Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s] already stopped!!", cloudInstance.Name, cloudInstance.Id, infrastructure.Name, infrastructure.Id))
		}
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{cloudInstance.Id})
	} else {
		if instance.Disabled {
			return errors.New(fmt.Sprintf("Local Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s] is disabled!!", instance.Name, instance.Id, infrastructure.Name, infrastructure.Id))
		}
		if instanceState == procedures.Machine_State_Stopped {
			return errors.New(fmt.Sprintf("Local Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s] already stopped!!", instance.Name, instance.Id, infrastructure.Name, infrastructure.Id))
		}
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{instance.Id})
	}
	stopCouples, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.StopMachine, exclusionIdList)
	if err != nil {
		return err
	}
	if len(stopCouples) == 0 {
		return errors.New("No Instance available for stop procedure...")
	}
	NumThreads := 1
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))
	errorsList := ExecuteInfrastructureActions(infrastructure, stopCouples, NumThreads, func(task tasks.SchedulerTask) {})
	if len(errorsList) > 0 {
		return errorsList[0]
	}
	return nil
}

func DisableInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool, instanceState procedures.MachineState, descriptor model.ProjectsDescriptor) error {
	var err error

	if isCloud {
		if cloudInstance.Disabled {
			return errors.New(fmt.Sprintf("Cloud Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s] already disabled!!", cloudInstance.Name, cloudInstance.Id, infrastructure.Name, infrastructure.Id))
		}
	} else {
		if instance.Disabled {
			return errors.New(fmt.Sprintf("Local Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s] already disabled!!", instance.Name, instance.Id, infrastructure.Name, infrastructure.Id))
		}
	}

	if instanceState == procedures.Machine_State_Running {
		var exclusionIdList []string = make([]string, 0)
		var actionCouples []tasks.ActivityCouple = make([]tasks.ActivityCouple, 0)
		var machineId, instanceId string
		if isCloud {
			machineId = cloudInstance.MachineId
			instanceId = cloudInstance.Id
			exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{cloudInstance.Id})
		} else {
			machineId = instance.MachineId
			instanceId = instance.Id
			exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{instance.Id})
		}

		utils.PrintlnWarning(fmt.Sprintf("Stopping Instance Id: '%s', Machine Id : '%s' from Project : '%s' [Id: %s]...", instanceId, machineId, descriptor.Name, descriptor.Id))

		stopCouples, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.StopMachine, exclusionIdList)
		if err != nil {
			return err
		}
		actionCouples = append(actionCouples, stopCouples...)
		if len(actionCouples) == 0 {
			return errors.New("No Instance available for destroy procedure...")
		}
		NumThreads := 1
		utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))
		errorsList := ExecuteInfrastructureActions(infrastructure, actionCouples, NumThreads, func(task tasks.SchedulerTask) {})

		if len(errorsList) > 0 {
			return errorsList[0]
		}
	}

	if isCloud {
		utils.PrintlnWarning(fmt.Sprintf("Disabling Cloud Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s]", cloudInstance.Name, cloudInstance.Id, infrastructure.Name, infrastructure.Id))
		for i := 0; i < len(infrastructure.Domains); i++ {
			for j := 0; j < len(infrastructure.Domains[i].Networks); j++ {
				for k := 0; k < len(infrastructure.Domains[i].Networks[j].CloudInstances); k++ {
					if infrastructure.Domains[i].Networks[j].CloudInstances[k].Id == cloudInstance.Id {
						infrastructure.Domains[i].Networks[j].CloudInstances[k].Disabled = true
					}
				}
			}
		}
	} else {
		utils.PrintlnWarning(fmt.Sprintf("Disabling Local Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s]", instance.Name, instance.Id, infrastructure.Name, infrastructure.Id))
		for i := 0; i < len(infrastructure.Domains); i++ {
			for j := 0; j < len(infrastructure.Domains[i].Networks); j++ {
				for k := 0; k < len(infrastructure.Domains[i].Networks[j].LocalInstances); k++ {
					if infrastructure.Domains[i].Networks[j].LocalInstances[k].Id == instance.Id {
						infrastructure.Domains[i].Networks[j].LocalInstances[k].Disabled = true
					}
				}
			}
		}
	}

	utils.PrintlnImportant("No changes allowed for this instance ...")

	utils.PrintlnWarning(fmt.Sprintf("Saving information to Infrastructure : %s [Id: %s]", infrastructure.Name, infrastructure.Id))

	err = vmio.SaveInfrastructure(infrastructure)

	if err != nil {
		return err
	}
	return nil
}

func EnableInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	var err error

	if isCloud {
		if !cloudInstance.Disabled {
			return errors.New(fmt.Sprintf("Cloud Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s] already enabled!!", cloudInstance.Name, cloudInstance.Id, infrastructure.Name, infrastructure.Id))
		}
		utils.PrintlnWarning(fmt.Sprintf("Enabling Cloud Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s]", cloudInstance.Name, cloudInstance.Id, infrastructure.Name, infrastructure.Id))
		for i := 0; i < len(infrastructure.Domains); i++ {
			for j := 0; j < len(infrastructure.Domains[i].Networks); j++ {
				for k := 0; k < len(infrastructure.Domains[i].Networks[j].CloudInstances); k++ {
					if infrastructure.Domains[i].Networks[j].CloudInstances[k].Id == cloudInstance.Id {
						infrastructure.Domains[i].Networks[j].CloudInstances[k].Disabled = false
					}
				}
			}
		}
	} else {
		if !instance.Disabled {
			return errors.New(fmt.Sprintf("Local Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s] already enabled!!", instance.Name, instance.Id, infrastructure.Name, infrastructure.Id))
		}
		utils.PrintlnWarning(fmt.Sprintf("Enabling Local Instance '%s' [Id: %s] part of Infrastructure : %s [Id: %s]", instance.Name, instance.Id, infrastructure.Name, infrastructure.Id))
		for i := 0; i < len(infrastructure.Domains); i++ {
			for j := 0; j < len(infrastructure.Domains[i].Networks); j++ {
				for k := 0; k < len(infrastructure.Domains[i].Networks[j].LocalInstances); k++ {
					if infrastructure.Domains[i].Networks[j].LocalInstances[k].Id == instance.Id {
						infrastructure.Domains[i].Networks[j].LocalInstances[k].Disabled = false
					}
				}
			}
		}
	}

	utils.PrintlnImportant("Now changes available for this instance ...")

	utils.PrintlnWarning(fmt.Sprintf("Saving information to Infrastructure : %s [Id: %s]", infrastructure.Name, infrastructure.Id))

	err = vmio.SaveInfrastructure(infrastructure)
	if err != nil {
		return err
	}
	return nil
}

func RecreateInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool, instanceState procedures.MachineState, descriptor model.ProjectsDescriptor) error {
	utils.PrintlnWarning(fmt.Sprintf("Load Project : '%s' for inspection ...", descriptor.Name))

	project, err := vmio.LoadProject(descriptor.Id)

	if err != nil {
		return err
	}

	var exclusionIdList []string = make([]string, 0)
	var actionCouples []tasks.ActivityCouple = make([]tasks.ActivityCouple, 0)
	var machineId, instanceId string
	if isCloud {
		machineId = cloudInstance.MachineId
		instanceId = cloudInstance.Id
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{cloudInstance.Id})
	} else {
		machineId = instance.MachineId
		instanceId = instance.Id
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{instance.Id})
	}

	if instanceState == procedures.Machine_State_Running {
		stopCouples, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.StopMachine, exclusionIdList)
		if err != nil {
			return err
		}
		actionCouples = append(actionCouples, stopCouples...)
	}

	destroyCouples, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.DestroyMachine, exclusionIdList)
	if err != nil {
		return err
	}
	actionCouples = append(actionCouples, destroyCouples...)

	createCouples, err := tasks.GetTaskActivitiesExclusion(project, infrastructure, tasks.CreateMachine, exclusionIdList)

	if err != nil {
		return err
	}
	actionCouples = append(actionCouples, createCouples...)

	inspectCouples, err := tasks.GetTaskActivitiesExclusion(project, infrastructure, tasks.MachineInspect, exclusionIdList)

	if err != nil {
		return err
	}
	actionCouples = append(actionCouples, inspectCouples...)

	ipAddressCouples, err := tasks.GetTaskActivitiesExclusion(project, infrastructure, tasks.MachineIPAddress, exclusionIdList)

	if err != nil {
		return err
	}
	actionCouples = append(actionCouples, ipAddressCouples...)

	stopCouples, err := tasks.GetTaskActivitiesExclusion(project, infrastructure, tasks.StopMachine, exclusionIdList)

	if err != nil {
		return err
	}
	actionCouples = append(actionCouples, stopCouples...)

	extendsDiskCouples, err := tasks.GetTaskActivitiesExclusion(project, infrastructure, tasks.MachineExtendsDisk, exclusionIdList)

	if err != nil {
		return err
	}
	actionCouples = append(actionCouples, extendsDiskCouples...)

	utils.PrintlnWarning(fmt.Sprintf("Recreating Machine Id: '%s' (Instance ID: %s) in Project : '%s' [Id: %s]  and Infrastructure : '%s' [Id: %s]...", machineId, instanceId, project.Name, project.Id, infrastructure.Name, infrastructure.Id))

	if len(actionCouples) == 0 {
		return errors.New("No Instance available for re-create procedure...")
	}
	NumThreads := 1
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))
	errorsList := ExecuteInfrastructureActions(infrastructure, actionCouples, NumThreads, func(task tasks.SchedulerTask) {})

	if len(errorsList) > 0 {
		return errorsList[0]
	}

	return nil
}

func DestroyInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool, instanceState procedures.MachineState, descriptor model.ProjectsDescriptor) error {
	utils.PrintlnWarning(fmt.Sprintf("Load Project : '%s' for inspection ...", descriptor.Name))

	project, err := vmio.LoadProject(descriptor.Id)

	if err != nil {
		return err
	}

	var exclusionIdList []string = make([]string, 0)
	var actionCouples []tasks.ActivityCouple = make([]tasks.ActivityCouple, 0)
	var machineId, instanceId string
	var instanceLogs model.LogStorage
	if isCloud {
		machineId = cloudInstance.MachineId
		instanceId = cloudInstance.Id
		instanceLogs = cloudInstance.Logs
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{cloudInstance.Id})
	} else {
		machineId = instance.MachineId
		instanceId = instance.Id
		instanceLogs = instance.Logs
		exclusionIdList = tasks.GetExclusionListExceptInstanceList(infrastructure, []string{instance.Id})
	}

	if instanceState == procedures.Machine_State_Running {
		stopCouples, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.StopMachine, exclusionIdList)
		if err != nil {
			return err
		}
		actionCouples = append(actionCouples, stopCouples...)
	}

	destroyCouples, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.DestroyMachine, exclusionIdList)
	if err != nil {
		return err
	}
	actionCouples = append(actionCouples, destroyCouples...)
	utils.PrintlnWarning(fmt.Sprintf("Removing Machine Id: '%s' from Project : '%s' [Id: %s]...", machineId, project.Name, project.Id))
	if len(actionCouples) == 0 {
		return errors.New("No Instance available for destroy procedure...")
	}
	err = RemoveProjectMachineById(&project, machineId)
	if err != nil {
		return err
	}
	utils.PrintlnWarning(fmt.Sprintf("Removing Instance Id: %s from Infrastructure : '%s' [Id: %s]...", instanceId, infrastructure.Name, infrastructure.Id))
	err = RemoveInfrastructureInstanceById(&infrastructure, instanceId)
	if err != nil {
		return err
	}

	NumThreads := 1
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))
	errorsList := ExecuteInfrastructureActions(infrastructure, actionCouples, NumThreads, func(task tasks.SchedulerTask) {})

	if len(errorsList) > 0 {
		return errorsList[0]
	}

	err = DeleteInfrastructureLogs(instanceLogs)
	if err != nil {
		return err
	}

	utils.PrintlnWarning(fmt.Sprintf("Overwriting Project : %s [Id: %s]", project.Name, project.Id))

	err = vmio.SaveProject(project)
	if err != nil {
		return err
	}

	utils.PrintlnWarning(fmt.Sprintf("Overwriting Infrastructure : %s [Id: %s]", infrastructure.Name, infrastructure.Id))

	err = vmio.SaveInfrastructure(infrastructure)
	if err != nil {
		return err
	}
	return nil
}

func AutoFixInfrastructureInstances(infrastructure model.Infrastructure) error {
	//TODO: Implement Autofix Procedure
	return errors.New("Alter Infrastructure, Autofix instancs not yet implemented ...")
}
