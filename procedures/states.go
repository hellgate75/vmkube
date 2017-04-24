package procedures

type MachineState int

const (
	Machine_State_None MachineState = iota
	Machine_State_Running
	Machine_State_Paused
	Machine_State_Saved
	Machine_State_Stopped
	Machine_State_Stopping
	Machine_State_Starting
	Machine_State_Error
	Machine_State_Timeout
)

var machine_states = []string{
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

var machine_states_enum = []MachineState{
	Machine_State_None,
	Machine_State_Running,
	Machine_State_Paused,
	Machine_State_Saved,
	Machine_State_Stopped,
	Machine_State_Stopping,
	Machine_State_Starting,
	Machine_State_Error,
	Machine_State_Timeout,
}

// Given a State type, returns its string representation
func (s MachineState) String() string {
	if int(s) >= 0 && int(s) < len( machine_states) {
		return  machine_states[s]
	}
	return ""
}


// Given a State type, returns its string representation
func GetStateFromMachineAnswer(state string) MachineState {
	for i := 0; i < len(machine_states_enum); i++ {
		if machine_states_enum[i].String() == state {
			return machine_states_enum[i]
		}
	}
	return Machine_State_None
}

