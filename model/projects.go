package model

import "time"


/*
Describe Swarm Cluster feature, contains
	
  * UseAddress      (bool)      Use Machine IP Address
  
  * UseDiscovery    (bool)      Use Swarm Discovery Option
  
  * UseExperimental (bool)      Use Docker Experimental Features
  
  * IsMaster        (bool)      Is Swarm Master
  
  * Image           (string)    Swarm Image (e.g.: smarm:latest)
  
  * JoinOpts        ([]string)  Swarm Join Options
  
  * Strategy        (string)    Swarm Strategy
  
  * TLSSan          ([]string)  Swarm TLS SAN Options
*/
type ProjectSwarmOpt struct {
	Host              string  `json:"Host",xml:"Host"`
	UseAddress        bool    `json:"UseAddress",xml:"UseAddress"`
	UseDiscovery      bool    `json:"UseDiscovery",xml:"UseDiscovery"`
	UseExperimental   bool    `json:"UseExperimental",xml:"UseExperimental"`
	IsMaster          bool    `json:"IsMaster",xml:"IsMaster"`
	Image             string  `json:"Image",xml:"Image"`
	JoinOpts        []string  `json:"JoinOpts",xml:"JoinOpts"`
	Strategy          string  `json:"Strategy",xml:"Strategy"`
	TLSSan          []string  `json:"TLSSan",xml:"TLSSan"`
}

/*
Describe Docker Engine options, contains
	
  * Environment       ([]string)  Environment variables
  
  * InsecureRegistry  ([]string)  Insecure Registry Options
  
  * StorageDriver     (string)    Storage Driver
  
  * InstallURL        (string)    Docker Install URL (e.g.: https://get.docker.com)
  
  * Labels            ([]string)  Engine Labels
  
  * Options           ([]string)  Engine Options
*/
type ProjectEngineOpt struct {
	Environment       []string  `json:"Environment",xml:"Environment"`
	InsecureRegistry  []string  `json:"InsecureRegistry",xml:"InsecureRegistry"`
	StorageDriver       string  `json:"StorageDriver",xml:"StorageDriver"`
	InstallURL          string  `json:"InstallURL",xml:"InstallURL"`
	Labels            []string  `json:"Labels",xml:"Labels"`
	Options           []string  `json:"Options",xml:"Options"`
}


/*
Describe Server options, contains
	
  * Id        (int)        Unique Identifier
  
  * Name      (string)     Server Local Name
  
  * Driver    (string)     Server Driver (virtualbox,vmware,hyperv) ref: https://docs.docker.com/machine/drivers/
  
  * Memory    (int)        Memory Size MB
  
  * Cpus      (int)        Number of Logical Cores
  
  * Swarm     (SwarmOpt)   Swarm Options
  
  * Engine    (EngineOpt)  Engine Options
  
  * OSType    (string)     Machine OS Type (ref: https://docs.docker.com/machine/drivers/os-base/)
  
  * OSVersion (string)     Machine OS Version (ref: https://docs.docker.com/machine/drivers/os-base/)
  
  * NoShare   (string)     Do not mount home as shared folder
  
  * Options   ([][]string) Specific vendor option in format key value pairs array without signs (e.g.: --)
  
  * Hostname  (string)     Logical Server Hostname
*/
type ProjectServer struct {
	Id      	  	  		 int    `json:"Id",xml:"Id"`
	Name    		  		string    `json:"Name",xml:"Name"`
	Roles   				[]string    `json:"Roles",xml:"Roles"`
	Driver   		 			string    `json:"Driver",xml:"Driver"`
	Memory      				int     `json:"Memory",xml:"Memory"`
	Cpus        				int     `json:"Cpus",xml:"Cpus"`
	DiskSize 						int 		`json:"DiskSize",xml:"DiskSize"`
	Swarm   ProjectSwarmOpt     `json:"Swarm",xml:"Swarm"`
	Engine 	ProjectEngineOpt    `json:"Engine",xml:"Engine"`
	OSType    			string      `json:"OSType",xml:"OSType"`
	OSVersion 			string      `json:"OSVersion",xml:"OSVersion"`
	NoShare     			bool      `json:"NoShare",xml:"NoShare"`
	Options   	[][]string  		`json:"Options",xml:"Options"`
	Hostname  			string      `json:"Hostname",xml:"Hostname"`
}

/*
Describe Cloud Server options, contains
	
  * Id        (int)         Unique Identifier
  
  * Type      (string)      Cloud Server Type ()
  
  * Hostname  (string)      Logical Server Hostname
  
  * Options   ([][]string)  Cloud Server Options
  
	Refers to : https://docs.docker.com/machine/drivers/
*/
type ProjectCloudServer struct {
	Id          int       `json:"Id",xml:"Id"`
	Type      string      `json:"Type",xml:"Type"`
	Hostname  string      `json:"Hostname",xml:"Hostname"`
	Options   [][]string  `json:"Options",xml:"Options"`
}

/*
Describe Installation Role Enum
	
  * Server  Take Part to installation cluster as Main Server
  
  * Host    Take part to installation cluster as simple host
*/
type InstallationRole string

const(
	ServerRole InstallationRole = "Server";
	HostRole   InstallationRole = "Host";
)

/*
Describe System Role Enum
	
  * StandAlone      StandAlone Server Unit
  
  * Master          Master Server in a cluster
  
  * Slave           Dependant Server in a cluster
  
  * ClusterMember  Peer Role in a cluster
*/
type SystemRole string

const(
	StandAloneRole      SystemRole = "Stand-Alone";
	MasterRole          SystemRole = "Master";
	SlaveRole           SystemRole = "Slave";
	ClusterMemberRole  	SystemRole = "Cluster-Memeber";
)

/*
Describe Environment Project Environment Enum (Rancher OS)
	
  * Cattle         Cattle Environment
  
  * Kubernetes     Kubernetes Environment
  
  * Mesos          Mesos Environment
  
  * Swarm          Swarm Environment
  
  * Custom         Custom Environment
*/
type ProjectEnvironment string

const(
	CattleEnv        ProjectEnvironment = "Cattle";
	KubernetesEnv    ProjectEnvironment = "Kubernetes";
	MesosEnv         ProjectEnvironment = "Mesos";
	SwarmEnv         ProjectEnvironment = "Swarm";
	CustomEnv        ProjectEnvironment = "Custom";
)

/*
Describe Command Set Type Enum (Rancher OS)

  * VirtKube       Virtual Kube Command Set

  * Ansible        Ansible Command Set

  * Helm	        	Helm Command Set
*/
type CommandSet string

const(
	VirtKubeCmdSet   CommandSet = "VirtKube";
	AnsibleCmdSet    CommandSet = "Ansible";
	HelpCmdSet       CommandSet = "Helm";
)


/*
Describe Server Installation Plan, contains
	
  * Id          	(string)              Unique Identifier
  
  * ServerId    	(int)               	Target Server Id
  
  * IsCloud     	(bool)              	Is A Cloud Server
  
  * Type        	(InstallationRole)  	Installation Type
  
  * Environment (ProjectEnvironment)  Installation Environment (Rancher)
  
  * Role        				(SystemRole)	Installation Role

  * MainCommandSet      (CommandSet)	Command Set used for installation

  * MainCommandRef      (string)			Location of installation commands

  * ProvisionCommandSet (CommandSet)	Command Set used for provisioning

  * ProvisionCommandRef (string)			Location of provision commands
*/
type InstallationPlan struct {
	Id          		string          `json:"Id",xml:"Id"`
	ServerId    		int             `json:"ServerId",xml:"ServerId"`
	IsCloud     		bool            `json:"IsCloud",xml:"IsCloud"`
	Type        	InstallationRole  `json:"Type",xml:"Type"`
	Environment 				SystemRole  `json:"Environment",xml:"Environment"`
	Role        ProjectEnvironment	`json:"Role",xml:"Role"`
	MainCommandSet  		CommandSet	`json:"MainCommandSet",xml:"MainCommandSet"`
	MainCommandRef		  string			`json:"MainCommandRef",xml:"MainCommandRef"`
	ProvisionCommandSet CommandSet	`json:"ProvisionCommandSet",xml:"ProvisionCommandSet"`
	ProvisionCommandRef	string			`json:"ProvisionCommandRef",xml:"ProvisionCommandRef"`
}

/*
Describe Network options, contains
	
  * Id            (string)             		Unique Identifier
  
  * Name          (string)          			Network Name
  
  * Servers       ([]ProjectServer)       Server List
  
  * CServers      ([]ProjectCloudServer)  Cloud Server List
  
  * Installations ([]InstallationPlan)  	Server Installation Plans
*/
type ProjectNetwork struct {
	Id            string             			`json:"Id",xml:"Id"`
	Name          string          				`json:"Name",xml:"Name"`
	Servers       []ProjectServer        	`json:"Servers",xml:"Servers"`
	CServers      []ProjectCloudServer   	`json:"CServers",xml:"CServers"`
	Installations []InstallationPlan  		`json:"Installations",xml:"Installations"`
}

/*
Describe domain options, contains
	
  * Id            (string)            Unique Identifier
  
  * Name          (string)          	Domain Name
  
  * Networks      ([]ProjectNetwork)	Networks List
*/
type ProjectDomain struct {
	Id          string            `json:"Id",xml:"Id"`
	Name        string          	`json:"Name",xml:"Name"`
	Networks    []ProjectNetwork	`json:"Networks",xml:"Networks"`
}

/*
Describe Project, contains

  * Id          (string)            Unique Identifier

  * Name      	(string)  	 				Project Name

  * Domains      ([]ProjectDomain)  List Of Domains

  * State        (State)     				Creation State

  * Created      (time.Timer)      	Creation Date
	
  * Modified      (time.Timer)     	Last Modification Date
	
  * Errors       (bool)      				Error State
	
  * LastMessage  (string)    				Last Alternation Message
*/
type Project struct {
	Id          string            `json:"Id",xml:"Id"`
	Name     		string  					`json:"Name",xml:"Name"`
	Domains     []ProjectDomain  	`json:"Domains",xml:"Domains"`
	Created     time.Timer      	`json:"Created",xml:"Created"`
	Modified    time.Timer      	`json:"Modified",xml:"Modified"`
	Errors      bool      				`json:"Errors",xml:"Errors"`
	LastMessage string    				`json:"LastMessage",xml:"LastMessage"`
}

/*
Describe Project State in Index, contains

  * Id          (string)            Unique Identifier

  * Name      	(string)  	 				Project Name

  * Active      (bool)						  Active State of Project
*/
type ProjectsDescriptor struct {
	Id          string            `json:"Id",xml:"Id"`
	Name     		string  					`json:"Name",xml:"Name"`
	Active     bool  							`json:"Active",xml:"Active"`
}

/*
Describe Projects Index, contains

  * Projects    ([]ProjectsDescriptor)     Projects indexed in VMKube
*/
type ProjectsIndex struct {
	Projects		[]ProjectsDescriptor 		`json:"Projects",xml:"Projects"`
}
