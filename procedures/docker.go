package procedures

type DockerMachineState int

const (
	Docker_Machine_State_None DockerMachineState = iota
	Docker_Machine_State_Running
	Docker_Machine_State_Paused
	Docker_Machine_State_Saved
	Docker_Machine_State_Stopped
	Docker_Machine_State_Stopping
	Docker_Machine_State_Starting
	Docker_Machine_State_Error
	Docker_Machine_State_Timeout
)

var docker_machine_states = []string{
	"",
	"Running",
	"Paused",
	"Saved",
	"Stopped",
	"Stopping",
	"Starting",
	"Error",
	"Timeout",
}

var docker_machine_states_enum = []DockerMachineState{
	Docker_Machine_State_None,
	Docker_Machine_State_Running,
	Docker_Machine_State_Paused,
	Docker_Machine_State_Saved,
	Docker_Machine_State_Stopped,
	Docker_Machine_State_Stopping,
	Docker_Machine_State_Starting,
	Docker_Machine_State_Error,
	Docker_Machine_State_Timeout,
}

// Given a State type, returns its string representation
func (s DockerMachineState) String() string {
	if int(s) >= 0 && int(s) < len( docker_machine_states) {
		return  docker_machine_states[s]
	}
	return ""
}


// Given a State type, returns its string representation
func GetStateFromDockerAnswer(state string) DockerMachineState {
	for i := 0; i < len(docker_machine_states_enum); i++ {
		if docker_machine_states_enum[i].String() == state {
			return docker_machine_states_enum[i]
		}
	}
	return Docker_Machine_State_None
}

