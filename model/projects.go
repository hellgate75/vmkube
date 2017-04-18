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
	Enabled       		bool    `json:"Enabled",xml:"Enabled",mandatory:"yes",descr:"Enable Swarm Features"`
	Host              string  `json:"Host",xml:"Host",mandatory:"no",descr:"Swarm discovery Address (tcp://0.0.0.0:3376)"`
	UseAddress        bool    `json:"UseAddress",xml:"UseAddress",mandatory:"no",descr:"Use Machine IP Address"`
	DiscoveryToken    string   `json:"DiscoveryToken",xml:"DiscoveryToken",mandatory:"no",descr:"Use Swarm Discovery Option (token://<token>)"`
	UseExperimental   bool    `json:"UseExperimental",xml:"UseExperimental",mandatory:"no",descr:"Use Docker Experimental Features"`
	IsMaster          bool    `json:"IsMaster",xml:"IsMaster",mandatory:"no",descr:"Is Swarm Master"`
	Image             string  `json:"Image",xml:"Image",mandatory:"no",descr:"Swarm Image (e.g.: smarm:latest)"`
	JoinOpts        []string  `json:"JoinOpts",xml:"JoinOpts",mandatory:"no",descr:"Swarm Join Options"`
	Strategy          string  `json:"Strategy",xml:"Strategy",mandatory:"no",descr:"Swarm Strategy"`
	TLSSan          []string  `json:"TLSSan",xml:"TLSSan",mandatory:"no",descr:"Swarm TLS Specific Options (No overwrite, be sure of syntax, eg: [\"--my-tls-key my-tls-value\"])"`
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
	Environment       []string  `json:"Environment",xml:"Environment",mandatory:"no",descr:"Environment variables"`
	InsecureRegistry  []string  `json:"InsecureRegistry",xml:"InsecureRegistry",mandatory:"no",descr:"Insecure Registry Options"`
	RegistryMirror  	[]string  `json:"RegistryMirror",xml:"RegistryMirror",mandatory:"no",descr:"Registry Mirror Options"`
	StorageDriver       string  `json:"StorageDriver",xml:"StorageDriver",mandatory:"no",descr:"Storage Driver"`
	InstallURL          string  `json:"InstallURL",xml:"InstallURL",mandatory:"no",descr:"Docker Install URL (e.g.: https://get.docker.com)"`
	Labels            []string  `json:"Labels",xml:"Labels",mandatory:"no",descr:"Engine Labels"`
	Options           []string  `json:"Options",xml:"Options",mandatory:"no",descr:"Engine Options"`
}


/*
Describe Server options, contains
	
  * Id        (string)     Unique Identifier
  
  * Name      (string)     Server Local Name
  
  * Roles     ([]string)   Roles used in the deployment plan
  
  * Driver    (string)     Server Driver (virtualbox,vmware,hyperv) ref: https://docs.docker.com/machine/drivers/
  
  * Memory    (int)        Memory Size MB
  
  * Cpus      (int)        Number of Logical Cores
  
  * DiskSize  (int)        Dimension of disk root (in GB)
  
  * Swarm     (SwarmOpt)   Swarm Options
  
  * Engine    (EngineOpt)  Engine Options
  
  * OSType    (string)     Machine OS Type (ref: https://docs.docker.com/machine/drivers/os-base/)
  
  * OSVersion (string)     Machine OS Version (ref: https://docs.docker.com/machine/drivers/os-base/)
  
  * NoShare   (string)     Mount or not home as shared folder
  
  * Options   ([][]string) Specific vendor option in format key value pairs array without driver (i.e.: --<driver>-<option>), no value options are accepted
  
  * Hostname  (string)     Logical Server Hostname
*/
type ProjectServer struct {
	Id      	  	  	string    `json:"Id",xml:"Id",mandatory:"yes",descr:"Unique Identifier"`
	Name    		  		string    `json:"Name",xml:"Name",mandatory:"yes",descr:"Server Local Name"`
	Roles   				[]string    `json:"Roles",xml:"Roles",mandatory:"no",descr:"Roles used in the deployment plan"`
	Driver   		 			string    `json:"Driver",xml:"Driver",mandatory:"yes",descr:"Server Driver (virtualbox,vmware,hyperv) ref: https://docs.docker.com/machine/drivers/"`
	Memory      				int     `json:"Memory",xml:"Memory",mandatory:"no",descr:"Memory Size MB"`
	Cpus        				int     `json:"Cpus",xml:"Cpus",mandatory:"no",descr:"Number of Logical Cores"`
	DiskSize 						int 		`json:"DiskSize",xml:"DiskSize",mandatory:"no",descr:"Dimension of disk root (in GB)"`
	Swarm   ProjectSwarmOpt     `json:"Swarm",xml:"Swarm",mandatory:"no",descr:"Swarm Options"`
	Engine 	ProjectEngineOpt    `json:"Engine",xml:"Engine",mandatory:"no",descr:"Engine Options"`
	OSType    			string      `json:"OSType",xml:"OSType",mandatory:"yes",descr:"Machine OS Type (ref: https://docs.docker.com/machine/drivers/os-base/)"`
	OSVersion 			string      `json:"OSVersion",xml:"OSVersion",mandatory:"yes",descr:"Machine OS Version (ref: https://docs.docker.com/machine/drivers/os-base/)"`
	NoShare     			bool      `json:"NoShare",xml:"NoShare",mandatory:"no",descr:"Mount or not home as shared folder"`
	Options   	[][]string  		`json:"Options",xml:"Options",mandatory:"no",descr:"Specific vendor option in format key value pairs array without driver (i.e.: --<driver>-<option>), no value options are accepted"`
	Hostname  			string      `json:"Hostname",xml:"Hostname",mandatory:"yes",descr:"Logical Server Hostname"`
}

/*
Describe Cloud Server options, contains
	
  * Id        (string)      Unique Identifier
  
  * Name      (string)      Cloud Instance Name
  
  * Driver    (string)      Cloud Server Driver (amazonec2, digitalocean, azure, etc...)
  
  * Hostname  (string)      Logical Server Hostname
  
  * Roles     ([]string)    Roles used in the deployment plan
  
  * Options   ([][]string)  Cloud Server Options (vendor specific options) without  driver (i.e.: --<driver>-<option>), no value options are accepted
  
	Refers to : https://docs.docker.com/machine/drivers/
*/
type ProjectCloudServer struct {
	Id        string      `json:"Id",xml:"Id",mandatory:"yes",descr:"Unique Identifier"`
	Name      string      `json:"Name",xml:"Name",mandatory:"yes",descr:"Cloud Instance Name"`
	Driver    string      `json:"Driver",xml:"Driver",mandatory:"yes",descr:"Cloud Server Driver (amazonec2, digitalocean, azure, etc...)"`
	Hostname  string      `json:"Hostname",xml:"Hostname",mandatory:"yes",descr:"Logical Server Hostname"`
	Roles   	[]string    `json:"Roles",xml:"Roles",mandatory:"no",descr:"Roles used in the deployment plan"`
	Options   [][]string  `json:"Options",xml:"Options",mandatory:"yes",descr:"Cloud Server Options (vendor specific options) without  driver (i.e.: --<driver>-<option>), no value options are accepted"`
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
  
  * Type        	(InstallationRole)  	Installation Type (Server,Host)
  
  * Environment (ProjectEnvironment)    Installation Environment (Cattle,Kubernetes,Mesos,Swarm,Custom)
  
  * Role        				(SystemRole)	  Installation Role (Stand-Alone,Master,Slave,Cluster-Memeber)

  * MainCommandSet      (CommandSet)	  Command Set used for installation (VirtKube,Ansible,Helm)

  * MainCommandRef      (string)			  Location of installation commands (http,file,git protocols are accepted)

  * ProvisionCommandSet (CommandSet)	  Command Set used for provisioning (VirtKube,Ansible,Helm)

  * ProvisionCommandRef (string)			  Location of provision commands (http,file,git protocols are accepted)
*/
type InstallationPlan struct {
	Id          		string          `json:"Id",xml:"Id",mandatory:"yes",descr:"Unique Identifier"`
	ServerId    		string          `json:"ServerId",xml:"ServerId",mandatory:"yes",descr:"Target Server Id"`
	IsCloud     		bool            `json:"IsCloud",xml:"IsCloud",mandatory:"yes",descr:"Is A Cloud Server"`
	Type        	InstallationRole  `json:"Type",xml:"Type",mandatory:"yes",descr:"Installation Type (Server,Host)"`
	Role 				  SystemRole        `json:"Environment",xml:"Environment",mandatory:"yes",descr:"Installation Role (Stand-Alone,Master,Slave,Cluster-Memeber)"`
	Environment  ProjectEnvironment	`json:"Role",xml:"Role",mandatory:"yes",descr:"Installation Environment (Cattle,Kubernetes,Mesos,Swarm,Custom)"`
	MainCommandSet  		CommandSet	`json:"MainCommandSet",xml:"MainCommandSet",mandatory:"no",descr:"Command Set used for installation (VirtKube,Ansible,Helm)"`
	MainCommandRef		  string			`json:"MainCommandRef",xml:"MainCommandRef",mandatory:"no",descr:"Location of installation commands (http,file,git protocols are accepted)"`
	ProvisionCommandSet CommandSet	`json:"ProvisionCommandSet",xml:"ProvisionCommandSet",mandatory:"no",descr:"Command Set used for provisioning (VirtKube,Ansible,Helm)"`
	ProvisionCommandRef	string			`json:"ProvisionCommandRef",xml:"ProvisionCommandRef",mandatory:"no",descr:"Location of provision commands (http,file,git protocols are accepted)"`
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
	Id            string             			`json:"Id",xml:"Id",mandatory:"yes",descr:"Unique Identifier"`
	Name          string          				`json:"Name",xml:"Name",mandatory:"yes",descr:"Network Name"`
	Servers       []ProjectServer        	`json:"Servers",xml:"Servers",mandatory:"yes",descr:"Server List"`
	CServers      []ProjectCloudServer   	`json:"CServers",xml:"CServers",mandatory:"yes",descr:"Cloud Server List"`
	Installations []InstallationPlan  		`json:"Installations",xml:"Installations",mandatory:"no",descr:"Server Installation Plans"`
	Options     [][]string                `json:"Options",xml:"Options",mandatory:"no",descr:"Specific Network information (eg. cloud provider info or local info)"`
}

/*
Describe domain options, contains
	
  * Id            (string)            Unique Identifier
  
  * Name          (string)          	Domain Name
  
  * Networks      ([]ProjectNetwork)	Networks List

  * Options       ([][]string)        Specific Domain information (eg. cloud provider info or local info)
*/
type ProjectDomain struct {
	Id          string            `json:"Id",xml:"Id",mandatory:"yes",descr:"Unique Identifier"`
	Name        string          	`json:"Name",xml:"Name",mandatory:"yes",descr:"Domain Name"`
	Networks    []ProjectNetwork	`json:"Networks",xml:"Networks",mandatory:"yes",descr:"Networks List"`
	Options     [][]string        `json:"Options",xml:"Options",mandatory:"no",descr:"Specific Domain information (eg. cloud provider info or local info)"`
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
	Id          string            `json:"Id",xml:"Id",mandatory:"yes",descr:"Unique Identifier"`
	Name     		string  					`json:"Name",xml:"Name",mandatory:"yes",descr:"Project Name"`
	Open      	bool      				`json:"Open",xml:"Open",mandatory:"no",descr:"Writable State"`
	Domains     []ProjectDomain  	`json:"Domains",xml:"Domains",mandatory:"yes",descr:"List Of Domains"`
	Created     time.Time       	`json:"Created",xml:"Created",mandatory:"no",descr:"Creation Date"`
	Modified    time.Time       	`json:"Modified",xml:"Modified",mandatory:"no",descr:"Last Modification Date"`
	Errors      bool      				`json:"Errors",xml:"Errors",mandatory:"no",descr:"Error State"`
	LastMessage string    				`json:"LastMessage",xml:"LastMessage",mandatory:"no",descr:"Last Alternation Message"`
}

/*
Describe Project Import Model, contains

  * Id          (string)            Unique Identifier

  * Name      	(string)  	 				Project Name

  * Domains      ([]ProjectDomain)  List Of Domains
*/
type ProjectImport struct {
	Id          string            `json:"Id",xml:"Id",mandatory:"yes",descr:"Unique Identifier"`
	Name     		string  					`json:"Name",xml:"Name",mandatory:"yes",descr:"Project Name"`
	Domains     []ProjectDomain  	`json:"Domains",xml:"Domains",mandatory:"yes",descr:"List Of Domains"`
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
	Id          string            `json:"Id",xml:"Id",mandatory:"yes",descr:"Project Unique Identifier"`
	Name     		string  					`json:"Name",xml:"Name",mandatory:"yes",descr:"Project Name"`
	InfraId     string            `json:"InfraId",xml:"InfraId",mandatory:"yes",descr:"Infrastructure Unique Identifier"`
	InfraName   string  					`json:"InfraName",xml:"InfraName",mandatory:"yes",descr:"Infrastructure Project Name"`
	Active      bool  						`json:"Active",xml:"Active",mandatory:"yes",descr:"Active State of Project"`
	Open      	bool  						`json:"Open",xml:"Open",mandatory:"yes",descr:"Project Writable State"`
	Synced      bool  						`json:"Synced",xml:"Synced",mandatory:"no",descr:"Sync State of Project"`
}

/*
Describe Projects Index, contains

  * Projects    ([]ProjectsDescriptor)     Projects indexed in VMKube
*/
type ProjectsIndex struct {
	Projects		[]ProjectsDescriptor 		`json:"Projects",xml:"Projects",mandatory:"yes",descr:"Projects indexed in VMKube"`
}
