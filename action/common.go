package action

import (
	"vmkube/model"
	"errors"
	"io/ioutil"
	"encoding/json"
	"encoding/xml"
)

type Response struct {
	Status	bool
	Message	string
}

type ExportImportDomains struct {
	Domains []model.ProjectDomain
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
