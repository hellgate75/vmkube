package action

import (
	"fmt"
	"errors"
	"vmkube/utils"
	"vmkube/vmio"
	"vmkube/model"
)

type InfrastructureActions interface {
	CheckInfra() bool
	CreateInfra() (Response, error)
	AlterInfra() (Response, error)
	DeleteInfra() (Response, error)
	StartInfra() (Response, error)
	StopInfra() (Response, error)
	RestartInfra() (Response, error)
	ListInfras() (Response, error)
	StatusInfra() (Response, error)
}

func (request *CmdRequest) CheckInfra() bool {
	if len(request.Arguments.Helper.Options) > 0 {
		correctness := true
		for _,option := range request.Arguments.Helper.Options {
			if option.Mandatory {
				//Mandatory Option
				found := false
				for _,argument := range request.Arguments.Options {
					if CorrectInput(argument[0]) == option.Option {
						found = true
						break
					}
				}
				if !found {
					correctness = false
					utils.PrintlnError(fmt.Sprintf("Option '--%s' is mandatory!!\n", option.Option))
				}
			}
		}
		if !correctness {
			return  false
		}
	}
	return  true
}

func (request *CmdRequest) CreateInfra() (Response, error) {
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) AlterInfra() (Response, error) {
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) DeleteInfra() (Response, error) {
	Name := ""
	Force := false
	for _,option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "force" == CorrectInput(option[0]) {
			Force = GetBoolean(option[1])
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "Infrastrcuture Name not provided",
			Status: false,},errors.New("Unable to execute task")
	}
	if ! Force {
		AllowInfraDeletion := utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with deletion process for Infrastructure named '%s'?", Name))
		if ! AllowInfraDeletion {
			response := Response{
				Status: false,
				Message: "User task interruption",
			}
			return response, errors.New("Unable to execute task")
		}
	}
	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	iFaceInfra := vmio.IFaceInfra{
		Id: descriptor.InfraId,
		ProjectId: descriptor.Id,
	}
	iFaceInfra.WaitForUnlock()

	infrastructure, err := vmio.LoadInfrastructure(descriptor.Id)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	for _,domain := range infrastructure.Domains {
		for _,network := range domain.Networks {
			for _,instance := range network.Instances {
				err = DeleteInfrastructureLogs(instance.Logs)
				if err != nil {
					response := Response{
						Status: false,
						Message: err.Error(),
					}
					return  response, errors.New("Unable to execute task")
				}
			}
			for _,instance := range network.CInstances {
				err = DeleteInfrastructureLogs(instance.Logs)
				if err != nil {
					response := Response{
						Status: false,
						Message: err.Error(),
					}
					return  response, errors.New("Unable to execute task")
				}
			}
			for _,installation := range network.Installations {
				err = DeleteInfrastructureLogs(installation.Logs)
				if err != nil {
					response := Response{
						Status: false,
						Message: err.Error(),
					}
					return  response, errors.New("Unable to execute task")
				}
			}
		}
	}
	
	iFaceRollbackIndex := IFaceRollBackIndex{
		Id: descriptor.InfraId,
		ProjectId: descriptor.Id,
	}
	iFaceRollbackIndex.WaitForUnlock()
	
	indexInfo := ProjectRollbackIndexInfo{
		Format: "",
		Index: RollBackIndex{
			Id: descriptor.InfraId,
			ProjectId: descriptor.Id,
			IndexList: []RollBackSegmentIndex{},
		},
	}
	LockRollBackIndex(indexInfo.Index)
	err = indexInfo.Read()
	UnlockRollBackIndex(indexInfo.Index)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	iFaceRollbackIndex.WaitForUnlock()
	LockRollBackIndex(indexInfo.Index)
	err = indexInfo.Delete()
	UnlockRollBackIndex(indexInfo.Index)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	vmio.LockInfrastructureById(descriptor.Id, descriptor.InfraId)
	var info vmio.InfrastructureInfo =  vmio.InfrastructureInfo{
		Format: "",
		Infra: model.Infrastructure{
			ProjectId: descriptor.Id,
			Id: descriptor.InfraId,
		},
	}
	err = info.Delete()
	vmio.UnlockInfrastructureById(descriptor.Id, descriptor.InfraId)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
}

func (request *CmdRequest) BackupInfra() (Response, error) {
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) RecoverInfra() (Response, error) {
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) StartInfra() (Response, error) {
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) StopInfra() (Response, error) {
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) RestartInfra() (Response, error) {
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) ListInfras() (Response, error) {
	indexes, error := vmio.LoadIndex()
	if error != nil {
		response := Response{
			Status: false,
			Message: error.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	if len(indexes.Projects) > 0 {
		utils.PrintlnImportant(fmt.Sprintf("%s  %s  %s", utils.StrPad("Infrastructure Id", 40), utils.StrPad("Infrastructure Name", 40), utils.StrPad("Active", 6)))
	} else {
		utils.PrintlnImportant("No Infrastructures found")
	}
	for _,index := range indexes.Projects {
		active := "no"
		if index.Active {
			active = "yes"
		}
		fmt.Printf("%s  %s  %s\n", utils.StrPad(index.InfraId, 40), utils.StrPad(index.InfraName, 40), utils.StrPad("  "+active, 6))
		
	}
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
}

func (request *CmdRequest) StatusInfra() (Response, error) {
	Name := ""
	for _,option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = option[1]
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "Infrastrcuture Name not provided",
			Status: false,},errors.New("Unable to execute task")
	}
	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	iFaceInfra := vmio.IFaceInfra{
		Id: descriptor.InfraId,
		ProjectId: descriptor.Id,
	}
	iFaceInfra.WaitForUnlock()
	
	infrastructure, err := vmio.LoadInfrastructure(descriptor.Id)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	modified := "no"
	if infrastructure.Altered {
		modified = "yes"
	}
	errors := "no"
	if infrastructure.Errors {
		errors = "yes"
	}
	fmt.Printf("Id: %s\nInfrastructure: %s\nModified: %s\n", infrastructure.Id,infrastructure.Name, modified)
	created := "no"
	if infrastructure.Created {
		created = "yes"
	}
	fmt.Printf("Created [%s] : %d-%02d-%02d %02d:%02d:%02d\n", created,
		infrastructure.Creation.Year(), infrastructure.Creation.Month(), infrastructure.Creation.Day(),
		infrastructure.Creation.Hour(), infrastructure.Creation.Minute(), infrastructure.Creation.Second())
	fmt.Printf("Modified : %d-%02d-%02d %02d:%02d:%02d\n",
		infrastructure.Modified.Year(), infrastructure.Modified.Month(), infrastructure.Modified.Day(),
		infrastructure.Modified.Hour(), infrastructure.Modified.Minute(), infrastructure.Modified.Second())
	fmt.Printf("Errors: %s\nLast Message: %s\n", errors,infrastructure.LastMessage)
	fmt.Printf("Domains: %d\n", len(infrastructure.Domains))
	for _,domain := range infrastructure.Domains {
		num, options := vmio.StripOptions(domain.Options)
		fmt.Printf("Domain: %s (Id: %s) - Options [%d] :%s\n", domain.Name, domain.Id, num, options)
		fmt.Printf("Networks: %d\n", len(domain.Networks))
		for _,network := range domain.Networks {
			num, options := vmio.StripOptions(network.Options)
			fmt.Printf("   Network: %s (Id: %s) - Options [%d] :%s\n", network.Name, network.Id, num, options)
			fmt.Printf("   Instances: %d\n", len(network.Instances))
			serversMap := make(map[string]string)
			for _,server := range network.Instances {
				serversMap[server.Id] = server.Name
				fmt.Printf("      Instance: %s (Id: %s) - Driver: %s - OS : %s:%s - IP Address: %s\n", server.Name, server.Id, server.Driver, server.OSType, server.OSVersion, server.IPAddress)
			}
			fmt.Printf("   Cloud Instances: %d\n", len(network.CInstances))
			for _,server := range network.CInstances {
				serversMap[server.Id] = server.Name
				num, options := vmio.StripOptions(server.Options)
				fmt.Printf("      Cloud Instance: %s (Id: %s) - Driver: %s - IP Address: %s - Options [%d] :%s\n", server.Name, server.Id, server.Driver, server.IPAddress, num, options)
			}
			fmt.Printf("   Installation Plans: %d\n", len(network.Installations))
			for _,installation := range network.Installations {
				serverName,ok := serversMap[installation.InstanceId]
				if !ok {
					serverName = "<invalid>"
				}
				cloud := "no"
				if installation.IsCloud {
					cloud = "yes"
				}
				success := "no"
				if installation.Success {
					success = "yes"
				}
				fmt.Printf("      Plan: Id: %s - Instance: %s [Id: %s] - Success: %s - Cloud: %s - Envoronment : %s  Role: %s  Type: %s\n", installation.Id, serverName, installation.InstanceId, success, cloud, installation.Environment, installation.Role, installation.Type)
			}
		}
	}
	return Response{
		Message: "Success",
		Status: true,}, nil
}
