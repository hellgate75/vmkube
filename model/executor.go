package model


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
	}
	return  false
}
