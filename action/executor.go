package action

import "fmt"

func ExecuteRequest(request CmdRequest) bool{
	//println(request.TypeStr,request.Type)
	//println(request.SubTypeStr,request.HelpType)
	switch request.Type {
		case NoCommand: {
			if request.HelpType != NoCommand  {
				PrintCommandHelper(request.SubTypeStr, "")
			} else {
				PrintCommandHelper("", "")
			}
			return  true
		}
	case StartInfrastructure: {
		if ! request.CheckInfra() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			response, error := request.StartInfra()
			if ! response.Status {
				fmt.Printf("Error: %s, clause: %s\n", error, response.Message)
			}
			return  error == nil
		}
	}
	case StopInfrastructure: {
		if ! request.CheckInfra() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			response, error := request.StopInfra()
			if ! response.Status {
				fmt.Printf("Error: %s, clause: %s\n", error, response.Message)
			}
			return  error == nil
		}
	}
	case RestartInfrastructure: {
		if ! request.CheckInfra() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			response, error := request.RestartInfra()
			if ! response.Status {
				fmt.Printf("Error: %s, clause: %s\n", error, response.Message)
			}
			return  error == nil
		}
	}
	case DestroyInfrastructure: {
		if ! request.CheckInfra() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			response, error := request.DeleteInfra()
			if ! response.Status {
				fmt.Printf("Error: %s, clause: %s\n", error, response.Message)
			}
			return  error == nil
		}
	}
	case ListInfrastructure: {
		if ! request.CheckInfra() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			response, error := request.StatusInfra()
			if ! response.Status {
				fmt.Printf("Error: %s, clause: %s\n", error, response.Message)
			}
			return  error == nil
		}
	}
	case ListInfrastructures: {
		if ! request.CheckInfra() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			response, error := request.ListInfras()
			if ! response.Status {
				fmt.Printf("Error: %s, clause: %s\n", error, response.Message)
			}
			return  error == nil
		}
	}
	case ListConfigs: {
		if ! request.CheckProject() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			response, error := request.ListProjects()
			if ! response.Status {
				fmt.Printf("Error: %s, clause: %s\n", error, response.Message)
			}
			return  error == nil
		}
	}
	case StatusConfig: {
		if ! request.CheckProject() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			response, error := request.StatusProject()
			if ! response.Status {
				fmt.Printf("Error: %s, clause: %s\n", error, response.Message)
			}
			return  error == nil
		}
	}
	case ImportConfig: {
		if ! request.CheckProject() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			response, error := request.ImportProject()
			if ! response.Status {
				fmt.Printf("Error: %s, clause: %s\n", error, response.Message)
			}
			return  error == nil
		}
	}
	case ExportConfig: {
		if ! request.CheckProject() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			response, error := request.ExportProject()
			if ! response.Status {
				fmt.Printf("Error: %s, clause: %s\n", error, response.Message)
			}
			return  error == nil
		}
	}
	case DefineConfig: {
		if ! request.CheckProject() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			response, error := request.CreateProject()
			if ! response.Status {
				fmt.Printf("Error: %s, clause: %s\n", error, response.Message)
			}
			return  error == nil
		}
	}
	case DeleteConfig: {
		if ! request.CheckProject() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			response, error := request.DeleteProject()
			if ! response.Status {
				fmt.Printf("Error: %s, clause: %s\n", error, response.Message)
			}
			return  error == nil
		}
	}
	case AlterConfig: {
		if ! request.CheckProject() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			response, error := request.AlterProject()
			if ! response.Status {
				fmt.Printf("Error: %s, clause: %s\n", error, response.Message)
			}
			return  error == nil
		}
	}
	case InfoConfig: {
		if ! request.CheckProject() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			request.InfoProject()
		}
		return  true
	}
	default:
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
	}
	return  false
}
