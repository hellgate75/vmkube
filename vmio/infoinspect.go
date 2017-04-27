package vmio

import (
	"vmkube/model"
)

type ProjectItem int

const(
	None ProjectItem = iota;
	MachineElement ProjectItem = iota + 1;
	CloudMachineElement ProjectItem = iota + 1;
	PlanElement ProjectItem = iota + 1;
	NetworkElement ProjectItem = iota + 1;
	DomainElement ProjectItem = iota + 1;
	ProjectElement ProjectItem = iota + 1;
)

type TypeDefineField struct {
	Name          string
	Type          string
	Mandatory     bool
	Description   string
}

type TypeDefine struct {
	Name          string
	Description   string
	Sample        interface{}
	Fields        []TypeDefineField
	Type          ProjectItem
}

type DefineList []TypeDefine

var MachineSample model.LocalMachine = model.LocalMachine{
	Id: "",
	Name: "MyMachine",
	Driver: "virtualbox",
	Hostname: "mymachine",
	Cpus: 2,
	Memory: 1024,
	DiskSize: 80,
	NoShare: true,
	Options: [][]string{
		[]string{"myoption", "myvalue"},
	},
	OSType: "rancheros",
	OSVersion: "0.9.0",
	Roles: []string{"machine","master","rancher-host","rancher-machine"},
	Engine: model.ProjectEngineOpt{
		Environment: []string{"MY=ENV-VAR=MY-VALUE"},
	},
	Swarm: model.ProjectSwarmOpt{},
}

var CloudMachineSample model.CloudMachine = model.CloudMachine{
	Id: "",
	Name: "MyCloudMachine",
	Driver: "virtualbox",
	Hostname: "mymachine",
	Roles: []string{"machine","master","rancher-host","rancher-machine"},
	Options: [][]string{
		[]string{"my-provider-option", "my-provider-option-value"},
	},
}

var InstallationPlanSample model.InstallationPlan = model.InstallationPlan{
	Id: "",
	Environment: model.KubernetesEnv,
	IsCloud: false,
	MachineId: "MyMachine",
	MainCommandRef: "https://github.com/myrepo/something/mycommand.git",
	MainCommandSet: model.AnsibleCmdSet,
	ProvisionCommandRef: "https://site.to.my.commands/something/mycommand.tgz",
	ProvisionCommandSet: model.HelmCmdSet,
	Role: model.MasterRole,
	Type: model.HostDeployment,
}

var InstallationPlanSample2 model.InstallationPlan = model.InstallationPlan{
	Id: "",
	Environment: model.CattleEnv,
	IsCloud: true,
	MachineId: "MyCloudMachine",
	MainCommandRef: "https://github.com/myrepo/something/mycommand.git",
	MainCommandSet: model.AnsibleCmdSet,
	ProvisionCommandRef: "https://site.to.my.commands/something/mycommand.tgz",
	ProvisionCommandSet: model.VirtKubeCmdSet,
	Role: model.StandAloneRole,
	Type: model.MachineDeployment,
}

var NetworkSample model.MachineNetwork = model.MachineNetwork{
	Id: "",
	Name: "MyNetwork",
	CloudMachines: []model.CloudMachine{CloudMachineSample},
	LocalMachines: []model.LocalMachine{MachineSample},
	Installations: []model.InstallationPlan{InstallationPlanSample, InstallationPlanSample2},
	Options: [][]string{
		[]string{"my-network-option", "my-network-option-value"},
	},
}

var DomainSample model.MachineDomain = model.MachineDomain{
	Id: "",
	Name: "MyDomain",
	Networks: []model.MachineNetwork{NetworkSample},
	Options: [][]string{
		[]string{"my-domain-option", "my-domain-option-value"},
	},
}

var ProjectSample model.ProjectImport = model.ProjectImport {
	Id: "",
	Name: "MyProject",
	Domains: []model.MachineDomain{DomainSample},
}

func ListProjectTypeDefines() DefineList {
	defineList := make(DefineList, 0)
	defineList = append(defineList, TypeDefine{
		Name: "Machine",
		Description: "Machine Element describes Instance configuration for local scope",
		Type: MachineElement,
		Fields: []TypeDefineField{},
		Sample: MachineSample,
	})
	defineList = append(defineList, TypeDefine{
		Name: "Cloud-Machine",
		Description: "Machine Element describes Instance configuration for remote/cloud scope",
		Type: CloudMachineElement,
		Fields: []TypeDefineField{},
		Sample: CloudMachineSample,
	})
	defineList = append(defineList, TypeDefine{
		Name: "Plan",
		Description: "Plan describes one single couple of installation and provisioning procedures for one project instance",
		Type: PlanElement,
		Fields: []TypeDefineField{},
		Sample: InstallationPlanSample,
	})
	defineList = append(defineList, TypeDefine{
		Name: "Network",
		Description: "Describe Network Inrastracture, composed by Machine/Cloud Machine Configurations and Plans",
		Type: NetworkElement,
		Fields: []TypeDefineField{},
		Sample: NetworkSample,
	})
	defineList = append(defineList, TypeDefine{
		Name: "Domain",
		Description: "Describe Domain Inrastracture, composed by Networks, and defining and order in the Infrastructure",
		Type: DomainElement,
		Fields: []TypeDefineField{},
		Sample: DomainSample,
	})
	defineList = append(defineList, TypeDefine{
		Name: "Project",
		Description: "Top Level Design Element containing multiple Domains",
		Type: ProjectElement,
		Fields: []TypeDefineField{},
		Sample: ProjectSample,
	})
	
	return defineList
}