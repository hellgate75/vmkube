package action

import (
	"errors"
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

	logFileName := baseFolder + string(os.PathSeparator) + ".project-" + utils.IdToFileFormat(info.Logs.ProjectId) + ".infra-" + utils.IdToFileFormat(info.Logs.InfraId) + ".elem-" + utils.IdToFileFormat(info.Logs.ElementId) + ".log"
	_, err = os.Stat(logFileName)
	if err == nil {
		ifaceLog := IFaceLogStorage{
			InfraId:   info.Logs.InfraId,
			ProjectId: info.Logs.ProjectId,
			ElementId: info.Logs.ProjectId,
		}
		ifaceLog.WaitForLogFileUnlock(0)

		LockLogFile(info.Logs, 0)

		logs, err := zipReadMultiPart(logFileName)
		
		if err != nil {
			return err
		}

		UnlockLogFile(info.Logs, 0)
		
		
		for _,log := range logs {
			bytes := log.Body
			lines := strings.Split(string(bytes), "\n")
			info.Logs.LogLines = append(info.Logs.LogLines, lines...)
		}
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
	
	logFileName := baseFolder + string(os.PathSeparator) + ".project-" + utils.IdToFileFormat(info.Logs.ProjectId) + ".infra-" + utils.IdToFileFormat(info.Logs.InfraId) + ".elem-" + utils.IdToFileFormat(info.Logs.ElementId) + ".log"
	
	var logs []CompressorData = make([]CompressorData, 0)
	if split {
		var i int = 0
		var sliceStart int = i * MAX_LINES_IN_LOG
		var sliceEnd int = sliceStart + MAX_LINES_IN_LOG
		for sliceStart < lineLength {
			if lineLength-sliceStart <= MAX_LINES_IN_LOG {
				data := []byte{}
				lines := info.Logs.LogLines[sliceStart:sliceEnd]
				for _,line := range lines {
					// Prevent empty lines
					if strings.TrimSpace(line) != "" {
						line += "\n"
						data = append(data, []byte(line)...)
					}
				}
				logs = append(logs, CompressorData{
					Descriptor: strconv.Itoa(i),
					Body: data,
				})
			}
			i++
			sliceStart = i * MAX_LINES_IN_LOG
			sliceEnd = sliceStart + MAX_LINES_IN_LOG
		}

	} else {
		data := []byte{}
		lines := info.Logs.LogLines
		for _,line := range lines {
			// Prevent empty lines
			if strings.TrimSpace(line) != "" {
				line += "\n"
				data = append(data, []byte(line)...)
			}
		}
		logs = append(logs, CompressorData{
			Descriptor: "0",
			Body: data,
		})
	}
	ifaceLog := IFaceLogStorage{
		InfraId:   info.Logs.InfraId,
		ProjectId: info.Logs.ProjectId,
		ElementId: info.Logs.ProjectId,
	}
	ifaceLog.WaitForLogFileUnlock(0)
	LockLogFile(info.Logs, 0)
	err = zipWriteMultiPart(logFileName, logs)
	if err != nil {
		return err
	}
	UnlockLogFile(info.Logs, 0)
	return err
}

func (info *InfrastructureLogsInfo) DeleteLogFile() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".logs"
	err := model.MakeFolderIfNotExists(baseFolder)
	if err != nil {
		return err
	}
	logFileName := baseFolder + string(os.PathSeparator) + ".project-" + utils.IdToFileFormat(info.Logs.ProjectId) + ".infra-" + utils.IdToFileFormat(info.Logs.InfraId) + ".elem-" + utils.IdToFileFormat(info.Logs.ElementId) + ".log"
	_, err = os.Stat(logFileName)
	if err == nil {
		ifaceLog := IFaceLogStorage{
			InfraId:   info.Logs.InfraId,
			ProjectId: info.Logs.ProjectId,
			ElementId: info.Logs.ProjectId,
		}
		ifaceLog.WaitForLogFileUnlock(0)
		if err != nil {
			return err
		}
		LockLogFile(info.Logs, 0)
		
		model.DeleteIfExists(logFileName)
		
		UnlockLogFile(info.Logs, 0)
		
	}
	return nil
}

func (info *InfrastructureLogsInfo) Exists() bool {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".data"
	err := model.MakeFolderIfNotExists(baseFolder)
	if err != nil {
		return false
	}
	fileName := baseFolder + string(os.PathSeparator) + "." + utils.IdToFileFormat(info.Logs.ProjectId) + ".infra-" + utils.IdToFileFormat(info.Logs.InfraId) + ".elem-" + utils.IdToFileFormat(info.Logs.ElementId) + ".infralogs"
	if _, err = os.Stat(fileName); err != nil {
		return  false
	}
	return  true
}

func (info *InfrastructureLogsInfo) Read() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".data"
	err := model.MakeFolderIfNotExists(baseFolder)
	if err != nil {
		return err
	}
	fileName := baseFolder + string(os.PathSeparator) + "." + utils.IdToFileFormat(info.Logs.ProjectId) + ".infra-" + utils.IdToFileFormat(info.Logs.InfraId) + ".elem-" + utils.IdToFileFormat(info.Logs.ElementId) + ".infralogs"
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
