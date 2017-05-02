package action

import (
	"fmt"
	"vmkube/term"
	"vmkube/utils"
)

func ExecuteRequest(request CmdRequest) bool {
	switch request.Type {
	case NoCommand:
		{
			if request.HelpType != NoCommand {
				PrintCommandHelper(request.SubTypeStr, "")
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
			} else {
				PrintCommandHelper("", "")
			}
			return true
		}
	case StartInfrastructure:
		{
			if !request.CheckInfra() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.StartInfra()
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
	case StopInfrastructure:
		{
			if !request.CheckInfra() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.StopInfra()
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
	case RestartInfrastructure:
		{
			if !request.CheckInfra() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.RestartInfra()
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
	case DestroyInfrastructure:
		{
			if !request.CheckInfra() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.DeleteInfra()
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
	case AlterInfrastructure:
		{
			if !request.CheckInfra() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.AlterInfra()
				if !response.Status {
					utils.PrintlnError(term.Screen.Bold(fmt.Sprintf("Error: %s, clause: %s", error, response.Message)))
					term.Screen.Flush()
				}
				if !utils.NO_COLORS {
					term.Screen.ShowCursor()
				}
				return error == nil
			}
		}
	case BackupInfrastructure:
		{
			if !request.CheckInfra() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.BackupInfra()
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
	case RecoverInfrastructure:
		{
			if !request.CheckInfra() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.RecoverInfra()
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
	case ListInfrastructure:
		{
			if !request.CheckInfra() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.StatusInfra()
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
	case ListInfrastructures:
		{
			if !request.CheckInfra() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.ListInfras()
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
	case ListConfigs:
		{
			if !request.CheckProject() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.ListProjects()
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
	case StatusConfig:
		{
			if !request.CheckProject() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.StatusProject()
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
	case BuildConfig:
		{
			if !request.CheckProject() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.BuildProject()
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
	case ImportConfig:
		{
			if !request.CheckProject() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.ImportProject()
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
	case ExportConfig:
		{
			if !request.CheckProject() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.ExportProject()
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
	case DefineConfig:
		{
			if !request.CheckProject() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.CreateProject()
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
	case DeleteConfig:
		{
			if !request.CheckProject() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.DeleteProject()
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
	case AlterConfig:
		{
			if !request.CheckProject() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				response, error := request.AlterProject()
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
	case InfoConfig:
		{
			if !request.CheckProject() {
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return false
			} else {
				request.InfoProject()
			}
			if !utils.NO_COLORS {
				term.Screen.ShowCursor()
			}
			return true
		}
	default:
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
	}
	return false
}
