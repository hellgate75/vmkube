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
	Enabled         bool     `json:"Enabled" xml:"Enabled" mandatory:"yes" descr:"Enable Swarm Features" type:"boolean"`
	Host            string   `json:"Host" xml:"Host" mandatory:"no" descr:"Swarm discovery Address (tcp://0.0.0.0:3376)" type:"text"`
	UseAddress      bool     `json:"UseAddress" xml:"UseAddress" mandatory:"no" descr:"Use Machine IP Address" type:"boolean"`
	DiscoveryToken  string   `json:"DiscoveryToken" xml:"DiscoveryToken" mandatory:"no" descr:"Use Swarm Discovery Option (token://<token>)" type:"text"`
	UseExperimental bool     `json:"UseExperimental" xml:"UseExperimental" mandatory:"no" descr:"Use Docker Experimental Features" type:"boolean"`
	IsMaster        bool     `json:"IsMaster" xml:"IsMaster" mandatory:"no" descr:"Is Swarm Master" type:"boolean"`
	Image           string   `json:"Image" xml:"Image" mandatory:"no" descr:"Swarm Image (e.g.: smarm:latest)" type:"text"`
	JoinOpts        []string `json:"JoinOpts" xml:"JoinOpts" mandatory:"no" descr:"Swarm Join Options" type:"text list"`
	Strategy        string   `json:"Strategy" xml:"Strategy" mandatory:"no" descr:"Swarm Strategy" type:"text"`
	TLSSan          []string `json:"TLSSan" xml:"TLSSan" mandatory:"no" descr:"Swarm TLS Specific Options (No overwrite, be sure of syntax, eg: [\"--my-tls-key my-tls-value\"])" type:"text list"`
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
	Environment      []string `json:"Environment" xml:"Environment" mandatory:"no" descr:"Environment variables" type:"text list"`
	InsecureRegistry []string `json:"InsecureRegistry" xml:"InsecureRegistry" mandatory:"no" descr:"Insecure Registry Options" type:"text list"`
	RegistryMirror   []string `json:"RegistryMirror" xml:"RegistryMirror" mandatory:"no" descr:"Registry Mirror Options" type:"text list"`
	StorageDriver    string   `json:"StorageDriver" xml:"StorageDriver" mandatory:"no" descr:"Storage Driver" type:"text"`
	InstallURL       string   `json:"InstallURL" xml:"InstallURL" mandatory:"no" descr:"Docker Install URL (e.g.: https://get.docker.com)" type:"text"`
	Labels           []string `json:"Labels" xml:"Labels" mandatory:"no" descr:"Engine Labels" type:"text list"`
	Options          []string `json:"Options" xml:"Options" mandatory:"no" descr:"Engine Options" type:"text list"`
}

/*
Describe Machine options, contains

  * Id        (string)     Unique Identifier

  * Name      (string)     Machine Local Name

  * Roles     ([]string)   Roles used in the deployment plan

  * Driver    (string)     Machine Driver (virtualbox,vmware,hyperv) ref: https://docs.docker.com/machine/drivers/

  * Memory    (int)        Memory Size MB

  * Cpus      (int)        Number of Logical Cores

  * DiskSize  (int)        Dimension of disk root (in GB)

  * Swarm     (SwarmOpt)   Swarm Options

  * Engine    (EngineOpt)  Engine Options

  * OSType    (string)     Machine OS Type (ref: https://docs.docker.com/machine/drivers/os-base/)

  * OSVersion (string)     Machine OS Version (ref: https://docs.docker.com/machine/drivers/os-base/)

  * NoShare   (string)     Mount or not home as shared folder

  * Options   ([][]string) Specific vendor option in format key value pairs array without driver (i.e.: --<driver>-<option>), no value options are accepted

  * Hostname  (string)     Logical Machine Hostname
*/
type LocalMachine struct {
	Id        string           `json:"Id" xml:"Id" mandatory:"yes" descr:"Unique Identifier" type:"text"`
	Name      string           `json:"Name" xml:"Name" mandatory:"yes" descr:"Machine Local Name" type:"text"`
	Roles     []string         `json:"Roles" xml:"Roles" mandatory:"no" descr:"Roles used in the deployment plan" type:"text list"`
	Driver    string           `json:"Driver" xml:"Driver" mandatory:"yes" descr:"Machine Driver (virtualbox,vmware,hyperv) ref: https://docs.docker.com/machine/drivers/" type:"text"`
	Memory    int              `json:"Memory" xml:"Memory" mandatory:"no" descr:"Memory Size MB" type:"integer"`
	Cpus      int              `json:"Cpus" xml:"Cpus" mandatory:"no" descr:"Number of Logical Cores" type:"integer"`
	DiskSize  int              `json:"DiskSize" xml:"DiskSize" mandatory:"no" descr:"Dimension of disk root (in GB)" type:"integer"`
	Swarm     ProjectSwarmOpt  `json:"Swarm" xml:"Swarm" mandatory:"no" descr:"Swarm Options" type:"object Swarm"`
	Engine    ProjectEngineOpt `json:"Engine" xml:"Engine" mandatory:"no" descr:"Engine Options" type:"object Engine"`
	OSType    string           `json:"OSType" xml:"OSType" mandatory:"yes" descr:"Machine OS Type (ref: https://docs.docker.com/machine/drivers/os-base/)" type:"text"`
	OSVersion string           `json:"OSVersion" xml:"OSVersion" mandatory:"yes" descr:"Machine OS Version (ref: https://docs.docker.com/machine/drivers/os-base/)" type:"text"`
	NoShare   bool             `json:"NoShare" xml:"NoShare" mandatory:"no" descr:"Mount or not home as shared folder" type:"boolean"`
	Options   [][]string       `json:"Options" xml:"Options" mandatory:"no" descr:"Specific vendor option in format key value pairs array without driver (i.e.: --<driver>-<option>), no value options are accepted" type:"text list of couple text list"`
	Hostname  string           `json:"Hostname" xml:"Hostname" mandatory:"yes" descr:"Logical Machine Hostname" type:"text"`
}

/*
Describe Cloud Machine options, contains

  * Id        (string)      Unique Identifier

  * Name      (string)      Cloud Instance Name

  * Driver    (string)      Cloud Machine Driver (amazonec2, digitalocean, azure, etc...)

  * Hostname  (string)      Logical Machine Hostname

  * Roles     ([]string)    Roles used in the deployment plan

  * Options   ([][]string)  Cloud Machine Options (vendor specific options) without  driver (i.e.: --<driver>-<option>), no value options are accepted

	Refers to : https://docs.docker.com/machine/drivers/
*/
type CloudMachine struct {
	Id       string     `json:"Id" xml:"Id" mandatory:"yes" descr:"Unique Identifier" type:"text"`
	Name     string     `json:"Name" xml:"Name" mandatory:"yes" descr:"Cloud Instance Name" type:"text"`
	Driver   string     `json:"Driver" xml:"Driver" mandatory:"yes" descr:"Cloud Machine Driver (amazonec2, digitalocean, azure, etc...)" type:"text"`
	Hostname string     `json:"Hostname" xml:"Hostname" mandatory:"yes" descr:"Logical Machine Hostname" type:"text"`
	Roles    []string   `json:"Roles" xml:"Roles" mandatory:"no" descr:"Roles used in the deployment plan" type:"text list"`
	Options  [][]string `json:"Options" xml:"Options" mandatory:"yes" descr:"Cloud Machine Options (vendor specific options) without  driver (i.e.: --<driver>-<option>), no value options are accepted" type:"text list of couple text list"`
}

/*
Describe Installation Role Enum

  * Machine  Take Part to installation cluster as Main Machine

  * Host    Take part to installation cluster as simple host
*/
type DeploymentRole string

const (
	MachineDeployment DeploymentRole = "Machine"
	HostDeployment    DeploymentRole = "Host"
)

/*
Describe System Role Enum

  * StandAlone      StandAlone Machine Unit

  * Master          Master Machine in a cluster

  * Slave           Dependant Machine in a cluster

  * ClusterMember  Peer Role in a cluster
*/
type MachineRole string

const (
	StandAloneRole    MachineRole = "Stand-Alone"
	MasterRole        MachineRole = "Master"
	SlaveRole         MachineRole = "Slave"
	ClusterMemberRole MachineRole = "Cluster-Memeber"
)

/*
Describe Environment Project Environment Enum (Rancher OS)

  * Cattle         Cattle Environment

  * Kubernetes     Kubernetes Environment

  * Mesos          Mesos Environment

  * Swarm          Swarm Environment

  * Custom         Custom Environment
*/
type MachineEnvironment string

const (
	CattleEnv     MachineEnvironment = "Cattle"
	KubernetesEnv MachineEnvironment = "Kubernetes"
	MesosEnv      MachineEnvironment = "Mesos"
	SwarmEnv      MachineEnvironment = "Swarm"
	CustomEnv     MachineEnvironment = "Custom"
)

/*
Describe Command Set Type Enum (Rancher OS)

  * VirtKube       Virtual Kube Command Set

  * Ansible        Ansible Command Set

  * Helm	        	Helm Command Set
*/
type CommandSet string

const (
	VirtKubeCmdSet CommandSet = "VirtKube"
	AnsibleCmdSet  CommandSet = "Ansible"
	HelmCmdSet     CommandSet = "Helm"
)

/*
Describe Machine Installation Plan, contains

  * Id          	(string)              Unique Identifier

  * MachineId    	(string)              Target Machine Id

  * IsCloud     	(bool)              	Is A Cloud Machine

  * Type        	(InstallationRole)  	Installation Type (Machine,Host)

  * Environment (ProjectEnvironment)    Installation Environment (Cattle,Kubernetes,Mesos,Swarm,Custom)

  * Role        				(SystemRole)	  Installation Role (Stand-Alone,Master,Slave,Cluster-Memeber)

  * MainCommandSet      (CommandSet)	  Command Set used for installation (VirtKube,Ansible,Helm)

  * MainCommandRef      (string)			  Location of installation commands (http,file,git protocols are accepted)

  * ProvisionCommandSet (CommandSet)	  Command Set used for provisioning (VirtKube,Ansible,Helm)

  * ProvisionCommandRef (string)			  Location of provision commands (http,file,git protocols are accepted)
*/
type InstallationPlan struct {
	Id                  string             `json:"Id" xml:"Id" mandatory:"yes" descr:"Unique Identifier" type:"text"`
	MachineId           string             `json:"MachineId" xml:"MachineId" mandatory:"yes" descr:"Target Machine Id" type:"text"`
	IsCloud             bool               `json:"IsCloud" xml:"IsCloud" mandatory:"yes" descr:"Is A Cloud Machine" type:"boolean"`
	Type                DeploymentRole     `json:"Type" xml:"Type" mandatory:"yes" descr:"Installation Type (Machine,Host)" type:"text"`
	Role                MachineRole        `json:"Environment" xml:"Environment" mandatory:"yes" descr:"Installation Role (Stand-Alone,Master,Slave,Cluster-Memeber)" type:"text"`
	Environment         MachineEnvironment `json:"Role" xml:"Role" mandatory:"yes" descr:"Installation Environment (Cattle,Kubernetes,Mesos,Swarm,Custom)" type:"text"`
	MainCommandSet      CommandSet         `json:"MainCommandSet" xml:"MainCommandSet" mandatory:"no" descr:"Command Set used for installation (VirtKube,Ansible,Helm)" type:"text"`
	MainCommandRef      string             `json:"MainCommandRef" xml:"MainCommandRef" mandatory:"no" descr:"Location of installation commands (http,file,git protocols are accepted)" type:"text"`
	ProvisionCommandSet CommandSet         `json:"ProvisionCommandSet" xml:"ProvisionCommandSet" mandatory:"no" descr:"Command Set used for provisioning (VirtKube,Ansible,Helm)" type:"text"`
	ProvisionCommandRef string             `json:"ProvisionCommandRef" xml:"ProvisionCommandRef" mandatory:"no" descr:"Location of provision commands (http,file,git protocols are accepted)" type:"text"`
}

/*
Describe Network options, contains

  * Id            (string)             		Unique Identifier

  * Name          (string)          			Network Name

  * Machines       ([]ProjectMachine)       Machine List

  * CMachines      ([]ProjectCloudMachine)  Cloud Machine List

  * Installations ([]InstallationPlan)  	Machine Installation Plans

  * Options       ([][]string)            Specific Network information (eg. cloud provider info or local info)
*/
type MachineNetwork struct {
	Id            string             `json:"Id" xml:"Id" mandatory:"yes" descr:"Unique Identifier" type:"text"`
	Name          string             `json:"Name" xml:"Name" mandatory:"yes" descr:"Network Name" type:"text"`
	LocalMachines []LocalMachine     `json:"LocalMachines" xml:"LocalMachines" mandatory:"yes" descr:"Machine List" type:"object Machines list"`
	CloudMachines []CloudMachine     `json:"CloudMachines" xml:"CloudMachines" mandatory:"yes" descr:"Cloud Machine List" type:"object CMachines list"`
	Installations []InstallationPlan `json:"Installations" xml:"Installations" mandatory:"no" descr:"Machine Installation Plans" type:"object Installations list"`
	Options       [][]string         `json:"Options" xml:"Options" mandatory:"no" descr:"Specific Network information (eg. cloud provider info or local info)" type:"text list of couple text list"`
}

/*
Describe domain options, contains

  * Id            (string)            Unique Identifier

  * Name          (string)          	Domain Name

  * Networks      ([]ProjectNetwork)	Networks List

  * Options       ([][]string)        Specific Domain information (eg. cloud provider info or local info)
*/
type MachineDomain struct {
	Id       string           `json:"Id" xml:"Id" mandatory:"yes" descr:"Unique Identifier" type:"text"`
	Name     string           `json:"Name" xml:"Name" mandatory:"yes" descr:"Domain Name" type:"text"`
	Networks []MachineNetwork `json:"Networks" xml:"Networks" mandatory:"yes" descr:"Networks List" type:"ovject Networks list"`
	Options  [][]string       `json:"Options" xml:"Options" mandatory:"no" descr:"Specific Domain information (eg. cloud provider info or local info)" type:"text list of couple text list"`
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
	Id          string          `json:"Id" xml:"Id" mandatory:"yes" descr:"Unique Identifier" type:"text"`
	Name        string          `json:"Name" xml:"Name" mandatory:"yes" descr:"Project Name" type:"text"`
	Open        bool            `json:"Open" xml:"Open" mandatory:"no" descr:"Writable State" type:"boolean"`
	Domains     []MachineDomain `json:"Domains" xml:"Domains" mandatory:"yes" descr:"List Of Domains" type:"object Domains list"`
	Created     time.Time       `json:"Created" xml:"Created" mandatory:"no" descr:"Creation Date" type:"datetime"`
	Modified    time.Time       `json:"Modified" xml:"Modified" mandatory:"no" descr:"Last Modification Date" type:"datetime"`
	Errors      bool            `json:"Errors" xml:"Errors" mandatory:"no" descr:"Error State" type:"boolean"`
	LastMessage string          `json:"LastMessage" xml:"LastMessage" mandatory:"no" descr:"Last Alternation Message" type:"text"`
}

/*
Describe Project Import Model, contains

  * Id          (string)            Unique Identifier

  * Name      	(string)  	 				Project Name

  * Domains      ([]ProjectDomain)  List Of Domains
*/
type ProjectImport struct {
	Id      string          `json:"Id" xml:"Id" mandatory:"yes" descr:"Unique Identifier" type:"text"`
	Name    string          `json:"Name" xml:"Name" mandatory:"yes" descr:"Project Name" type:"text"`
	Domains []MachineDomain `json:"Domains" xml:"Domains" mandatory:"yes" descr:"List Of Domains" type:"object Domains list"`
}

/*
Describe Project State in Index, contains

  * Id          (string)            Project Descriptor Unique Identifier

  * Name      	(string)  	 				Project Name

  * InfraId     (string)            Infrastructure Unique Identifier

  * InfraName   (string)  	 				Infrastructure Project Name

	* Open       	(bool)      				Project Writable State

  * Active      (bool)						  Active State of Infrastructure

  * Synced      (bool)						  Sync State of Project
*/
type ProjectsDescriptor struct {
	Id        string `json:"Id" xml:"Id" mandatory:"yes" descr:"Project Unique Identifier" type:"text"`
	Name      string `json:"Name" xml:"Name" mandatory:"yes" descr:"Project Name" type:"text"`
	InfraId   string `json:"InfraId" xml:"InfraId" mandatory:"yes" descr:"Infrastructure Unique Identifier" type:"text"`
	InfraName string `json:"InfraName" xml:"InfraName" mandatory:"yes" descr:"Infrastructure Project Name" type:"text"`
	Active    bool   `json:"Active" xml:"Active" mandatory:"yes" descr:"Active State of Infrastructure" type:"boolean"`
	Open      bool   `json:"Open" xml:"Open" mandatory:"yes" descr:"Project Writable State" type:"boolean"`
	Synced    bool   `json:"Synced" xml:"Synced" mandatory:"no" descr:"Sync State of Project" type:"boolean"`
}

/*
Describe Projects Index, contains

  * Id          (string)                  Indexes Unique Identifier

  * Projects    ([]ProjectsDescriptor)    Projects indexed in VMKube
*/
type ProjectsIndex struct {
	Id       string               `json:"Id" xml:"Id" mandatory:"yes" descr:"Project Unique Identifier" type:"text"`
	Projects []ProjectsDescriptor `json:"Projects" xml:"Projects" mandatory:"yes" descr:"Projects indexed in VMKube" type:"object Projects list"`
}
