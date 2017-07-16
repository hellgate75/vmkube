package nuterm

import (
	"fmt"
	"github.com/satori/go.uuid"
	"time"
	"github.com/hellgate75/vmkube/term"
)

func NewUUIDString() string {
	return uuid.NewV4().String()
}

func main() {
	var elems []term.KeyValueElement = make([]term.KeyValueElement, 0)
	for i := 0; i < 10; i++ {
		elems = append(elems, term.KeyValueElement{
			Id:    NewUUIDString(),
			Name:  fmt.Sprintf("Test Line Number %d", (i + 1)),
			Value: "waiting...",
		})
	}
	manager := term.KeyValueScreenManager{
		Elements:      elems,
		MessageMaxLen: 25,
		Separator:     " ... ",
		OffsetCols:    0,
		OffsetRows:    0,
		TextLen:       0,
		BoldValue:     false,
	}
	manager.Init()
	manager.Start()
	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)
		elems[i].Value = "processing"
		elems[i].State = term.StateColorYellow
		manager.CommChannel <- elems[i]
		time.Sleep(2 * time.Second)
		if i%2 == 0 {
			elems[i].Value = term.Screen.Bold("success!")
			elems[i].State = term.StateColorGreen
			manager.CommChannel <- elems[i]
		} else {
			elems[i].Value = term.Screen.Bold("failed!")
			elems[i].State = term.StateColorRed
			manager.CommChannel <- elems[i]
		}
	}
	manager.Stop(false)
}
