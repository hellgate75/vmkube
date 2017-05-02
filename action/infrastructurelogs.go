package action

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"vmkube/model"
	"vmkube/utils"
)

const MAX_LINES_IN_LOG = 2000

type InfrastructureLogsInfo struct {
	Format string
	Logs   model.LogStorage
}

func (info *InfrastructureLogsInfo) ReadLogFiles() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".logs"
	err := model.MakeFolderIfNotExists(baseFolder)
	if err != nil {
		return err
	}

	var i int = 0

	fileName := baseFolder + string(os.PathSeparator) + ".project-" + utils.IdToFileFormat(info.Logs.ProjectId) + ".infra-" + utils.IdToFileFormat(info.Logs.InfraId) + ".elem-" + utils.IdToFileFormat(info.Logs.ElementId) + "-" + strconv.Itoa(i) + ".log"
	_, err = os.Stat(fileName)
	for err == nil {
		ifaceLog := IFaceLogStorage{
			InfraId:   info.Logs.InfraId,
			ProjectId: info.Logs.ProjectId,
			ElementId: info.Logs.ProjectId,
		}
		ifaceLog.WaitForLogFileUnlock(i)

		LockLogFile(info.Logs, i)

		bytes, err := ioutil.ReadFile(fileName)

		UnlockLogFile(info.Logs, i)

		if err == nil {
			return err
		}
		lines := strings.Split(string(bytes), "\n")
		info.Logs.LogLines = append(info.Logs.LogLines, lines...)
		i++
		fileName = baseFolder + string(os.PathSeparator) + ".project-" + utils.IdToFileFormat(info.Logs.ProjectId) + ".infra-" + utils.IdToFileFormat(info.Logs.InfraId) + ".elem-" + utils.IdToFileFormat(info.Logs.ElementId) + "-" + strconv.Itoa(i) + ".log"
		_, err = os.Stat(fileName)
	}
	return nil
}

func (info *InfrastructureLogsInfo) SaveLogFile() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".logs"
	err := model.MakeFolderIfNotExists(baseFolder)
	if err != nil {
		return err
	}
	var lineLength int = len(info.Logs.LogLines)

	var split bool = (lineLength > MAX_LINES_IN_LOG)

	if split {
		var i int = 0
		var sliceStart int = i * MAX_LINES_IN_LOG
		var sliceEnd int = sliceStart + MAX_LINES_IN_LOG
		for sliceStart < lineLength {
			data := []byte(strings.Join(info.Logs.LogLines[sliceStart:sliceEnd], "\n"))
			fileName := baseFolder + string(os.PathSeparator) + ".project-" + utils.IdToFileFormat(info.Logs.ProjectId) + ".infra-" + utils.IdToFileFormat(info.Logs.InfraId) + ".elem-" + utils.IdToFileFormat(info.Logs.ElementId) + "-" + strconv.Itoa(i) + ".log"
			_, err = os.Stat(fileName)
			if err != nil || lineLength-sliceStart <= MAX_LINES_IN_LOG {

				ifaceLog := IFaceLogStorage{
					InfraId:   info.Logs.InfraId,
					ProjectId: info.Logs.ProjectId,
					ElementId: info.Logs.ProjectId,
				}
				ifaceLog.WaitForLogFileUnlock(i)

				LockLogFile(info.Logs, i)

				if err == nil {
					model.DeleteIfExists(fileName)
				}
				err = ioutil.WriteFile(fileName, data, 0777)

				UnlockLogFile(info.Logs, i)

				if err != nil {
					return err
				}
			}
			i++
			sliceStart = i * MAX_LINES_IN_LOG
			sliceEnd = sliceStart + MAX_LINES_IN_LOG
		}

	} else {
		fileName := baseFolder + string(os.PathSeparator) + ".project-" + utils.IdToFileFormat(info.Logs.ProjectId) + ".infra-" + utils.IdToFileFormat(info.Logs.InfraId) + ".elem-" + utils.IdToFileFormat(info.Logs.ElementId) + "-0.log"
		data := []byte(strings.Join(info.Logs.LogLines, "\n"))
		err = ioutil.WriteFile(fileName, data, 0777)
		return err
	}
	return nil
}

func (info *InfrastructureLogsInfo) DeleteLogFile() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".logs"
	err := model.MakeFolderIfNotExists(baseFolder)
	if err != nil {
		return err
	}
	var i int = 0
	var errorCount int = 0

	fileName := baseFolder + string(os.PathSeparator) + ".project-" + utils.IdToFileFormat(info.Logs.ProjectId) + ".infra-" + utils.IdToFileFormat(info.Logs.InfraId) + ".elem-" + utils.IdToFileFormat(info.Logs.ElementId) + "-" + strconv.Itoa(i) + ".log"
	_, err = os.Stat(fileName)
	for err == nil || errorCount < 5 {
		if err == nil {
			ifaceLog := IFaceLogStorage{
				InfraId:   info.Logs.InfraId,
				ProjectId: info.Logs.ProjectId,
				ElementId: info.Logs.ProjectId,
			}
			ifaceLog.WaitForLogFileUnlock(i)

			LockLogFile(info.Logs, i)

			model.DeleteIfExists(fileName)

			UnlockLogFile(info.Logs, i)

		} else {
			errorCount++
		}
		i++
		fileName = baseFolder + string(os.PathSeparator) + ".project-" + utils.IdToFileFormat(info.Logs.ProjectId) + ".infra-" + utils.IdToFileFormat(info.Logs.InfraId) + ".elem-" + utils.IdToFileFormat(info.Logs.ElementId) + "-" + strconv.Itoa(i) + ".log"
		_, err = os.Stat(fileName)
	}
	return nil
}

func (info *InfrastructureLogsInfo) Read() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".data"
	err := model.MakeFolderIfNotExists(baseFolder)
	if err != nil {
		return err
	}
	fileName := baseFolder + string(os.PathSeparator) + "." + utils.IdToFileFormat(info.Logs.ProjectId) + ".infra-" + utils.IdToFileFormat(info.Logs.InfraId) + ".elem-" + utils.IdToFileFormat(info.Logs.ElementId) + ".actionindex"
	if _, err = os.Stat(fileName); err != nil {
		info.Logs = model.LogStorage{
			InfraId:   info.Logs.InfraId,
			ProjectId: info.Logs.ProjectId,
			ElementId: info.Logs.ElementId,
			LogLines:  []string{},
		}
		return nil
	}

	ifaceLog := IFaceLogStorage{
		InfraId:   info.Logs.InfraId,
		ProjectId: info.Logs.ProjectId,
		ElementId: info.Logs.ProjectId,
	}
	ifaceLog.WaitForUnlock()

	LockLog(info.Logs)

	err = info.Logs.Load(fileName)

	UnlockLog(info.Logs)

	return err
}

func (info *InfrastructureLogsInfo) Write() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + "." + utils.IdToFileFormat(info.Logs.ProjectId) + ".infra-" + utils.IdToFileFormat(info.Logs.InfraId) + ".elem-" + utils.IdToFileFormat(info.Logs.ElementId) + ".vmkubelog"

	ifaceLog := IFaceLogStorage{
		InfraId:   info.Logs.InfraId,
		ProjectId: info.Logs.ProjectId,
		ElementId: info.Logs.ProjectId,
	}
	ifaceLog.WaitForUnlock()

	LockLog(info.Logs)

	emptyLog := model.LogStorage{
		InfraId:   info.Logs.InfraId,
		ProjectId: info.Logs.ProjectId,
		ElementId: info.Logs.ElementId,
		LogLines:  []string{},
	}

	if _, err := os.Stat(fileName); err == nil {
		model.DeleteIfExists(fileName)
	}

	err := emptyLog.Save(fileName)

	UnlockLog(info.Logs)

	return err
}

func (info *InfrastructureLogsInfo) Import(file string, format string) error {
	return info.Logs.Import(file, format)
}

func (info *InfrastructureLogsInfo) Delete() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + "." + utils.IdToFileFormat(info.Logs.ProjectId) + ".infra-" + utils.IdToFileFormat(info.Logs.InfraId) + ".elem-" + utils.IdToFileFormat(info.Logs.ElementId) + ".vmkubelog"
	_, err := os.Stat(fileName)
	if err == nil {
		ifaceLog := IFaceLogStorage{
			InfraId:   info.Logs.InfraId,
			ProjectId: info.Logs.ProjectId,
			ElementId: info.Logs.ProjectId,
		}
		ifaceLog.WaitForUnlock()
		LockLog(info.Logs)
		err = model.DeleteIfExists(fileName)
		UnlockLog(info.Logs)
		if err != nil {
			return err
		}
		return info.DeleteLogFile()
	}
	return nil
}

func (info *InfrastructureLogsInfo) Export(prettify bool) ([]byte, error) {
	emptyLog := model.LogStorage{
		InfraId:   info.Logs.InfraId,
		ProjectId: info.Logs.ProjectId,
		ElementId: info.Logs.ElementId,
		LogLines:  []string{},
	}
	if "json" == info.Format {
		return utils.GetJSONFromElem(emptyLog, prettify)
	} else if "xml" == info.Format {
		return utils.GetXMLFromElem(emptyLog, prettify)
	} else {
		return []byte{}, errors.New("Format type : " + info.Format + " not provided ...")
	}
}
