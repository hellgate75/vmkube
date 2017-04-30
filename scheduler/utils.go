package scheduler

import (
	"io/ioutil"
	"fmt"
	"strings"
	"os"
	"vmkube/model"
)

func DumpData(file string, data interface{}, overwrite bool) {
	text := ""
	if !overwrite {
		if strings.Index(file, string(os.PathSeparator)) < 0  {
			file = model.HomeFolder() + string(os.PathSeparator) + file
		}
		bytes, err := ioutil.ReadFile(file)
		if err == nil {
			text = string(bytes)
			text = fmt.Sprintf("%s", text)
		}
	}
	text += fmt.Sprintf("%s\n", data)
	ioutil.WriteFile(file, []byte(text), 0777)

}
