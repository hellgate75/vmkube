package scheduler

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"vmkube/model"
)

var mutex sync.RWMutex

func DumpData(file string, data interface{}, overwrite bool) {
	text := ""
	if !overwrite {
		if strings.Index(file, string(os.PathSeparator)) < 0 {
			file = model.HomeFolder() + string(os.PathSeparator) + file
		}
		mutex.RLock()
		bytes, err := ioutil.ReadFile(file)
		mutex.RUnlock()
		if err == nil {
			text = string(bytes)
			text = fmt.Sprintf("%s", text)
		}
	}
	text += fmt.Sprintf("%s\n", data)
	mutex.Lock()
	ioutil.WriteFile(file, []byte(text), 0777)
	mutex.Unlock()
}
