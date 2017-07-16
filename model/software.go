package model

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"github.com/hellgate75/vmkube/utils"
)

type MachineISO struct {
	Name            string `json:"name" xml:"name"`
	BaseURL         string `json:"baseurl" xml:"baseurl"`
	ISOName         string `json:"isoname" xml:"isoname"`
	FolderName      string `json:"folder" xml:"folder"`
	FinalNamePrefix string `json:"fileprefix" xml:"fileprefix"`
	FinalNameSuffix string `json:"filesuffix" xml:"filesuffix"`
}

func HomeFolder() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
		return os.TempDir()
	}
	return usr.HomeDir
}

func VMBaseFolder() string {
	home := HomeFolder()
	return home + string(os.PathSeparator) + ".vmkube"
}

func GetEmergencyFolder() string {
	return HomeFolder()
}

type MachineActions interface {
	Download(v string) bool
	Check(v string) bool
	Path(v string) string
}

func (isoTemplate *MachineISO) Path(version string) string {
	home := VMBaseFolder()
	folder := home + string(os.PathSeparator) + "images" + string(os.PathSeparator) + isoTemplate.FolderName
	fileName := folder + string(os.PathSeparator) + isoTemplate.FinalNamePrefix + version + isoTemplate.FinalNameSuffix
	return fileName
}

func (isoTemplate *MachineISO) Check(version string) bool {
	home := VMBaseFolder()
	folder := home + string(os.PathSeparator) + "images" + string(os.PathSeparator) + isoTemplate.FolderName
	fileName := folder + string(os.PathSeparator) + isoTemplate.FinalNamePrefix + version + isoTemplate.FinalNameSuffix
	_, error := os.Stat(fileName)
	return !os.IsNotExist(error)
}

func (isoTemplate *MachineISO) Download(version string) bool {
	url := isoTemplate.BaseURL + version + isoTemplate.ISOName
	home := VMBaseFolder()
	folder := home + string(os.PathSeparator) + "images" + string(os.PathSeparator) + isoTemplate.FolderName
	os.MkdirAll(folder, 0777)
	fileName := folder + string(os.PathSeparator) + isoTemplate.FinalNamePrefix + version + isoTemplate.FinalNameSuffix
	fmt.Printf("Downloading %s to %s\n", url, fileName)

	// Check() for Download Prevent File Presence on the Disk
	err := utils.CreateNewEmptyFile(fileName)
	if err != nil {
		fmt.Printf("Error while creating %s - %s\n", fileName, err)
		return false
	}

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error while downloading %s - %s\n", url, err)
		return false
	}
	defer response.Body.Close()

	output, err := os.OpenFile(fileName, os.O_RDWR, 0777)
	if err != nil {
		return false
	}
	defer output.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Printf("Error while downloading %s - %s\n", url, err)
		return false
	}
	fmt.Printf("%d bytes downloaded.\n", n)
	_, error := os.Stat(fileName)
	return !os.IsNotExist(error)
}

func GetMachineAction(name string) (*MachineISO, error) {
	switch name {
	case "rancheros":
		return &MachineISO{
			Name:            "rancheros",
			BaseURL:         "https://github.com/rancher/os/releases/download/v",
			ISOName:         "/rancheros.iso",
			FinalNamePrefix: "rancheros-",
			FinalNameSuffix: ".iso",
			FolderName:      "rancheros",
		}, nil
	default:
		return &MachineISO{}, errors.New("Unbable to discover machine type : " + name)
	}
}
