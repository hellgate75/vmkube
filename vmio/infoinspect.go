package vmio

import (
	"vmkube/model"
)

type ProjectItem int

const(
	None ProjectItem = iota;
	ServerElement ProjectItem = iota + 1;
	CloudServerElement ProjectItem = iota + 1;
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

var ServerSample model.ProjectServer = model.ProjectServer{
	Id: "",
	Name: "MyServer",
	Driver: "virtualbox",
	Hostname: "myserver",
	Cpus: 2,
	Memory: 1024,
	DiskSize: 80,
	NoShare: true,
	Options: [][]string{
		[]string{"myoption", "myvalue"},
	},
	OSType: "rancheros",
	OSVersion: "0.9.0",
	Roles: []string{"server","master","rancher-host","rancher-server"},
	Engine: model.ProjectEngineOpt{
		Environment: []string{"MY=ENV-VAR=MY-VALUE"},
	},
	Swarm: model.ProjectSwarmOpt{},
}

var CloudServerSample model.ProjectCloudServer = model.ProjectCloudServer{
	Id: "",
	Name: "MyCloudServer",
	Driver: "virtualbox",
	Hostname: "myserver",
	Roles: []string{"server","master","rancher-host","rancher-server"},
	Options: [][]string{
		[]string{"my-provider-option", "my-provider-option-value"},
	},
}

var InstallationPlanSample model.InstallationPlan = model.InstallationPlan{
	Id: "",
	Environment: model.KubernetesEnv,
	IsCloud: false,
	ServerId: "MyServer",
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
	ServerId: "MyCloudServer",
	MainCommandRef: "https://github.com/myrepo/something/mycommand.git",
	MainCommandSet: model.AnsibleCmdSet,
	ProvisionCommandRef: "https://site.to.my.commands/something/mycommand.tgz",
	ProvisionCommandSet: model.VirtKubeCmdSet,
	Role: model.StandAloneRole,
	Type: model.ServerRole,
}

var NetworkSample model.ProjectNetwork = model.ProjectNetwork{
	Id: "",
	Name: "MyNetwork",
	CServers: []model.ProjectCloudServer{CloudServerSample},
	Servers: []model.ProjectServer{ServerSample},
	Installations: []model.InstallationPlan{InstallationPlanSample, InstallationPlanSample2},
	Options: [][]string{
		[]string{"my-network-option", "my-network-option-value"},
	},
}

var DomainSample model.ProjectDomain = model.ProjectDomain {
	Id: "",
	Name: "MyDomain",
	Networks: []model.ProjectNetwork{NetworkSample},
	Options: [][]string{
		[]string{"my-domain-option", "my-domain-option-value"},
	},
}

var ProjectSample model.ProjectImport = model.ProjectImport {
	Id: "",
	Name: "MyProject",
	Domains: []model.ProjectDomain{DomainSample},
}

func ListProjectTypeDefines() DefineList {
	defineList := make(DefineList, 0)
	defineList = append(defineList, TypeDefine{
		Name: "Server",
		Description: "Server Element describes Instance configuration for local scope",
		Type: ServerElement,
		Fields: []TypeDefineField{},
		Sample: ServerSample,
	})
	defineList = append(defineList, TypeDefine{
		Name: "Cloud-Server",
		Description: "Server Element describes Instance configuration for remote/cloud scope",
		Type: CloudServerElement,
		Fields: []TypeDefineField{},
		Sample: CloudServerSample,
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
		Description: "Describe Network Inrastracture, composed by Server/Cloud Server Configurations and Plans",
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