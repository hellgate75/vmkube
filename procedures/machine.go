package procedures

import (
	"fmt"
	"vmkube/model"
	"log"
	"encoding/json"
	"io/ioutil"
	"encoding/xml"
)

func DownloadISO(machineType string, version string) bool {
	machineAction, error := model.GetMachineAction(machineType)
	if error == nil {
		if ! machineAction.Check(version) {
			fmt.Println("Machine",machineType,"Version",version,"not present, downloading from internet...")
			downloaded := machineAction.Download(version)
			fmt.Println("Machine",machineType,"Version",version,"dowanloaded:",downloaded)
			return downloaded
		} else {
			fmt.Println("Machine",machineType,"Version",version,"already dowanloaded...")
			return true
		}
	} else {
		log.Fatal("Machine",machineType,"not found!!","-","error:", error)
		return  false
	}
}

func ImportJSONInfrastructure(file string) (model.Infrastructure, error) {
	bytes, error := ioutil.ReadFile(file)
	infrastructure := model.Infrastructure{}
	if error == nil  {
		error = json.Unmarshal(bytes, infrastructure)
		if error == nil  {
			return  infrastructure, nil
		} else {
			//error unmarhalling json file
			return  infrastructure, error
		}
	} else  {
		//error reading file
		return  infrastructure, error
	}
}

func ImportXMLInfrastructure(file string) (model.Infrastructure, error) {
	bytes, error := ioutil.ReadFile(file)
	infrastructure := model.Infrastructure{}
	if error == nil  {
		error = xml.Unmarshal(bytes, infrastructure)
		if error == nil  {
			return  infrastructure, nil
		} else {
			//error unmarhalling xml file
			return  infrastructure, error
		}
	} else  {
		//error reading file
		return  infrastructure, error
	}
}