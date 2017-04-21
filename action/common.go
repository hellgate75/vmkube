package action

import (
	"vmkube/model"
	"errors"
	"io/ioutil"
	"encoding/json"
	"encoding/xml"
	"vmkube/vmio"
)

type Response struct {
	Status	bool
	Message	string
}

type ExportImportDomains struct {
	Domains []model.ProjectDomain
}

var ImportDomainSample ExportImportDomains = ExportImportDomains{
	Domains: []model.ProjectDomain{
		vmio.DomainSample,
		vmio.DomainSample,
	},
}

func UserImportDomains(File string, Format string) (ExportImportDomains, error) {
	domains := ExportImportDomains{
		Domains: []model.ProjectDomain{},
	}
	if ! model.ExistsFile(File) {
		return  domains, errors.New("File "+File+" doesn't exist!!")
	}
	if Format != "json" && Format != "xml" {
		return  domains, errors.New("Format "+Format+" not available!!")
	}
	byteArray, err := ioutil.ReadFile(File)
	if err != nil {
		return  domains, err
	}
	if Format == "json" {
		err = json.Unmarshal(byteArray, &domains)
	} else  {
		err = xml.Unmarshal(byteArray, &domains)
	}
	return domains, err
}

type DomainReference struct {
	DomainId      string
	DomainName    string
}

type ExportImportNetwork struct {
	Domain      DomainReference
	Networks    []model.ProjectNetwork
}

var ImportNetworkSample ExportImportNetwork = ExportImportNetwork{
	Domain:      DomainReference{
		DomainId: NewUUIDString(),
		DomainName: vmio.DomainSample.Name,
	},
	Networks:    []model.ProjectNetwork{
		vmio.NetworkSample,
		vmio.NetworkSample,
	},
}

func UserImportNetworks(File string, Format string) ([]ExportImportNetwork, error) {
	networks := make([]ExportImportNetwork, 0)
	if ! model.ExistsFile(File) {
		return networks, errors.New("File "+File+" doesn't exist!!")
	}
	if Format != "json" && Format != "xml" {
		return networks, errors.New("Format "+Format+" not available!!")
	}
	byteArray, err := ioutil.ReadFile(File)
	if err != nil {
		return networks, err
	}
	if Format == "json" {
		err = json.Unmarshal(byteArray, &networks)
	} else  {
		err = xml.Unmarshal(byteArray, &networks)
	}
	return networks, err
}

type NetworkReference struct {
	DomainId        string
	DomainName      string
	NetworkId       string
	NetworkName     string
}

type ExportImportLocalServers struct {
	Network      NetworkReference
	Servers      []model.ProjectServer
}

var ImportLocalServerSample ExportImportLocalServers = ExportImportLocalServers{
	Network:      NetworkReference{
		DomainId:     NewUUIDString(),
		DomainName:   vmio.DomainSample.Name,
		NetworkId:    NewUUIDString(),
		NetworkName:  vmio.NetworkSample.Name,
	},
	Servers:    []model.ProjectServer{
		vmio.ServerSample,
		vmio.ServerSample,
	},
}

func UserImportLocalServers(File string, Format string) ([]ExportImportLocalServers, error) {
	servers := make([]ExportImportLocalServers, 0)
	if ! model.ExistsFile(File) {
		return servers, errors.New("File "+File+" doesn't exist!!")
	}
	if Format != "json" && Format != "xml" {
		return servers, errors.New("Format "+Format+" not available!!")
	}
	byteArray, err := ioutil.ReadFile(File)
	if err != nil {
		return servers, err
	}
	if Format == "json" {
		err = json.Unmarshal(byteArray, &servers)
	} else  {
		err = xml.Unmarshal(byteArray, &servers)
	}
	return servers, err
}

type ExportImportCloudServers struct {
	Network      NetworkReference
	Servers      []model.ProjectCloudServer
}

var ImportCloudServerSample ExportImportCloudServers = ExportImportCloudServers{
	Network:      NetworkReference{
		DomainId:     NewUUIDString(),
		DomainName:   vmio.DomainSample.Name,
		NetworkId:    NewUUIDString(),
		NetworkName:  vmio.NetworkSample.Name,
	},
	Servers:    []model.ProjectCloudServer{
		vmio.CloudServerSample,
		vmio.CloudServerSample,
	},
}

func UserImportCloudServers(File string, Format string) ([]ExportImportCloudServers, error) {
	servers := make([]ExportImportCloudServers, 0)
	if ! model.ExistsFile(File) {
		return servers, errors.New("File "+File+" doesn't exist!!")
	}
	if Format != "json" && Format != "xml" {
		return servers, errors.New("Format "+Format+" not available!!")
	}
	byteArray, err := ioutil.ReadFile(File)
	if err != nil {
		return servers, err
	}
	if Format == "json" {
		err = json.Unmarshal(byteArray, &servers)
	} else  {
		err = xml.Unmarshal(byteArray, &servers)
	}
	return servers, err
}

type ExportImportPlans struct {
	Network      NetworkReference
	Plans       []model.InstallationPlan
}

var InstallationPlanSample model.InstallationPlan = model.InstallationPlan{
	Id: "",
	Environment: model.KubernetesEnv,
	IsCloud: false,
	ServerId: NewUUIDString(),
	MainCommandRef: "https://github.com/myrepo/something/mycommand.git",
	MainCommandSet: model.AnsibleCmdSet,
	ProvisionCommandRef: "https://site.to.my.commands/something/mycommand.tgz",
	ProvisionCommandSet: model.HelmCmdSet,
	Role: model.MasterRole,
	Type: model.HostRole,
}

var InstallationPlanSample2 model.InstallationPlan = model.InstallationPlan{
	Id: "",
	Environment: model.CattleEnv,
	IsCloud: true,
	ServerId: NewUUIDString(),
	MainCommandRef: "https://github.com/myrepo/something/mycommand.git",
	MainCommandSet: model.AnsibleCmdSet,
	ProvisionCommandRef: "https://site.to.my.commands/something/mycommand.tgz",
	ProvisionCommandSet: model.VirtKubeCmdSet,
	Role: model.StandAloneRole,
	Type: model.ServerRole,
}

var ImportPlansSample ExportImportPlans = ExportImportPlans{
	Network:      NetworkReference{
		DomainId:     NewUUIDString(),
		DomainName:   vmio.DomainSample.Name,
		NetworkId:    NewUUIDString(),
		NetworkName:  vmio.NetworkSample.Name,
	},
	Plans:    []model.InstallationPlan{
		InstallationPlanSample,
		InstallationPlanSample2,
	},
}

func UserImportPlans(File string, Format string) ([]ExportImportPlans, error) {
	servers := make([]ExportImportPlans, 0)
	if ! model.ExistsFile(File) {
		return servers, errors.New("File "+File+" doesn't exist!!")
	}
	if Format != "json" && Format != "xml" {
		return servers, errors.New("Format "+Format+" not available!!")
	}
	byteArray, err := ioutil.ReadFile(File)
	if err != nil {
		return servers, err
	}
	if Format == "json" {
		err = json.Unmarshal(byteArray, &servers)
	} else  {
		err = xml.Unmarshal(byteArray, &servers)
	}
	return servers, err
}
