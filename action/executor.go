package action

import (
	"fmt"
	"github.com/hellgate75/vmkube/common"
	"github.com/hellgate75/vmkube/term"
	"github.com/hellgate75/vmkube/utils"
)

func ExecuteRequest(request common.CmdRequest) bool {
	switch request.Type {
	case common.NoCommand:
		{
			if request.HelpType != common.NoCommand {
				common.PrintCommandHelper(request.SubTypeStr, "", GetArgumentHelpers)
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
			} else {
				common.PrintCommandHelper("", "", GetArgumentHelpers)
			}
			return true
		}
	case common.StartInfrastructure:
		{
			if !request.CheckInfra(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.StartInfra(GetArgumentHelpers)
				if !response.Status {
					term.Screen.Println(term.Screen.Color(fmt.Sprintf("Error: %s, clause: %s", error, response.Message), term.RED))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.StopInfrastructure:
		{
			if !request.CheckInfra(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.StopInfra(GetArgumentHelpers)
				if !response.Status {
					term.Screen.Println(term.Screen.Color(fmt.Sprintf("Error: %s, clause: %s", error, response.Message), term.RED))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.RestartInfrastructure:
		{
			if !request.CheckInfra(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.RestartInfra(GetArgumentHelpers)
				if !response.Status {
					term.Screen.Println(term.Screen.Color(fmt.Sprintf("Error: %s, clause: %s", error, response.Message), term.RED))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.DestroyInfrastructure:
		{
			if !request.CheckInfra(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.DeleteInfra(GetArgumentHelpers)
				if !response.Status {
					term.Screen.Println(term.Screen.Color(fmt.Sprintf("Error: %s, clause: %s", error, response.Message), term.RED))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.AlterInfrastructure:
		{
			if !request.CheckInfra(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.AlterInfra(GetArgumentHelpers)
				if !response.Status {
					utils.PrintlnBoldError(fmt.Sprintf("Error: %s, clause: %s", error, response.Message))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.BackupInfrastructure:
		{
			if !request.CheckInfra(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.BackupInfra(GetArgumentHelpers)
				if !response.Status {
					term.Screen.Println(term.Screen.Color(fmt.Sprintf("Error: %s, clause: %s", error, response.Message), term.RED))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.RecoverInfrastructure:
		{
			if !request.CheckInfra(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.RecoverInfra(GetArgumentHelpers)
				if !response.Status {
					term.Screen.Println(term.Screen.Color(fmt.Sprintf("Error: %s, clause: %s", error, response.Message), term.RED))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.ListInfrastructure:
		{
			if !request.CheckInfra(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.StatusInfra(GetArgumentHelpers)
				if !response.Status {
					term.Screen.Println(term.Screen.Color(fmt.Sprintf("Error: %s, clause: %s", error, response.Message), term.RED))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.ListInfrastructures:
		{
			if !request.CheckInfra(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.ListInfras(GetArgumentHelpers)
				if !response.Status {
					term.Screen.Println(term.Screen.Color(fmt.Sprintf("Error: %s, clause: %s", error, response.Message), term.RED))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.ListConfigs:
		{
			if !request.CheckProject(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.ListProjects(GetArgumentHelpers)
				if !response.Status {
					term.Screen.Println(term.Screen.Color(fmt.Sprintf("Error: %s, clause: %s", error, response.Message), term.RED))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.StatusConfig:
		{
			if !request.CheckProject(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.StatusProject(GetArgumentHelpers)
				if !response.Status {
					term.Screen.Println(term.Screen.Color(fmt.Sprintf("Error: %s, clause: %s", error, response.Message), term.RED))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.BuildConfig:
		{
			if !request.CheckProject(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.BuildProject(GetArgumentHelpers)
				if !response.Status {
					term.Screen.Println(term.Screen.Color(fmt.Sprintf("Error: %s, clause: %s", error, response.Message), term.RED))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.ImportConfig:
		{
			if !request.CheckProject(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.ImportProject(GetArgumentHelpers)
				if !response.Status {
					term.Screen.Println(term.Screen.Color(fmt.Sprintf("Error: %s, clause: %s", error, response.Message), term.RED))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.ExportConfig:
		{
			if !request.CheckProject(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.ExportProject(GetArgumentHelpers)
				if !response.Status {
					term.Screen.Println(term.Screen.Color(fmt.Sprintf("Error: %s, clause: %s", error, response.Message), term.RED))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.DefineConfig:
		{
			if !request.CheckProject(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.CreateProject(GetArgumentHelpers)
				if !response.Status {
					term.Screen.Println(term.Screen.Color(fmt.Sprintf("Error: %s, clause: %s", error, response.Message), term.RED))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.DeleteConfig:
		{
			if !request.CheckProject(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.DeleteProject(GetArgumentHelpers)
				if !response.Status {
					term.Screen.Println(term.Screen.Color(fmt.Sprintf("Error: %s, clause: %s", error, response.Message), term.RED))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.AlterConfig:
		{
			if !request.CheckProject(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				response, error := request.AlterProject(GetArgumentHelpers)
				if !response.Status {
					term.Screen.Println(term.Screen.Color(fmt.Sprintf("Error: %s, clause: %s", error, response.Message), term.RED))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case common.InfoConfig:
		{
			if !request.CheckProject(GetArgumentHelpers) {
				common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
				return false
			} else {
				request.InfoProject(GetArgumentHelpers)
			}
			if !utils.NO_COLORS {
				term.Screen.ShowCursor()
			}
			return true
		}
	default:
		common.PrintCommandHelper(request.TypeStr, request.SubTypeStr, GetArgumentHelpers)
	}
	return false
}
