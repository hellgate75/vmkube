package common

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"github.com/hellgate75/vmkube/model"
	"github.com/hellgate75/vmkube/vmio"
)

type ExportImportDomains struct {
	Domains []model.MachineDomain
}

var ImportDomainSample ExportImportDomains = ExportImportDomains{
	Domains: []model.MachineDomain{
		vmio.DomainSample,
		vmio.DomainSample,
	},
}

func UserImportDomains(File string, Format string) (ExportImportDomains, error) {
	domains := ExportImportDomains{
		Domains: []model.MachineDomain{},
	}
	if !model.ExistsFile(File) {
		return domains, errors.New("File " + File + " doesn't exist!!")
	}
	if Format != "json" && Format != "xml" {
		return domains, errors.New("Format " + Format + " not available!!")
	}
	byteArray, err := ioutil.ReadFile(File)
	if err != nil {
		return domains, err
	}
	if Format == "json" {
		err = json.Unmarshal(byteArray, &domains)
	} else {
		err = xml.Unmarshal(byteArray, &domains)
	}
	return domains, err
}

type DomainReference struct {
	DomainId   string
	DomainName string
}

type ExportImportNetwork struct {
	Domain   DomainReference
	Networks []model.MachineNetwork
}

var ImportNetworkSample ExportImportNetwork = ExportImportNetwork{
	Domain: DomainReference{
		DomainId:   NewUUIDString(),
		DomainName: vmio.DomainSample.Name,
	},
	Networks: []model.MachineNetwork{
		vmio.NetworkSample,
		vmio.NetworkSample,
	},
}

func UserImportNetworks(File string, Format string) ([]ExportImportNetwork, error) {
	networks := make([]ExportImportNetwork, 0)
	if !model.ExistsFile(File) {
		return networks, errors.New("File " + File + " doesn't exist!!")
	}
	if Format != "json" && Format != "xml" {
		return networks, errors.New("Format " + Format + " not available!!")
	}
	byteArray, err := ioutil.ReadFile(File)
	if err != nil {
		return networks, err
	}
	if Format == "json" {
		err = json.Unmarshal(byteArray, &networks)
	} else {
		err = xml.Unmarshal(byteArray, &networks)
	}
	return networks, err
}

type NetworkReference struct {
	DomainId    string
	DomainName  string
	NetworkId   string
	NetworkName string
}

type ExportImportLocalMachines struct {
	Network  NetworkReference
	Machines []model.LocalMachine
}

var ImportLocalMachineSample ExportImportLocalMachines = ExportImportLocalMachines{
	Network: NetworkReference{
		DomainId:    NewUUIDString(),
		DomainName:  vmio.DomainSample.Name,
		NetworkId:   NewUUIDString(),
		NetworkName: vmio.NetworkSample.Name,
	},
	Machines: []model.LocalMachine{
		vmio.MachineSample,
		vmio.MachineSample,
	},
}

func UserImportLocalMachines(File string, Format string) ([]ExportImportLocalMachines, error) {
	machines := make([]ExportImportLocalMachines, 0)
	if !model.ExistsFile(File) {
		return machines, errors.New("File " + File + " doesn't exist!!")
	}
	if Format != "json" && Format != "xml" {
		return machines, errors.New("Format " + Format + " not available!!")
	}
	byteArray, err := ioutil.ReadFile(File)
	if err != nil {
		return machines, err
	}
	if Format == "json" {
		err = json.Unmarshal(byteArray, &machines)
	} else {
		err = xml.Unmarshal(byteArray, &machines)
	}
	return machines, err
}

type ExportImportCloudMachines struct {
	Network  NetworkReference
	Machines []model.CloudMachine
}

var ImportCloudMachineSample ExportImportCloudMachines = ExportImportCloudMachines{
	Network: NetworkReference{
		DomainId:    NewUUIDString(),
		DomainName:  vmio.DomainSample.Name,
		NetworkId:   NewUUIDString(),
		NetworkName: vmio.NetworkSample.Name,
	},
	Machines: []model.CloudMachine{
		vmio.CloudMachineSample,
		vmio.CloudMachineSample,
	},
}

func UserImportCloudMachines(File string, Format string) ([]ExportImportCloudMachines, error) {
	machines := make([]ExportImportCloudMachines, 0)
	if !model.ExistsFile(File) {
		return machines, errors.New("File " + File + " doesn't exist!!")
	}
	if Format != "json" && Format != "xml" {
		return machines, errors.New("Format " + Format + " not available!!")
	}
	byteArray, err := ioutil.ReadFile(File)
	if err != nil {
		return machines, err
	}
	if Format == "json" {
		err = json.Unmarshal(byteArray, &machines)
	} else {
		err = xml.Unmarshal(byteArray, &machines)
	}
	return machines, err
}

type ExportImportPlans struct {
	Network NetworkReference
	Plans   []model.InstallationPlan
}

var InstallationPlanSample model.InstallationPlan = model.InstallationPlan{
	Id:                  "",
	Environment:         model.KubernetesEnv,
	IsCloud:             false,
	MachineId:           NewUUIDString(),
	MainCommandRef:      "https://github.com/myrepo/something/mycommand.git",
	MainCommandSet:      model.AnsibleCmdSet,
	ProvisionCommandRef: "https://site.to.my.commands/something/mycommand.tgz",
	ProvisionCommandSet: model.HelmCmdSet,
	Role:                model.MasterRole,
	Type:                model.HostDeployment,
}

var InstallationPlanSample2 model.InstallationPlan = model.InstallationPlan{
	Id:                  "",
	Environment:         model.CattleEnv,
	IsCloud:             true,
	MachineId:           NewUUIDString(),
	MainCommandRef:      "https://github.com/myrepo/something/mycommand.git",
	MainCommandSet:      model.AnsibleCmdSet,
	ProvisionCommandRef: "https://site.to.my.commands/something/mycommand.tgz",
	ProvisionCommandSet: model.VirtKubeCmdSet,
	Role:                model.StandAloneRole,
	Type:                model.MachineDeployment,
}

var ImportPlansSample ExportImportPlans = ExportImportPlans{
	Network: NetworkReference{
		DomainId:    NewUUIDString(),
		DomainName:  vmio.DomainSample.Name,
		NetworkId:   NewUUIDString(),
		NetworkName: vmio.NetworkSample.Name,
	},
	Plans: []model.InstallationPlan{
		InstallationPlanSample,
		InstallationPlanSample2,
	},
}

func UserImportPlans(File string, Format string) ([]ExportImportPlans, error) {
	machines := make([]ExportImportPlans, 0)
	if !model.ExistsFile(File) {
		return machines, errors.New("File " + File + " doesn't exist!!")
	}
	if Format != "json" && Format != "xml" {
		return machines, errors.New("Format " + Format + " not available!!")
	}
	byteArray, err := ioutil.ReadFile(File)
	if err != nil {
		return machines, err
	}
	if Format == "json" {
		err = json.Unmarshal(byteArray, &machines)
	} else {
		err = xml.Unmarshal(byteArray, &machines)
	}
	return machines, err
}
