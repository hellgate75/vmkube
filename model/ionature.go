package model

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"vmkube/utils"
)

func ExistsFile(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func DeleteIfExists(file string) error {
	_, err := os.Stat(file)
	if err == nil {
		return os.Remove(file)
	}
	return err
}

func MakeFolderIfNotExists(folder string) error {
	if _, err := os.Stat(folder); err != nil {
		err := os.MkdirAll(folder, 0777)
		return err
	}
	return nil
}

type IONature interface {
	Load(file string) error
	Import(file string, format string) error
	PostImport() error
	Save(file string) error
	Validate() []error
}

func GetLockFile(id string) string {
	folder := VMBaseFolder() + string(os.PathSeparator) + ".lock"
	os.MkdirAll(folder, 0777)
	return folder + string(os.PathSeparator) + strings.Replace(id, "-", "_", len(id)) + ".lock"

}

func readLocks(containerId string) ([]string, error) {
	fileName := GetLockFile(containerId)
	if !ExistsFile(fileName) {
		return []string{}, nil
	}
	file, err := os.OpenFile(fileName, os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func addLock(containerId string, newline string) error {
	fileName := GetLockFile(containerId)
	DeleteIfExists(fileName)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(newline + "\n")
	return err
}

func overwriteLocks(containerId string, lines []string) error {
	fileName := GetLockFile(containerId)
	DeleteIfExists(fileName)
	err := utils.CreateNewEmptyFile(fileName)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(fileName, os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func WriteLock(containerId string, resourceId string) bool {
	if strings.TrimSpace(resourceId) != "" && strings.TrimSpace(resourceId) != "" {
		return addLock(containerId, resourceId) == nil
	} else {
		return false
	}
}

func RemoveLock(containerId string, resourceId string) bool {
	if strings.TrimSpace(resourceId) != "" && strings.TrimSpace(resourceId) != "" {
		lines, err := readLocks(containerId)
		if err != nil {
			return false
		}
		newLines := make([]string, 0)
		for _, line := range lines {
			if resourceId != line && strings.TrimSpace(resourceId) != "" && strings.TrimSpace(line) != "" {
				newLines = append(newLines, line)
			}
		}
		if len(newLines) == 0 {
			return os.Remove(GetLockFile(containerId)) == nil
		} else {
			return overwriteLocks(containerId, newLines) == nil
		}
	} else {
		return false
	}
}

func HasLock(containerId string, resourceId string) bool {
	if strings.TrimSpace(resourceId) != "" && strings.TrimSpace(resourceId) != "" {
		lines, err := readLocks(containerId)
		if err != nil {
			return false
		}
		for _, line := range lines {
			if resourceId == line {
				return true
			}
		}
		return false
	} else {
		return false
	}
}
