package model

import "time"

/*
Describe Disk feature, contains
	
  * Id    (string)  Unique Identifier
  
  * Name  (string)  Local Disk Name
  
  * Type  (int)     Disk Type
  
  * Size  (int)     Disk Size in MB
*/
type Disk struct {
	Id   string  `json:"Id",xml:"Id"`
	Name string `json:"Name",xml:"Name"`
	Type    int `json:"Type",xml:"Type"`
	Size    int `json:"Size",xml:"Size"`
}

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
  
  * TLSSan          ([]string)  Swarm TLS Specific Options (No overwrite, be sure of syntax)
*/
type SwarmOpt struct {
	Enabled       		bool    `json:"Enabled",xml:"Enabled"`
	Host              string  `json:"Host",xml:"Host"`
	UseAddress        bool    `json:"UseAddress",xml:"UseAddress"`
	DiscoveryToken    string   `json:"DiscoveryToken",xml:"DiscoveryToken"`
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

  * RegistryMirror		([]string)	Registry Mirror Options

  * StorageDriver     (string)    Storage Driver
  
  * InstallURL        (string)    Docker Install URL (e.g.: https://get.docker.com)
  
  * Labels            ([]string)  Engine Labels
  
  * Options           ([]string)  Engine Options
*/
type EngineOpt struct {
	Environment       []string  `json:"Environment",xml:"Environment"`
	InsecureRegistry  []string  `json:"InsecureRegistry",xml:"InsecureRegistry"`
	RegistryMirror  	[]string  `json:"RegistryMirror",xml:"RegistryMirror"`
	StorageDriver       string  `json:"StorageDriver",xml:"StorageDriver"`
	InstallURL          string  `json:"InstallURL",xml:"InstallURL"`
	Labels            []string  `json:"Labels",xml:"Labels"`
	Options           []string  `json:"Options",xml:"Options"`
}


/*
Describe Server options, contains
	
  * Id        (string)     Unique Identifier
  
  * Name      (string)     Server Local Name
  
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
type Server struct {
	Id        string      `json:"Id",xml:"Id"`
	Name      string      `json:"Name",xml:"Name"`
	Roles   []string      `json:"Roles",xml:"Roles"`
	Driver    string      `json:"Driver",xml:"Driver"`
	Memory      int       `json:"Memory",xml:"Memory"`
	Cpus        int       `json:"Cpus",xml:"Cpus"`
	Disks     []Disk      `json:"Disks",xml:"Disks"`
	Swarm   SwarmOpt      `json:"Swarm",xml:"Swarm"`
	Engine EngineOpt      `json:"Engine",xml:"Engine"`
	OSType    string      `json:"OSType",xml:"OSType"`
	OSVersion string      `json:"OSVersion",xml:"OSVersion"`
	NoShare     bool      `json:"NoShare",xml:"NoShare"`
	Options   [][]string  `json:"Options",xml:"Options"`
	Hostname  string      `json:"Hostname",xml:"Hostname"`
}

/*
Describe Cloud Server options, contains
	
  * Id        (string)      Unique Identifier
  
  * Type      (string)      Cloud Server Type ()
  
  * Hostname  (string)      Logical Server Hostname
  
  * Options   ([][]string)  Cloud Server Options
  
	Refers to : https://docs.docker.com/machine/drivers/
*/
type CloudServer struct {
	Id        string      `json:"Id",xml:"Id"`
	Type      string      `json:"Type",xml:"Type"`
	Hostname  string      `json:"Hostname",xml:"Hostname"`
	Options   [][]string  `json:"Options",xml:"Options"`
}

/*
Describe Installation Type Enum
	
  * Server  Take Part to installation cluster as Main Server
  
  * Host    Take part to installation cluster as simple host
*/
type InstallationType int

const(
	ServerType InstallationType = iota + 1; // value 1
	HostType   InstallationType = iota + 1; // value 2
)

/*
Describe Role Type Enum
	
  * StandAlone      StandAlone Server Unit
  
  * Master          Master Server in a cluster
  
  * Slave           Dependant Server in a cluster
  
  * ClusterMember  Peer Role in a cluster
*/
type RoleType int

const(
	StandAlone      RoleType = iota + 1; // value 1
	Master          RoleType = iota + 1; // value 2
	Slave           RoleType = iota + 1; // value 3
	ClusterMember  	RoleType = iota + 1; // value 4
)

/*
Describe Environment Type Enum (Rancher OS)
	
  * Cattle         Cattle Environment
  
  * Kubernetes     Kubernetes Environment
  
  * Mesos          Mesos Environment
  
  * Swarm          Swarm Environment
  
  * Custom         Custom Environment
*/
type EnvironmentType int

const(
	Cattle        EnvironmentType = iota + 1; // value 1
	Kubernetes    EnvironmentType = iota + 1; // value 2
	Mesos         EnvironmentType = iota + 1; // value 3
	Swarm         EnvironmentType = iota + 1; // value 4
	Custom        EnvironmentType = iota + 1; // value 5
)

/*
Describe Server Installation options, contains
	
  * Id          	(string)            Unique Identifier
  
  * ServerId    	(int)               Target Server Id
  
  * IsCloud     	(bool)              Is A Cloud Server
  
  * Type        	(InstallationType)  Installation Type
  
  * Environment 	(RoleType)          Installation Environment (Rancher)
  
  * Role        	(EnvironmentType)   Installation Role

  * Plan     			(InstallationPlan) Reference to installation plan

  * LastExecution	(time.Timer) 				Last Execution Date

  * Success      	(bool)       				Success State

  * Errors       (bool)       				Error State

  * LastMessage  (string)     			Last Error Message

  * LogsPath  		(string)     			Path to Log file
*/
type Installation struct {
	Id          	string            `json:"Id",xml:"Id"`
	ServerId    	int               `json:"ServerId",xml:"ServerId"`
	IsCloud     	bool              `json:"IsCloud",xml:"IsCloud"`
	Type        	InstallationType  `json:"Type",xml:"Type"`
	Environment 	RoleType          `json:"Environment",xml:"Environment"`
	Role        	EnvironmentType   `json:"Role",xml:"Role"`
	Plan					InstallationPlan	`json:"Plan",xml:"Plan"`
	LastExecution	time.Timer 				`json:"LastExecution",xml:"LastExecution"`
	Success     	bool       				`json:"Success",xml:"Success"`
	Errors      	bool       				`json:"Errors",xml:"Errors"`
	LastMessage 	string     				`json:"LastMessage",xml:"LastMessage"`
	LogsPath 			string     				`json:"LogsPath",xml:"LogsPath"`
}

/*
Describe Network options, contains
	
  * Id            (string)          Unique Identifier
  
  * Name          (string)          Network Name
  
  * Servers       ([]Server)        Server List
  
  * CServers      ([]CloudServer)   Cloud Server List
  
  * Installations ([]Installation)  Server Installations
*/
type Network struct {
	Id          	string          `json:"Id",xml:"Id"`
	Name          string          `json:"Name",xml:"Name"`
	Servers       []Server        `json:"Servers",xml:"Servers"`
	CServers      []CloudServer   `json:"CServers",xml:"CServers"`
	Installations []Installation  `json:"Installations",xml:"Installations"`
}

/*
Describe domain options, contains
	
  * Id            (string)          Unique Identifier
  
  * Name          (string)          Domain Name
  
  * Networks      ([]Networks)      Networks List
*/
type Domain struct {
	Id          string          `json:"Id",xml:"Id"`
	Name        string          `json:"Name",xml:"Name"`
	Networks    []Network       `json:"Networks",xml:"Networks"`
}

/*
Describe Server State, contains
	
  * Id           (string)     State Unique Identifier
  
  * Hostname     (string)     Defined Host Name
  
  * IPAddresses  ([]string)   Computed IP Address
  
  * ServerId     (int)        Target Server Id
  
  * IsCloud      (bool)       Is A Cloud Server
  
  * NetworkId    (int)        Target Network Id
  
  * DomainId     (int)        Target Domain Id

  * Creation     (time.Timer) Creation Date

  * Modified     (time.Timer) Last Modification Date

  * Created      (bool)       Creation State
  
  * Altered      (bool)       Alteration State
  
  * Errors       (bool)       Error State
  
  * LastMessage  (string)     Last Alternation Message
*/
type ServerState struct {
	Id          string     `json:"Id",xml:"Id"`
	Hostname    string     `json:"Hostname",xml:"Hostname"`
	IPAddresses []string   `json:"IPAddresses",xml:"IPAddresses"`
	ServerId    int        `json:"ServerId",xml:"ServerId"`
	IsCloud     bool       `json:"IsCloud",xml:"IsCloud"`
	NetworkId   int        `json:"NetworkId",xml:"NetworkId"`
	DomainId    int        `json:"DomainId",xml:"DomainId"`
	Creation    time.Timer `json:"Creation",xml:"Creation"`
	Modified 		time.Timer `json:"Modified",xml:"Modified"`
	Created     bool       `json:"Created",xml:"Created"`
	Altered     bool       `json:"Altered",xml:"Altered"`
	Errors      bool       `json:"Errors",xml:"Errors"`
	LastMessage string     `json:"LastMessage",xml:"LastMessage"`
}

/*
Describe Network State, contains
	
  * Id           (string)         State Unique Identifier
  
  * NetworkId    (int)            Target Network Id
  
  * DomainId     (int)            Target Domain Id
  
  * Servers      ([]ServerState)  List Of Server State

  * Creation     (time.Timer)     Creation Date

  * Modified     (time.Timer)     Last Modification Date

  * Created      (bool)           Creation State
  
  * Altered      (bool)           Alteration State
  
  * Errors       (bool)           Error State
  
  * LastMessage  (string)         Last Alternation Message
*/
type NetworkState struct {
	Id          string     			`json:"Id",xml:"Id"`
	NetworkId   int             `json:"NetworkId",xml:"NetworkId"`
	DomainId    int             `json:"DomainId",xml:"DomainId"`
	Servers     []ServerState   `json:"Servers",xml:"Servers"`
	Creation    time.Timer      `json:"Creation",xml:"Creation"`
	Modified 		time.Timer      `json:"Modified",xml:"Modified"`
	Created     bool            `json:"Created",xml:"Created"`
	Altered     bool            `json:"Altered",xml:"Altered"`
	Errors      bool            `json:"Errors",xml:"Errors"`
	LastMessage string          `json:"LastMessage",xml:"LastMessage"`
}

/*
Describe Domain State, contains
	
  * Id           (string)         State Unique Identifier
  
  * DomainId     (int)            Target Domain Id
  
  * Networks     ([]NetworkState) List Of Network State

  * Creation     (time.Timer) 		Creation Date

  * Modified     (time.Timer) 		Last Modification Date

  * Created      (bool)           Creation State
  
  * Altered      (bool)           Alteration State
  
  * Errors       (bool)           Error State
  
  * LastMessage  (string)         Last Alternation Message
*/
type DomainState struct {
	Id          string     			`json:"Id",xml:"Id"`
	DomainId    int             `json:"DomainId",xml:"DomainId"`
	Networks    []NetworkState  `json:"Networks",xml:"Networks"`
	Creation    time.Timer      `json:"Creation",xml:"Creation"`
	Modified 		time.Timer      `json:"Modified",xml:"Modified"`
	Created     bool            `json:"Created",xml:"Created"`
	Altered     bool            `json:"Altered",xml:"Altered"`
	Errors      bool            `json:"Errors",xml:"Errors"`
	LastMessage string          `json:"LastMessage",xml:"LastMessage"`
}

/*
Describe State, contains

	* Id           (string)         State Unique Identifier

  * Domains      ([]DomainState)  List Of Domain State

  * Creation     (time.Timer) 		Creation Date

  * Modified     (time.Timer) 		Last Modification Date

  * Created      (bool)           Creation State
	
  * Altered      (bool)           Alteration State
	
  * Errors       (bool)           Error State
	
  * LastMessage  (string)         Last Alternation Message
*/
type State struct {
	Id          string     			`json:"Id",xml:"Id"`
	Domains     []DomainState   `json:"Domains",xml:"Domains"`
	Creation    time.Timer      `json:"Creation",xml:"Creation"`
	Modified 		time.Timer      `json:"Modified",xml:"Modified"`
	Created     bool            `json:"Created",xml:"Created"`
	Altered     bool            `json:"Altered",xml:"Altered"`
	Errors      bool            `json:"Errors",xml:"Errors"`
	LastMessage string          `json:"LastMessage",xml:"LastMessage"`
}

/*
Describe Entire Infrastructure, contains

	* Id          (string)      Infrastructure Unique Identifier

	* ProjectId   (string)      Related Project Unique Identifier

  * Name      	(string)  	 	Infrastructure Name

  * Domains      ([]Domain)  	List Of Domains

  * State        (State)     	Creation State

  * Creation     (time.Timer) Creation Date

  * Modified     (time.Timer) Last Modification Date

  * Created      (bool)      	Creation State
	
  * Altered      (bool)      	Alteration State
	
  * Errors       (bool)      	Error State
	
  * LastMessage  (string)    	Last Alternation Message
*/
type Infrastructure struct {
	Id          string     				`json:"Id",xml:"Id"`
	ProjectId   string     				`json:"ProjectId",xml:"ProjectId"`
	Name     		string  					`json:"Name",xml:"Name"`
	Domains     []Domain  				`json:"Domains",xml:"Domains"`
	State       State     				`json:"State",xml:"State"`
	Creation    time.Timer      	`json:"Creation",xml:"Creation"`
	Modified 		time.Timer      	`json:"Modified",xml:"Modified"`
	Created     bool      				`json:"Created",xml:"Created"`
	Altered     bool      				`json:"Altered",xml:"Altered"`
	Errors      bool      				`json:"Errors",xml:"Errors"`
	LastMessage string    				`json:"LastMessage",xml:"LastMessage"`
}
