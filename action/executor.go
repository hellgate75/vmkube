package action

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
			request.StartInfra()
		}
		return  true
	}
	case StopInfrastructure: {
		if ! request.CheckInfra() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			request.StopInfra()
		}
		return  true
	}
	case RestartInfrastructure: {
		if ! request.CheckInfra() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			request.RestartInfra()
		}
		return  true
	}
	case DestroyInfrastructure: {
		if ! request.CheckInfra() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			request.DeleteInfra()
		}
		return  true
	}
	case ListInfrastructure: {
		if ! request.CheckInfra() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			request.StatusInfra()
		}
		return  true
	}
	case ListInfrastructures: {
		if ! request.CheckInfra() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			request.ListInfras()
		}
		return  true
	}
	case ListConfigs: {
		if ! request.CheckProject() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			request.ListProjects()
		}
		return  true
	}
	case StatusConfig: {
		if ! request.CheckProject() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			request.StatusProject()
		}
		return  true
	}
	case ImportConfig: {
		if ! request.CheckProject() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			request.ImportProject()
		}
		return  true
	}
	case ExportConfig: {
		if ! request.CheckProject() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			request.ExportProject()
		}
		return  true
	}
	case DefineConfig: {
		if ! request.CheckProject() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			request.CreateProject()
		}
		return  true
	}
	case DeleteConfig: {
		if ! request.CheckProject() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			request.DeleteProject()
		}
		return  true
	}
	case AlterConfig: {
		if ! request.CheckProject() {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return  false
		} else  {
			request.AlterProject()
		}
		return  true
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
