package tasks

import (
	"strings"
	"os"
	"vmkube/model"
	"io/ioutil"
	"fmt"
	"sync"
)

func ConvertActivityTaskInString(task ActivityTask) string {
	operation := "Create Virtual"
	switch task {
	case DestroyMachine:
		operation = "Destroy Virtual"
		break
	case StartMachine:
		operation = "Start Virtual"
		break
	case StopMachine:
		operation = "Stop Virtual"
		break
	case RestartMachine:
		operation = "Restart Virtual"
		break
	case MachineStatus:
		operation = "Get Status of"
		break
	case MachineEnv:
		operation = "Get Environment for"
		break
	case MachineInspect:
		operation = "Get Descriptor for"
		break
	case MachineIPAddress:
		operation = "Get IP Address for"
		break
	case MachineExtendsDisk:
		operation = "Extends Disk of"
		break
	}
	return operation
}

func ConvertSubActivityTaskInString(task ActivityTask) string {
	operation := "creating"
	switch task {
	case DestroyMachine:
		operation = "destroying"
		break
	case StartMachine:
		operation = "starting"
		break
	case StopMachine:
		operation = "stopping"
		break
	case RestartMachine:
		operation = "restarting"
		break
	case MachineStatus:
		operation = "informaing"
		break
	case MachineEnv:
		operation = "environment"
		break
	case MachineInspect:
		operation = "inspecting"
		break
	case MachineIPAddress:
		operation = "ip address"
		break
	case MachineExtendsDisk:
		operation = "extending disk"
		break
	}
	return operation
}

var mutex sync.RWMutex

func DumpData(file string, data interface{}, overwrite bool) {
	text := ""
	if !overwrite {
		if strings.Index(file, string(os.PathSeparator)) < 0  {
			file = model.HomeFolder() + string(os.PathSeparator) + file
		}
		mutex.RLock()
		bytes, err := ioutil.ReadFile(file)
		mutex.RUnlock()

		if err == nil {
			text = string(bytes)
			text = fmt.Sprintf("%s", text)
		}
	}
	text += fmt.Sprintf("%s\n", data)
	mutex.Lock()
	ioutil.WriteFile(file, []byte(text), 0777)
	mutex.Unlock()
}
