package term

import (
	"fmt"
	"vmkube/utils"
	"time"
)

type TextColorState int

const(
	StateColorWhite TextColorState = iota
	StateColorYellow
	StateColorGreen
	StateColorRed
	StateColorBlack
	StateColorBlue
	StateColorCyan
	StateColorMagenta
	
)

type KeyValueElement struct {
	Id    string
	Name  string
	Value string
	State TextColorState
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
		case StateColorWhite:
			return WHITE
		case StateColorYellow:
			return YELLOW
		case StateColorGreen:
			return GREEN
		case StateColorRed:
			return RED
		case StateColorBlack:
			return BLACK
		case StateColorBlue:
			return BLUE
		case StateColorCyan:
			return CYAN
		default:
			return MAGENTA
	}
	return WHITE
}


func (screenData *KeyValueScreenManager) drawGrid() {
	ScreenClear() // Clear current screen
	if screenData.TextLen == 0 {
		for _,elem := range screenData.Elements {
			if len(elem.Name) > screenData.TextLen {
				screenData.TextLen = len(elem.Name)
			}
		}
	}
	for i := 0; i< len(screenData.Elements); i++  {
		ScreenMoveCursor(i + screenData.OffsetCols + 1, screenData.OffsetRows + 1)
		text := ScreenColor(fmt.Sprintf("%s%s%s", utils.StrPad(screenData.Elements[i].Name, screenData.TextLen), screenData.Separator, ScreenBold(utils.StrPad(screenData.Elements[i].Value, screenData.MessageMaxLen))), screenData.getElementScreenColor(screenData.Elements[i]))
		ScreenPrintln(text)
		ScreenFlush()
	}
	go func(screenData *KeyValueScreenManager){
		for screenData.Active {
			update := <- screenData.CommChannel
			index := screenData.IndexOf(update)
			if index >= 0 {
				screenData.Elements[index] = update
				ScreenMoveCursor(index + screenData.OffsetCols + 1, screenData.OffsetRows + 1)
				text := ScreenColor(fmt.Sprintf("%s%s%s", utils.StrPad(screenData.Elements[index].Name, screenData.TextLen), screenData.Separator, ScreenBold(utils.StrPad(screenData.Elements[index].Value, screenData.MessageMaxLen))), screenData.getElementScreenColor(screenData.Elements[index]))
				ScreenPrintln(text)
				ScreenFlush()
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
	ScreenClear()
}

func (screenData *KeyValueScreenManager) Start() {
	screenData.Active = true
	screenData.drawGrid()
}

//func ProjectOnScreen(grid []ProjectBoxElem, channel chan ProjectBoxElem, onOfCChannel chan bool) {
//	Clear() // Clear current screen
//	for _,elem := range grid {
//		Println(fmt.Sprintf())
//	}
//
//	for {
//		// By moving cursor to top-left position we ensure that console output
//		// will be overwritten each time, instead of adding new.
//		MoveCursor(1, 1)
//
//		Println("Current Time:", time.Now().Format(time.RFC1123))
//
//		Flush() // Call it every time at the end of rendering
//
//		time.Sleep(time.Second)
//	}
//
//}
