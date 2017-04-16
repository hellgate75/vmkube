package model

/*
Describe Disk feature, contains
	
  * Id    (int)     Unique Identifier
  
  * Name  (string)  Local Disk Name
  
  * Type  (int)     Disk Type
  
  * Size  (int)     Disk Size in MB
*/
type Disk struct {
	Id      int `json:"uid",xml:"uid"`
	Name string `json:"name",xml:"name"`
	Type    int `json:"type",xml:"type"`
	Size    int `json:"size",xml:"size"`
}

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
type SwarmOpt struct {
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
type EngineOpt struct {
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
type Server struct {
	Id          int       `json:"uid",xml:"uid"`
	Name      string      `json:"name",xml:"name"`
	Roles   []string      `json:"roles",xml:"roles"`
	Driver    string      `json:"driver",xml:"driver"`
	Memory      int       `json:"memory",xml:"memory"`
	Cpus        int       `json:"cpus",xml:"cpus"`
	Disks     []Disk      `json:"disks",xml:"disks"`
	Swarm   SwarmOpt      `json:"swarm",xml:"swarm"`
	Engine EngineOpt      `json:"engine",xml:"engine"`
	OSType    string      `json:"ostype",xml:"ostype"`
	OSVersion string      `json:"osversion",xml:"osversion"`
	NoShare     bool      `json:"noshare",xml:"noshare"`
	Options   [][]string  `json:"options",xml:"options"`
	Hostname  string      `json:"hostname",xml:"hostname"`
}

/*
Describe Cloud Server options, contains
	
  * Id        (int)         Unique Identifier
  
  * Type      (string)      Cloud Server Type ()
  
  * Hostname  (string)      Logical Server Hostname
  
  * Options   ([][]string)  Cloud Server Options
  
	Refers to : https://docs.docker.com/machine/drivers/
*/
type CloudServer struct {
	Id          int       `json:"uid",xml:"uid"`
	Type      string      `json:"type",xml:"type"`
	Hostname  string      `json:"hostname",xml:"hostname"`
	Options   [][]string  `json:"options",xml:"options"`
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
	ClusterMember  RoleType = iota + 1; // value 4
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
	
  * Id          (int)               Unique Identifier
  
  * ServerId    (int)               Target Server Id
  
  * IsCloud     (bool)              Is A Cloud Server
  
  * Type        (InstallationType)  Installation Type
  
  * Environment (RoleType)          Installation Environment (Rancher)
  
  * Role        (EnvironmentType)   Installation Role
*/
type Installation struct {
	Id          int               `json:"uid",xml:"uid"`
	ServerId    int               `json:"serverid",xml:"serverid"`
	IsCloud     bool              `json:"iscloud",xml:"iscloud"`
	Type        InstallationType  `json:"type",xml:"type"`
	Environment RoleType          `json:"environment",xml:"environment"`
	Role        EnvironmentType   `json:"role",xml:"role"`
}

/*
Describe Network options, contains
	
  * Id            (int)             Unique Identifier
  
  * Name          (string)          Network Name
  
  * Servers       ([]Server)        Server List
  
  * CServers      ([]CloudServer)   Cloud Server List
  
  * Installations ([]Installation)  Server Installations
*/
type Network struct {
	Id            int             `json:"uid",xml:"uid"`
	Name          string          `json:"name",xml:"name"`
	Servers       []Server        `json:"servers",xml:"servers"`
	CServers      []CloudServer   `json:"cloudservers",xml:"cloudservers"`
	Installations []Installation  `json:"installations",xml:"installations"`
}

/*
Describe domain options, contains
	
  * Id            (int)             Unique Identifier
  
  * Name          (string)          Domain Name
  
  * Networks      ([]Networks)      Networks List
*/
type Domain struct {
	Id          int             `json:"uid",xml:"uid"`
	Name        string          `json:"name",xml:"name"`
	Networks    []Network       `json:"networks",xml:"networks"`
}

/*
Describe Server State, contains
	
  * Id           (int)        State Unique Identifier
  
  * Hostname     (string)     Defined Host Name
  
  * IPAddresses  ([]string)   Computed IP Address
  
  * ServerId     (int)        Target Server Id
  
  * IsCloud      (bool)       Is A Cloud Server
  
  * NetworkId    (int)        Target Network Id
  
  * DomainId     (int)        Target Domain Id
  
  * Created      (bool)       Creation State
  
  * Altered      (bool)       Alteration State
  
  * Errors       (bool)       Error State
  
  * LastMessage  (string)     Last Alternation Message
*/
type ServerState struct {
  Id          int        `json:"uid",xml:"uid"`
	Hostname    string     `json:"hostname",xml:"hostname"`
	IPAddresses []string   `json:"ipaddresses",xml:"ipaddresses"`
	ServerId    int        `json:"serverid",xml:"serverid"`
	IsCloud     bool       `json:"iscloud",xml:"iscloud"`
	NetworkId   int        `json:"networkid",xml:"networkid"`
	DomainId    int        `json:"domainid",xml:"domainid"`
	Created     bool       `json:"wascreated",xml:"wascreated"`
	Altered     bool       `json:"wasaltered",xml:"wasaltered"`
	Errors      bool       `json:"haserrors",xml:"haserrors"`
	LastMessage string     `json:"lastmessage",xml:"lastmessage"`
}

/*
Describe Network State, contains
	
  * Id           (int)            State Unique Identifier
  
  * NetworkId    (int)            Target Network Id
  
  * DomainId     (int)            Target Domain Id
  
  * Servers      ([]ServerState)  List Of Server State
  
  * Created      (bool)           Creation State
  
  * Altered      (bool)           Alteration State
  
  * Errors       (bool)           Error State
  
  * LastMessage  (string)         Last Alternation Message
*/
type NetworkState struct {
  Id          int             `json:"uid",xml:"uid"`
	NetworkId   int             `json:"networkid",xml:"networkid"`
	DomainId    int             `json:"domainid",xml:"domainid"`
	Servers     []ServerState   `json:"servers",xml:"servers"`
	Created     bool            `json:"wascreated",xml:"wascreated"`
	Altered     bool            `json:"wasaltered",xml:"wasaltered"`
	Errors      bool            `json:"haserrors",xml:"haserrors"`
	LastMessage string          `json:"lastmessage",xml:"lastmessage"`
}

/*
Describe Domain State, contains
	
  * Id           (int)            State Unique Identifier
  
  * DomainId     (int)            Target Domain Id
  
  * Networks     ([]NetworkState) List Of Network State
  
  * Created      (bool)           Creation State
  
  * Altered      (bool)           Alteration State
  
  * Errors       (bool)           Error State
  
  * LastMessage  (string)         Last Alternation Message
*/
type DomainState struct {
	Id          int             `json:"uid",xml:"uid"`
	DomainId    int             `json:"domainid",xml:"domainid"`
	Networks    []NetworkState  `json:"networks",xml:"networks"`
	Created     bool            `json:"wascreated",xml:"wascreated"`
	Altered     bool            `json:"wasaltered",xml:"wasaltered"`
	Errors      bool            `json:"haserrors",xml:"haserrors"`
	LastMessage string          `json:"lastmessage",xml:"lastmessage"`
}

/*
Describe State, contains
	
  * Domains      ([]DomainState)  List Of Domain State
	
  * Created      (bool)           Creation State
	
  * Altered      (bool)           Alteration State
	
  * Errors       (bool)           Error State
	
  * LastMessage  (string)         Last Alternation Message
*/
type State struct {
	Domains     []DomainState   `json:"domains",xml:"domains"`
	Created     bool            `json:"wascreated",xml:"wascreated"`
	Altered     bool            `json:"wasaltered",xml:"wasaltered"`
	Errors      bool            `json:"haserrors",xml:"haserrors"`
	LastMessage string          `json:"lastmessage",xml:"lastmessage"`
}

/*
Describe Entire Infrastructure, contains

  * Name      	(string)  	 Infrastructure Name

  * Domains      ([]Domain)  List Of Domains

  * State        (State)     Creation State

  * Created      (bool)      Creation State
	
  * Altered      (bool)      Alteration State
	
  * Errors       (bool)      Error State
	
  * LastMessage  (string)    Last Alternation Message
*/
type Infrastructure struct {
	Name     		string  	`json:"name",xml:"name"`
	Domains     []Domain  `json:"domains",xml:"domains"`
	State       State     `json:"state",xml:"state"`
	Created     bool      `json:"wascreated",xml:"wascreated"`
	Altered     bool      `json:"wasaltered",xml:"wasaltered"`
	Errors      bool      `json:"haserrors",xml:"haserrors"`
	LastMessage string    `json:"lastmessage",xml:"lastmessage"`
}
