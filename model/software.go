package model

import (
	"fmt"
	"os"
	"io"
	"net/http"
	"os/user"
	"log"
	"errors"
)

type MachineISO struct {
	Name						string `json:"name",xml:"name"`
	BaseURL					string `json:"baseurl",xml:"baseurl"`
	ISOName					string `json:"isoname",xml:"isoname"`
	FolderName			string `json:"folder",xml:"folder"`
	FinalNamePrefix	string `json:"fileprefix",xml:"fileprefix"`
	FinalNameSuffix	string `json:"filesuffix",xml:"filesuffix"`
}

func Homefolder() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
		return  string(os.PathSeparator) + "temp"
	}
	return usr.HomeDir
}

func VMBaseFolder() string {
	home := Homefolder()
	return home + string(os.PathSeparator) + ".vmkube"
}

type MachineActions interface {
	Download(v string) bool
	Check(v string) bool
	Path(v string) string
}

func (isoTemplate *MachineISO)	Path(version string) string {
	home := VMBaseFolder()
	folder := home + string(os.PathSeparator) + "images" + string(os.PathSeparator) + isoTemplate.FolderName
	fileName := folder + string(os.PathSeparator) + isoTemplate.FinalNamePrefix + version + isoTemplate.FinalNameSuffix
	return fileName
}

func (isoTemplate *MachineISO)	Check(version string) bool {
	home := VMBaseFolder()
	folder := home + string(os.PathSeparator) + "images" + string(os.PathSeparator) + isoTemplate.FolderName
	fileName := folder + string(os.PathSeparator) + isoTemplate.FinalNamePrefix + version + isoTemplate.FinalNameSuffix
	_, error := os.Stat(fileName)
	return ! os.IsNotExist(error)
}

func (isoTemplate *MachineISO)	Download(version string) bool {
	url := isoTemplate.BaseURL + version + isoTemplate.ISOName
	home := VMBaseFolder()
	folder := home + string(os.PathSeparator) + "images" + string(os.PathSeparator) + isoTemplate.FolderName
	os.MkdirAll(folder, 0666)
	fileName := folder + string(os.PathSeparator) + isoTemplate.FinalNamePrefix + version + isoTemplate.FinalNameSuffix
	fmt.Printf("Downloading %s to %s\n", url, fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error while creating %s - %s\n", fileName, err)
		return false
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error while downloading %s - %s\n", url, err)
		return false
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Printf("Error while downloading %s - %s\n", url, err)
		return false
	}
	fmt.Printf("%d bytes downloaded.", n)
	_, error := os.Stat(fileName)
	return ! os.IsNotExist(error)
}

func GetMachineAction(name string) (*MachineISO, error) {
	switch name {
	case "rancheros":
		return  &MachineISO{
			Name: "rancheros",
			BaseURL: "https://github.com/rancher/os/releases/download/v",
			ISOName: "/rancheros.iso",
			FinalNamePrefix: "rancheros-",
			FinalNameSuffix: ".iso",
			FolderName: "rancheros",
		}, nil
	default:
		return &MachineISO{}, errors.New("Unbable to discover machine type : " + name)
	}
}