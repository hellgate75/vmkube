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
	Host              string  `json:"swarmhost",xml:"swarmhost"`
	UseAddress        bool    `json:"useaddress",xml:"useaddress"`
	UseDiscovery      bool    `json:"usediscovery",xml:"usediscovery"`
	UseExperimental   bool    `json:"useexperimental",xml:"useexperimental"`
	IsMaster          bool    `json:"ismaster",xml:"ismaster"`
	Image             string  `json:"swarmimage",xml:"swarmimage"`
	JoinOpts        []string  `json:"joinopts",xml:"joinopts"`
	Strategy          string  `json:"swarmstrategy",xml:"swarmstrategy"`
	TLSSan          []string  `json:"swarmtlssan",xml:"swarmtlssan"`
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
	Environment       []string  `json:"environment",xml:"environment"`
	InsecureRegistry  []string  `json:"insecureregistry",xml:"insecureregistry"`
	StorageDriver       string  `json:"storagedriver",xml:"storagedriver"`
	InstallURL          string  `json:"installurl",xml:"installurl"`
	Labels            []string  `json:"labels",xml:"labels"`
	Options           []string  `json:"options",xml:"options"`
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
	Id      	  	  		 int    `json:"uid",xml:"uid"`
	Name    		  		string    `json:"name",xml:"name"`
	Roles   				[]string    `json:"roles",xml:"roles"`
	Driver   		 			string    `json:"driver",xml:"driver"`
	Memory      				int     `json:"memory",xml:"memory"`
	Cpus        				int     `json:"cpus",xml:"cpus"`
	Disk    						int 		`json:"disksize",xml:"disksize"`
	Swarm   ProjectSwarmOpt     `json:"swarm",xml:"swarm"`
	Engine 	ProjectEngineOpt    `json:"engine",xml:"engine"`
	OSType    			string      `json:"ostype",xml:"ostype"`
	OSVersion 			string      `json:"osversion",xml:"osversion"`
	NoShare     			bool      `json:"noshare",xml:"noshare"`
	Options   	[][]string  		`json:"options",xml:"options"`
	Hostname  			string      `json:"hostname",xml:"hostname"`
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
	Id          int       `json:"uid",xml:"uid"`
	Type      string      `json:"type",xml:"type"`
	Hostname  string      `json:"hostname",xml:"hostname"`
	Options   [][]string  `json:"options",xml:"options"`
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
	Id          		string          `json:"uid",xml:"uid"`
	ServerId    		int             `json:"serverid",xml:"serverid"`
	IsCloud     		bool            `json:"iscloud",xml:"iscloud"`
	Type        	InstallationRole  `json:"type",xml:"type"`
	Environment 				SystemRole  `json:"environment",xml:"environment"`
	Role        ProjectEnvironment	`json:"role",xml:"role"`
	MainCommandSet  		CommandSet	`json:"cmdset",xml:"cmdset"`
	MainCommandRef		  string			`json:"commandref",xml:"commandref"`
	ProvisionCommandSet CommandSet	`json:"provisionset",xml:"provisionset"`
	ProvisionCommandRef	string			`json:"provisionref",xml:"provisionref"`
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
	Id            string             			`json:"uid",xml:"uid"`
	Name          string          				`json:"name",xml:"name"`
	Servers       []ProjectServer        	`json:"servers",xml:"servers"`
	CServers      []ProjectCloudServer   	`json:"cloudservers",xml:"cloudservers"`
	Installations []InstallationPlan  		`json:"installations",xml:"installations"`
}

/*
Describe domain options, contains
	
  * Id            (string)            Unique Identifier
  
  * Name          (string)          	Domain Name
  
  * Networks      ([]ProjectNetwork)	Networks List
*/
type ProjectDomain struct {
	Id          string            `json:"uid",xml:"uid"`
	Name        string          	`json:"name",xml:"name"`
	Networks    []ProjectNetwork	`json:"networks",xml:"networks"`
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
	Id          string            `json:"uid",xml:"uid"`
	Name     		string  					`json:"name",xml:"name"`
	Domains     []ProjectDomain  	`json:"domains",xml:"domains"`
	Created     time.Timer      	`json:"creation",xml:"creation"`
	Modified    time.Timer      	`json:"modified",xml:"modified"`
	Errors      bool      `json:"haserrors",xml:"haserrors"`
	LastMessage string    `json:"lastmessage",xml:"lastmessage"`
}
