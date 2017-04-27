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
	Id   string  `json:"Id" xml:"Id"`
	Name string `json:"Name" xml:"Name"`
	Type    int `json:"Type" xml:"Type"`
	Size    int `json:"Size" xml:"Size"`
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
	Enabled       		bool    `json:"Enabled" xml:"Enabled"`
	Host              string  `json:"Host" xml:"Host"`
	UseAddress        bool    `json:"UseAddress" xml:"UseAddress"`
	DiscoveryToken    string   `json:"DiscoveryToken" xml:"DiscoveryToken"`
	UseExperimental   bool    `json:"UseExperimental" xml:"UseExperimental"`
	IsMaster          bool    `json:"IsMaster" xml:"IsMaster"`
	Image             string  `json:"Image" xml:"Image"`
	JoinOpts        []string  `json:"JoinOpts" xml:"JoinOpts"`
	Strategy          string  `json:"Strategy" xml:"Strategy"`
	TLSSan          []string  `json:"TLSSan" xml:"TLSSan"`
}

func ToInstanceSwarmOpt(opt ProjectSwarmOpt) SwarmOpt {
	return SwarmOpt{
		Enabled: opt.Enabled,
		Host: opt.Host,
		UseAddress: opt.UseAddress,
		DiscoveryToken: opt.DiscoveryToken,
		UseExperimental: opt.UseExperimental,
		IsMaster: opt.IsMaster,
		Image: opt.Image,
		JoinOpts: opt.JoinOpts,
		Strategy: opt.Strategy,
		TLSSan: opt.TLSSan,
	}
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
	Environment       []string  `json:"Environment" xml:"Environment"`
	InsecureRegistry  []string  `json:"InsecureRegistry" xml:"InsecureRegistry"`
	RegistryMirror  	[]string  `json:"RegistryMirror" xml:"RegistryMirror"`
	StorageDriver       string  `json:"StorageDriver" xml:"StorageDriver"`
	InstallURL          string  `json:"InstallURL" xml:"InstallURL"`
	Labels            []string  `json:"Labels" xml:"Labels"`
	Options           []string  `json:"Options" xml:"Options"`
}

func ToInstanceEngineOpt(opt ProjectEngineOpt) EngineOpt {
	return EngineOpt{
		Environment: opt.Environment,
		InsecureRegistry: opt.InsecureRegistry,
		InstallURL: opt.InstallURL,
		Labels: opt.Labels,
		Options: opt.Options,
		RegistryMirror: opt.RegistryMirror,
		StorageDriver: opt.StorageDriver,
	}
}

/*
Describe Server options, contains
	
  * Id          (string)     Unique Identifier
	
  * ServerId    (string)     Project Server Unique Identifier
  
  * Name        (string)     Server Local Name
   
  * Roles       ([]string)   Roles used in the deployment plan
 
  * Driver      (string)     Server Driver (virtualbox,vmware,hyperv) ref: https://docs.docker.com/machine/drivers/
  
  * Memory      (int)        Memory Size MB
  
  * Cpus        (int)        Number of Logical Cores
  
  * Swarm       (SwarmOpt)   Swarm Options
  
  * Engine      (EngineOpt)  Engine Options
  
  * OSType      (string)     Machine OS Type (ref: https://docs.docker.com/machine/drivers/os-base/)
  
  * OSVersion   (string)     Machine OS Version (ref: https://docs.docker.com/machine/drivers/os-base/)
  
  * NoShare     (string)     Do not mount home as shared folder
  
  * Options     ([][]string) Specific vendor option in format key value pairs array without signs (e.g.: --<driver>)
  
  * Hostname    (string)     Network Server Hostname

  * IPAddress   (string)     Network IP Address

  * InspectJSON (string)     Inspection Data JSON

  * Logs        (LogStorage) Log information data
*/
type Instance struct {
	Id          string      `json:"Id" xml:"Id"`
	ServerId    string      `json:"ServerId" xml:"ServerId"`
	Name        string      `json:"Name" xml:"Name"`
	Roles     []string      `json:"Roles" xml:"Roles"`
	Driver      string      `json:"Driver" xml:"Driver"`
	Memory        int       `json:"Memory" xml:"Memory"`
	Cpus          int       `json:"Cpus" xml:"Cpus"`
	Disks       []Disk      `json:"Disks" xml:"Disks"`
	Swarm     SwarmOpt      `json:"Swarm" xml:"Swarm"`
	Engine    EngineOpt     `json:"Engine" xml:"Engine"`
	OSType      string      `json:"OSType" xml:"OSType"`
	OSVersion   string      `json:"OSVersion" xml:"OSVersion"`
	NoShare       bool      `json:"NoShare" xml:"NoShare"`
	Options     [][]string  `json:"Options" xml:"Options"`
	Hostname    string      `json:"Hostname" xml:"Hostname"`
	IPAddress   string      `json:"IPAddress" xml:"IPAddress"`
	InspectJSON string      `json:"IPAddress" xml:"IPAddress"`
	Logs     		LogStorage  `json:"Logs" xml:"Logs"`
}

/*
Describe Cloud Server options, contains
	
  * Id          (string)      Unique Identifier
	
  * ServerId    (string)      Project Cloud Server Unique Identifier
  
  * Name        (string)      Cloud Instance Name
  
  * Driver      (string)      Cloud Server Driver (amazonec2, digitalocean, azure, etc...)
  
  * Hostname    (string)      Cloud Server Hostname
  
  * Roles       ([]string)    Roles used in the deployment plan
  
  * Options     ([][]string)  Cloud Server Options
 
  * IPAddress   (string)      Cloud IP Address

  * InspectJSON (string)      Inspection Data JSON

  * Logs        (LogStorage)  Log information data
 
	Refers to : https://docs.docker.com/machine/drivers/
*/
type CloudInstance struct {
	Id          string      `json:"Id" xml:"Id"`
	ServerId    string      `json:"ServerId" xml:"ServerId"`
	Name        string      `json:"Name" xml:"Name"`
	Driver      string      `json:"Driver" xml:"Driver"`
	Hostname    string      `json:"Hostname" xml:"Hostname"`
	Roles   	  []string    `json:"Roles" xml:"Roles"`
	Options     [][]string  `json:"Options" xml:"Options"`
	IPAddress   string      `json:"IPAddress" xml:"IPAddress"`
	InspectJSON string      `json:"IPAddress" xml:"IPAddress"`
	Logs     		LogStorage  `json:"Logs" xml:"Logs"`
}

/*
Describe Installation Type Enum
	
  * Server  Take Part to installation cluster as Main Server
  
  * Host    Take part to installation cluster as simple host
*/
type InstallationType int

const(
	ServerType InstallationType = iota // value 0
	HostType
)

func ToInstanceInstallation(role InstallationRole) InstallationType {
	switch role {
	case ServerRole:
		return ServerType
	default:
		return HostType
	}
}
func InstanceInstallationToString(role InstallationType) string {
	switch role {
	case ServerType:
		return "Server"
	default:
		return "Host"
	}
}

/*
Describe Role Type Enum
	
  * StandAlone      StandAlone Server Unit
  
  * Master          Master Server in a cluster
  
  * Slave           Dependant Server in a cluster
  
  * ClusterMember  Peer Role in a cluster
*/
type RoleType int

const(
	StandAlone      RoleType = iota //value 0
	Master
	Slave
	ClusterMember
)

func ToInstanceRole(role SystemRole) RoleType {
	switch role {
	case StandAloneRole:
		return StandAlone
	case MasterRole:
		return Master
	case SlaveRole:
		return Slave
	default:
		return ClusterMember
	}
}

func InstanceRoleToString(role RoleType) string {
	switch role {
	case StandAlone:
		return "Stand-Alone"
	case Master:
		return "Master"
	case Slave:
		return "Slave"
	default:
		return "Cluster-Member"
	}
}


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
	Cattle        EnvironmentType = iota // value 0
	Kubernetes
	Mesos
	Swarm
	Custom
)

func ToInstanceEnvironment(env ProjectEnvironment) EnvironmentType {
	switch env {
		case CattleEnv:
			return Cattle
		case KubernetesEnv:
			return Kubernetes
		case MesosEnv:
			return Mesos
		case SwarmEnv:
			return Swarm
		default:
			return Custom
	}
}

func InstanceEnvironmentToString(env EnvironmentType) string {
	switch env {
	case Cattle:
		return "Cattle"
	case Kubernetes:
		return "Kubernetes"
	case Mesos:
		return "Mesos"
	case Swarm:
		return "Swarm"
	default:
		return "Custom"
	}
}

/*
Describe Log Storage, contains

  * ProjectId   (string)    Infrastructure Unique Identifier

  * ProjectId   (string)    Project Unique Identifier
  
  * ElementId   (string)    Infrastructure Element Unique Identifier
  
  * LogLines    ([]string)  Log information data
*/
type LogStorage struct {
	InfraId     string            `json:"Id" xml:"Id" mandatory:"yes" descr:"Infrastructure Unique Identifier" type:"text"`
	ProjectId   string            `json:"ProjectId" xml:"ProjectId" mandatory:"yes" descr:"Project Unique Identifier"`
	ElementId   string            `json:"ElementId" xml:"ElementId" mandatory:"yes" descr:"Infrastructure Element Unique Identifier"`
	LogLines    []string          `json:"LogLines" xml:"LogLines,omitempty" mandatory:"no" descr:"Log information data"`
}


/*
Describe Server Installation options, contains
	
  * Id          	(string)            Unique Identifier
  
  * InstanceId    (string)            Target Instance Id
  
  * IsCloud     	(bool)              Is A Cloud Instance
  
  * Type        	(InstallationType)  Installation Type
  
  * Environment 	(RoleType)          Installation Environment (Rancher)
  
  * Role        	(EnvironmentType)   Installation Role

  * Plan     			(InstallationPlan) Reference to installation plan

  * LastExecution	(time.Time ) 				Last Execution Date

  * Success      	(bool)       				Success State

  * Errors        (bool)       				Error State

  * LastMessage   (string)     			Last Error Message

  * Logs  		    (LogStorage)     			Installation Log File Descriptor
*/
type Installation struct {
	Id          	string            `json:"Id" xml:"Id"`
	InstanceId    string            `json:"InstanceId" xml:"InstanceId"`
	IsCloud     	bool              `json:"IsCloud" xml:"IsCloud"`
	Type        	InstallationType  `json:"Type" xml:"Type"`
	Environment 	EnvironmentType   `json:"Environment" xml:"Environment"`
	Role        	RoleType          `json:"Role" xml:"Role"`
	Plan					InstallationPlan	`json:"Plan" xml:"Plan"`
	LastExecution	time.Time  				`json:"LastExecution" xml:"LastExecution"`
	Success     	bool       				`json:"Success" xml:"Success"`
	Errors      	bool       				`json:"Errors" xml:"Errors"`
	LastMessage 	string     				`json:"LastMessage" xml:"LastMessage"`
	Logs     			LogStorage     		`json:"Logs" xml:"Logs"`
}

/*
Describe Network options, contains
	
  * Id            (string)          Unique Identifier
  
  * Name          (string)          Network Name
  
  * Instances     ([]Instance)      Instances List
  
  * CInstances    ([]CloudInstance) Cloud Instances List
  
  * Installations ([]Installation)  Server Executed Installations
  
  * Options       ([][]string)      Specific Network information (eg. cloud provider info or local info)

*/
type Network struct {
	Id          	string          `json:"Id" xml:"Id"`
	Name          string          `json:"Name" xml:"Name"`
	Instances     []Instance      `json:"Instances" xml:"Instances"`
	CInstances    []CloudInstance `json:"CInstances" xml:"CInstances"`
	Installations []Installation  `json:"Installations" xml:"Installations"`
	Options     [][]string        `json:"Options" xml:"Options"`
}

/*
Describe domain options, contains
	
  * Id            (string)          Unique Identifier
  
  * Name          (string)          Domain Name
  
  * Networks      ([]Networks)      Networks List

  * Options       ([][]string)      Specific Domain information (eg. cloud provider info or local info)
*/
type Domain struct {
	Id          string          `json:"Id" xml:"Id"`
	Name        string          `json:"Name" xml:"Name"`
	Networks    []Network       `json:"Networks" xml:"Networks"`
	Options   [][]string        `json:"Options" xml:"Options"`
}

/*
Describe Server State, contains
	
  * Id           (string)     State Unique Identifier
  
  * Hostname     (string)     Defined Host Name
  
  * IPAddresses  ([]string)   Computed IP Address
  
  * InstanceId   (string)     Target Instance Id
  
  * IsCloud      (bool)       Is A Cloud Server
  
  * NetworkId    (string)     Target Network Id
  
  * DomainId     (string)     Target Domain Id

  * Creation     (time.Time ) Creation Date

  * Modified     (time.Time ) Last Modification Date

  * Created      (bool)       Creation State
  
  * Altered      (bool)       Alteration State
  
  * Errors       (bool)       Error State
  
  * LastMessage  (string)     Last Alternation Message
*/
type InstanceState struct {
	Id          string     `json:"Id" xml:"Id"`
	Hostname    string     `json:"Hostname" xml:"Hostname"`
	IPAddresses []string   `json:"IPAddresses" xml:"IPAddresses"`
	InstanceId  string     `json:"InstanceId" xml:"InstanceId"`
	IsCloud     bool       `json:"IsCloud" xml:"IsCloud"`
	NetworkId   string     `json:"NetworkId" xml:"NetworkId"`
	DomainId    string     `json:"DomainId" xml:"DomainId"`
	Creation    time.Time  `json:"Creation" xml:"Creation"`
	Modified 		time.Time  `json:"Modified" xml:"Modified"`
	Created     bool       `json:"Created" xml:"Created"`
	Altered     bool       `json:"Altered" xml:"Altered"`
	Errors      bool       `json:"Errors" xml:"Errors"`
	LastMessage string     `json:"LastMessage" xml:"LastMessage"`
}

/*
Describe Network State, contains
	
  * Id              (string)         State Unique Identifier
  
  * NetworkId       (string)         Target Network Id
  
  * DomainId        (string)         Target Domain Id
  
  * InstanceStates  ([]ServerState)  List Of Instance States

  * Creation        (time.Time )     Creation Date

  * Modified        (time.Time )     Last Modification Date

  * Created         (bool)           Creation State
  
  * Altered         (bool)           Alteration State
  
  * Errors          (bool)           Error State
  
  * LastMessage     (string)         Last Alternation Message
*/
type NetworkState struct {
	Id              string     			`json:"Id" xml:"Id"`
	NetworkId       string          `json:"NetworkId" xml:"NetworkId"`
	DomainId        string          `json:"DomainId" xml:"DomainId"`
	InstanceStates  []InstanceState `json:"InstanceStates" xml:"InstanceStates"`
	Creation        time.Time       `json:"Creation" xml:"Creation"`
	Modified 		    time.Time       `json:"Modified" xml:"Modified"`
	Created         bool            `json:"Created" xml:"Created"`
	Altered         bool            `json:"Altered" xml:"Altered"`
	Errors          bool            `json:"Errors" xml:"Errors"`
	LastMessage     string          `json:"LastMessage" xml:"LastMessage"`
}

/*
Describe Domain State, contains
	
  * Id            (string)          State Unique Identifier
  
  * DomainId      (string)          Target Domain Id
  
  * NetworkStates ([]NetworkState)  List Of Network States

  * Creation      (time.Time ) 		  Creation Date

  * Modified      (time.Time ) 		  Last Modification Date

  * Created       (bool)            Creation State
  
  * Altered       (bool)            Alteration State
  
  * Errors        (bool)            Error State
  
  * LastMessage  (string)           Last Alternation Message
*/
type DomainState struct {
	Id            string     			`json:"Id" xml:"Id"`
	DomainId      string          `json:"DomainId" xml:"DomainId"`
	NetworkStates []NetworkState  `json:"NetworkStates" xml:"NetworkStates"`
	Creation      time.Time       `json:"Creation" xml:"Creation"`
	Modified 		  time.Time       `json:"Modified" xml:"Modified"`
	Created       bool            `json:"Created" xml:"Created"`
	Altered       bool            `json:"Altered" xml:"Altered"`
	Errors        bool            `json:"Errors" xml:"Errors"`
	LastMessage   string          `json:"LastMessage" xml:"LastMessage"`
}

/*
Describe State, contains

	* Id            (string)          State Unique Identifier

  * DomainStates  ([]DomainState)   List Of Domain States

  * Creation      (time.Time ) 		  Creation Date

  * Modified      (time.Time ) 		  Last Modification Date

  * Created       (bool)           Creation State
	
  * Altered       (bool)           Alteration State
	
  * Errors        (bool)           Error State
	
  * LastMessage   (string)         Last Alternation Message
*/
type State struct {
	Id            string     			`json:"Id" xml:"Id"`
	DomainStates  []DomainState   `json:"DomainStates" xml:"DomainStates"`
	Creation      time.Time       `json:"Creation" xml:"Creation"`
	Modified 		  time.Time       `json:"Modified" xml:"Modified"`
	Created       bool            `json:"Created" xml:"Created"`
	Altered       bool            `json:"Altered" xml:"Altered"`
	Errors        bool            `json:"Errors" xml:"Errors"`
	LastMessage   string          `json:"LastMessage" xml:"LastMessage"`
}

/*
Describe Entire Infrastructure, contains

	* Id          (string)      Infrastructure Unique Identifier

	* ProjectId   (string)      Related Project Unique Identifier

  * Name      	(string)  	 	Infrastructure Name

  * Domains      ([]Domain)  	List Of Domains

  * State        (State)     	Creation State

  * Creation     (time.Time ) Creation Date

  * Modified     (time.Time ) Last Modification Date

  * Created      (bool)      	Creation State
	
  * Altered      (bool)      	Alteration State
	
  * Errors       (bool)      	Error State
	
  * LastMessage  (string)    	Last Alternation Message
*/
type Infrastructure struct {
	Id          string     				`json:"Id" xml:"Id"`
	ProjectId   string     				`json:"ProjectId" xml:"ProjectId"`
	Name     		string  					`json:"Name" xml:"Name"`
	Domains     []Domain  				`json:"Domains" xml:"Domains"`
	State       State     				`json:"State" xml:"State"`
	Creation    time.Time       	`json:"Creation" xml:"Creation"`
	Modified 		time.Time       	`json:"Modified" xml:"Modified"`
	Created     bool      				`json:"Created" xml:"Created"`
	Altered     bool      				`json:"Altered" xml:"Altered"`
	Errors      bool      				`json:"Errors" xml:"Errors"`
	LastMessage string    				`json:"LastMessage" xml:"LastMessage"`
}
