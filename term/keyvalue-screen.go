package term

import (
	tm "github.com/buger/goterm"
	"fmt"
	"vmkube/utils"
	"time"
)

type ElementState int

const(
	State_Element_Waiting ElementState = iota
	State_Element_Partial
	State_Element_Complete
	State_Element_Error
)

type KeyValueElement struct {
	Id        string
	Name      string
	Value     string
	State     ElementState
}

type KeyValueScreenManager struct {
	Elements      []KeyValueElement
	CommChannel   chan KeyValueElement
	CtrlChannel   chan bool
	Active        bool
	TextLen       int
	MessageMaxLen int
	OffsetCols    int
	OffsetRows    int
	Separator     string
}

func (screenData *KeyValueScreenManager) getElementScreenColor(elem KeyValueElement) int {
	switch elem.State {
		case State_Element_Waiting:
			return tm.WHITE
		case State_Element_Partial:
			return tm.YELLOW
		case State_Element_Complete:
			return tm.GREEN
		default:
			return tm.RED
	}
	return tm.WHITE
}


func (screenData *KeyValueScreenManager) drawGrid() {
	tm.Clear() // Clear current screen
	if screenData.TextLen == 0 {
		for _,elem := range screenData.Elements {
			if len(elem.Name) > screenData.TextLen {
				screenData.TextLen = len(elem.Name)
			}
		}
	}
	for i := 0; i< len(screenData.Elements); i++  {
		tm.MoveCursor(i + screenData.OffsetCols + 1, screenData.OffsetRows + 1)
		text := tm.Color(fmt.Sprintf("%s%s%s", utils.StrPad(screenData.Elements[i].Name, screenData.TextLen), screenData.Separator, utils.StrPad(screenData.Elements[i].Value, screenData.MessageMaxLen)), screenData.getElementScreenColor(screenData.Elements[i]))
		tm.Println(text)
		tm.Flush()
	}
	go func(screenData *KeyValueScreenManager){
		for screenData.Active {
			update := <- screenData.CommChannel
			index := screenData.IndexOf(update)
			if index >= 0 {
				screenData.Elements[index] = update
				tm.MoveCursor(index + screenData.OffsetCols + 1, screenData.OffsetRows + 1)
				text := tm.Color(fmt.Sprintf("%s%s%s", utils.StrPad(screenData.Elements[index].Name, screenData.TextLen), screenData.Separator, utils.StrPad(screenData.Elements[index].Value, screenData.MessageMaxLen)), screenData.getElementScreenColor(screenData.Elements[index]))
				tm.Println(text)
				tm.Flush()
			}
		}
	}(screenData)
	go func(screenData *KeyValueScreenManager){
		screenData.Active = <- screenData.CtrlChannel
		if ! screenData.Active {
			close(screenData.CtrlChannel)
			close(screenData.CommChannel)
			screenData.Stop()
		} else {
			screenData.Start()
		}
	}(screenData)
	
}

func (screenData *KeyValueScreenManager) Init() {
	screenData.CommChannel = make(chan KeyValueElement)
	screenData.CtrlChannel = make(chan bool)
	if screenData.Separator == "" {
		screenData.Separator = " "
	}
}

func (screenData *KeyValueScreenManager) IndexOf(elem KeyValueElement) int {
	for i := 0; i< len(screenData.Elements); i++  {
		if screenData.Elements[i].Id == elem.Id {
			return i
		}
	}
	return -1
}

func (screenData *KeyValueScreenManager) Stop() {
	time.Sleep(1 * time.Second)
	screenData.CtrlChannel <- false
	tm.Clear()
}

func (screenData *KeyValueScreenManager) Start() {
	screenData.Active = true
	screenData.drawGrid()
}

//func ProjectOnScreen(grid []ProjectBoxElem, channel chan ProjectBoxElem, onOfCChannel chan bool) {
//	tm.Clear() // Clear current screen
//	for _,elem := range grid {
//		tm.Println(fmt.Sprintf())
//	}
//
//	for {
//		// By moving cursor to top-left position we ensure that console output
//		// will be overwritten each time, instead of adding new.
//		tm.MoveCursor(1, 1)
//
//		tm.Println("Current Time:", time.Now().Format(time.RFC1123))
//
//		tm.Flush() // Call it every time at the end of rendering
//
//		time.Sleep(time.Second)
//	}
//
//}
