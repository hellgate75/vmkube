package model

import "time"


/*
Describe Swarm Cluster feature, contains
	
  * Enabled      		(bool)      Enable Swarm Features

  * Host      			(string)    Swarm discovery Address (tcp://0.0.0.0:3376)

  * UseAddress      (bool)      Use Machine IP Address

  * DiscoveryToken  (bool)    	Use Swarm Discovery Option (token://<token>)
  
  * UseExperimental (bool)      Use Docker Experimental Features
  
  * IsMaster        (bool)      Is Swarm Master
  
  * Image           (string)    Swarm Image (e.g.: smarm:latest)
  
  * JoinOpts        ([]string)  Swarm Join Options
  
  * Strategy        (string)    Swarm Strategy
  
  * TLSSan          ([]string)  Swarm TLS Specific Options (No overwrite, be sure of syntax, eg: ["--my-tls-key my-tls-value"])
*/
type ProjectSwarmOpt struct {
	Enabled       		bool    `json:"Enabled",xml:"Enabled",mandatory:"yes"`
	Host              string  `json:"Host",xml:"Host",mandatory:"no"`
	UseAddress        bool    `json:"UseAddress",xml:"UseAddress",mandatory:"no"`
	DiscoveryToken    string   `json:"DiscoveryToken",xml:"DiscoveryToken",mandatory:"no"`
	UseExperimental   bool    `json:"UseExperimental",xml:"UseExperimental",mandatory:"no"`
	IsMaster          bool    `json:"IsMaster",xml:"IsMaster",mandatory:"no"`
	Image             string  `json:"Image",xml:"Image",mandatory:"no"`
	JoinOpts        []string  `json:"JoinOpts",xml:"JoinOpts",mandatory:"no"`
	Strategy          string  `json:"Strategy",xml:"Strategy",mandatory:"no"`
	TLSSan          []string  `json:"TLSSan",xml:"TLSSan",mandatory:"no"`
}

/*
Describe Docker Engine options, contains
	
  * Environment       ([]string)  Environment variables
  
  * InsecureRegistry  ([]string)  Insecure Registry Options

  * RegistryMirror		([]string)	Registry Mirror Options
  
  * StorageDriver     (string)    Storage Driver
  
  * InstallURL        (string)    Docker Install URL (e.g.: https://get.docker.com)
  
  * Labels            ([]string)  Engine Labels
  
  * Options           ([]string)  Engine Options
*/
type ProjectEngineOpt struct {
	Environment       []string  `json:"Environment",xml:"Environment",mandatory:"no"`
	InsecureRegistry  []string  `json:"InsecureRegistry",xml:"InsecureRegistry",mandatory:"no"`
	RegistryMirror  	[]string  `json:"RegistryMirror",xml:"RegistryMirror",mandatory:"no"`
	StorageDriver       string  `json:"StorageDriver",xml:"StorageDriver",mandatory:"no"`
	InstallURL          string  `json:"InstallURL",xml:"InstallURL",mandatory:"no"`
	Labels            []string  `json:"Labels",xml:"Labels",mandatory:"no"`
	Options           []string  `json:"Options",xml:"Options",mandatory:"no"`
}


/*
Describe Server options, contains
	
  * Id        (string)     Unique Identifier
  
  * Name      (string)     Server Local Name
  
  * Roles     ([]string)   Roles used in the deployment plan
  
  * Driver    (string)     Server Driver (virtualbox,vmware,hyperv) ref: https://docs.docker.com/machine/drivers/
  
  * Memory    (int)        Memory Size MB
  
  * Cpus      (int)        Number of Logical Cores
  
  * Swarm     (SwarmOpt)   Swarm Options
  
  * Engine    (EngineOpt)  Engine Options
  
  * OSType    (string)     Machine OS Type (ref: https://docs.docker.com/machine/drivers/os-base/)
  
  * OSVersion (string)     Machine OS Version (ref: https://docs.docker.com/machine/drivers/os-base/)
  
  * NoShare   (string)     Do not mount home as shared folder
  
  * Options   ([][]string) Specific vendor option in format key value pairs array without signs (e.g.: --<driver>)
  
  * Hostname  (string)     Logical Server Hostname
*/
type ProjectServer struct {
	Id      	  	  	string    `json:"Id",xml:"Id",mandatory:"yes"`
	Name    		  		string    `json:"Name",xml:"Name",mandatory:"yes"`
	Roles   				[]string    `json:"Roles",xml:"Roles",mandatory:"no"`
	Driver   		 			string    `json:"Driver",xml:"Driver",mandatory:"yes"`
	Memory      				int     `json:"Memory",xml:"Memory",mandatory:"no"`
	Cpus        				int     `json:"Cpus",xml:"Cpus",mandatory:"no"`
	DiskSize 						int 		`json:"DiskSize",xml:"DiskSize",mandatory:"no"`
	Swarm   ProjectSwarmOpt     `json:"Swarm",xml:"Swarm",mandatory:"no"`
	Engine 	ProjectEngineOpt    `json:"Engine",xml:"Engine",mandatory:"no"`
	OSType    			string      `json:"OSType",xml:"OSType",mandatory:"yes"`
	OSVersion 			string      `json:"OSVersion",xml:"OSVersion",mandatory:"yes"`
	NoShare     			bool      `json:"NoShare",xml:"NoShare",mandatory:"no"`
	Options   	[][]string  		`json:"Options",xml:"Options",mandatory:"no"`
	Hostname  			string      `json:"Hostname",xml:"Hostname",mandatory:"yes"`
}

/*
Describe Cloud Server options, contains
	
  * Id        (string)      Unique Identifier
  
  * Name      (string)      Cloud Instance Name
  
  * Driver    (string)      Cloud Server Driver (amazonec2, digitalocean, azure, etc...)
  
  * Hostname  (string)      Logical Server Hostname
  
  * Roles     ([]string)    Roles used in the deployment plan
  
  * Options   ([][]string)  Cloud Server Options
  
	Refers to : https://docs.docker.com/machine/drivers/
*/
type ProjectCloudServer struct {
	Id        string      `json:"Id",xml:"Id",mandatory:"yes"`
	Name      string      `json:"Name",xml:"Name",mandatory:"yes"`
	Driver    string      `json:"Driver",xml:"Driver",mandatory:"yes"`
	Hostname  string      `json:"Hostname",xml:"Hostname",mandatory:"yes"`
	Roles   	[]string    `json:"Roles",xml:"Roles",mandatory:"no"`
	Options   [][]string  `json:"Options",xml:"Options",mandatory:"yes"`
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
	HelmCmdSet       CommandSet = "Helm";
)


/*
Describe Server Installation Plan, contains
	
  * Id          	(string)              Unique Identifier
  
  * ServerId    	(string)              Target Server Id
  
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
	Id          		string          `json:"Id",xml:"Id",mandatory:"yes"`
	ServerId    		string          `json:"ServerId",xml:"ServerId",mandatory:"yes"`
	IsCloud     		bool            `json:"IsCloud",xml:"IsCloud",mandatory:"yes"`
	Type        	InstallationRole  `json:"Type",xml:"Type",mandatory:"yes"`
	Environment 				SystemRole  `json:"Environment",xml:"Environment",mandatory:"yes"`
	Role        ProjectEnvironment	`json:"Role",xml:"Role",mandatory:"yes"`
	MainCommandSet  		CommandSet	`json:"MainCommandSet",xml:"MainCommandSet",mandatory:"no"`
	MainCommandRef		  string			`json:"MainCommandRef",xml:"MainCommandRef",mandatory:"no"`
	ProvisionCommandSet CommandSet	`json:"ProvisionCommandSet",xml:"ProvisionCommandSet",mandatory:"no"`
	ProvisionCommandRef	string			`json:"ProvisionCommandRef",xml:"ProvisionCommandRef",mandatory:"no"`
}

/*
Describe Network options, contains
	
  * Id            (string)             		Unique Identifier
  
  * Name          (string)          			Network Name
  
  * Servers       ([]ProjectServer)       Server List
  
  * CServers      ([]ProjectCloudServer)  Cloud Server List
  
  * Installations ([]InstallationPlan)  	Server Installation Plans
  
  * Options       ([][]string)            Specific Network information (eg. cloud provider info or local info)
*/
type ProjectNetwork struct {
	Id            string             			`json:"Id",xml:"Id",mandatory:"yes"`
	Name          string          				`json:"Name",xml:"Name",mandatory:"yes"`
	Servers       []ProjectServer        	`json:"Servers",xml:"Servers",mandatory:"yes"`
	CServers      []ProjectCloudServer   	`json:"CServers",xml:"CServers",mandatory:"yes"`
	Installations []InstallationPlan  		`json:"Installations",xml:"Installations",mandatory:"no"`
	Options     [][]string                `json:"Options",xml:"Options",mandatory:"no"`
}

/*
Describe domain options, contains
	
  * Id            (string)            Unique Identifier
  
  * Name          (string)          	Domain Name
  
  * Networks      ([]ProjectNetwork)	Networks List

  * Options       ([][]string)        Specific Domain information (eg. cloud provider info or local info)
*/
type ProjectDomain struct {
	Id          string            `json:"Id",xml:"Id",mandatory:"yes"`
	Name        string          	`json:"Name",xml:"Name",mandatory:"yes"`
	Networks    []ProjectNetwork	`json:"Networks",xml:"Networks",mandatory:"yes"`
	Options     [][]string        `json:"Options",xml:"Options",mandatory:"no"`
}

/*
Describe Project, contains

  * Id          (string)            Unique Identifier

  * Name      	(string)  	 				Project Name

	* Open       	(bool)      				Writable State

  * Domains      ([]ProjectDomain)  List Of Domains

  * Created      (time.Time )      	Creation Date
	
  * Modified      (time.Time )     	Last Modification Date
	
  * Errors       (bool)      				Error State
	
  * LastMessage  (string)    				Last Alternation Message
*/
type Project struct {
	Id          string            `json:"Id",xml:"Id",mandatory:"yes"`
	Name     		string  					`json:"Name",xml:"Name",mandatory:"yes"`
	Open      	bool      				`json:"Open",xml:"Open",mandatory:"no"`
	Domains     []ProjectDomain  	`json:"Domains",xml:"Domains",mandatory:"yes"`
	Created     time.Time       	`json:"Created",xml:"Created",mandatory:"no"`
	Modified    time.Time       	`json:"Modified",xml:"Modified",mandatory:"no"`
	Errors      bool      				`json:"Errors",xml:"Errors",mandatory:"no"`
	LastMessage string    				`json:"LastMessage",xml:"LastMessage",mandatory:"no"`
}

/*
Describe Project Import Model, contains

  * Id          (string)            Unique Identifier

  * Name      	(string)  	 				Project Name

  * Domains      ([]ProjectDomain)  List Of Domains
*/
type ProjectImport struct {
	Id          string            `json:"Id",xml:"Id",mandatory:"yes"`
	Name     		string  					`json:"Name",xml:"Name",mandatory:"yes"`
	Domains     []ProjectDomain  	`json:"Domains",xml:"Domains",mandatory:"yes"`
}


/*
Describe Project State in Index, contains

  * Id          (string)            Project Unique Identifier

  * Name      	(string)  	 				Project Name

  * InfraId     (string)            Infrastructure Unique Identifier

  * InfraName   (string)  	 				Infrastructure Project Name

	* Open       	(bool)      				Project Writable State

  * Active      (bool)						  Active State of Project

  * Synced      (bool)						  Sync State of Project
*/
type ProjectsDescriptor struct {
	Id          string            `json:"Id",xml:"Id",mandatory:"yes"`
	Name     		string  					`json:"Name",xml:"Name",mandatory:"yes"`
	InfraId     string            `json:"InfraId",xml:"InfraId",mandatory:"yes"`
	InfraName   string  					`json:"InfraName",xml:"InfraName",mandatory:"yes"`
	Active      bool  						`json:"Active",xml:"Active",mandatory:"yes"`
	Open      	bool  						`json:"Open",xml:"Open",mandatory:"yes"`
	Synced      bool  						`json:"Synced",xml:"Synced",mandatory:"no"`
}

/*
Describe Projects Index, contains

  * Projects    ([]ProjectsDescriptor)     Projects indexed in VMKube
*/
type ProjectsIndex struct {
	Projects		[]ProjectsDescriptor 		`json:"Projects",xml:"Projects",mandatory:"yes"`
}
