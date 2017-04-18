package model

import (
	"os"
	"strings"
	"bufio"
	"fmt"
)

func existsFile(file string) bool {
	_,err := os.Stat(file)
	return  err == nil
}

func deleteIfExists(file string) error {
	_,err := os.Stat(file)
	if err != nil {
		return os.Remove(file)
	}
	return  err
}

func MakeFolderIfNotExists(folder string) error {
	err := os.MkdirAll(folder, 0666)
	return  err
}

type IONature interface {
	Load(file string) error
	Import(file string, format string) error
	Save(file string) error
	Validate() []error
}

func GetLockFile(id string) string {
	folder := VMBaseFolder() + "/.lock"
	os.MkdirAll(folder, 0666)
	return folder + "/" + strings.Replace(id, "-", "_", len(id)) + ".lock"

}

func readLocks(projectId string) ([]string, error) {
	fileName := GetLockFile(projectId)
	if !existsFile(fileName) {
		_, err := os.Create(fileName)
		if err != nil {
			return nil, err
		}
	}
	file, err := os.Open(fileName)
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

func addLock(projectId string, newline string) error {
	fileName := GetLockFile(projectId)
	deleteIfExists(fileName)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(newline+"\n")
	return err
}

func overwriteLocks(projectId string, lines []string) error {
	fileName := GetLockFile(projectId)
	deleteIfExists(fileName)
	file, err := os.Create(fileName)
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


func WriteLock(projectId string, resourceId string) bool {
	return addLock(projectId, resourceId) == nil
}

func RemoveLock(projectId string, resourceId string) bool {
	lines, err := readLocks(projectId)
	if err != nil {
		return false
	}
	newLines := lines[0:]
	for _,line := range lines {
		if resourceId != line {
			newLines = append(newLines, line)
		}
	}
	return overwriteLocks(projectId, newLines) == nil
}

