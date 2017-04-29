package term

import (
	"fmt"
	"time"
	"sync"
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
	Id      string
	Name    string
	Value   string
	State   TextColorState
	Ref     interface{}
	Actions int
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
	BoldValue     bool
	inited				bool
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

var mutex sync.Mutex

func (screenData *KeyValueScreenManager) drawGrid() {
	ScreenClear() // Clear current screen
	if screenData.TextLen == 0 {
		for _,elem := range screenData.Elements {
			if len(elem.Name) > screenData.TextLen {
				screenData.TextLen = len(elem.Name)
			}
		}
	}
	screenHeight := ScreenHeight()
	rows := len(screenData.Elements)
	for i := 0; i< rows; i++  {
		ScreenMoveCursor(i + screenData.OffsetCols + 1, screenData.OffsetRows + 1)
		var text string
		if screenData.BoldValue {
			text = ScreenColor(fmt.Sprintf("%s%s%s", StrPad(screenData.Elements[i].Name, screenData.TextLen), screenData.Separator, ScreenBold(StrPad(screenData.Elements[i].Value, screenData.MessageMaxLen))), screenData.getElementScreenColor(screenData.Elements[i]))
		} else {
			text = ScreenColor(fmt.Sprintf("%s%s%s", StrPad(screenData.Elements[i].Name, screenData.TextLen), screenData.Separator, StrPad(screenData.Elements[i].Value, screenData.MessageMaxLen)), screenData.getElementScreenColor(screenData.Elements[i]))
		}
		ScreenPrintln(text)
		ScreenFlush()
	}
	if rows > screenHeight {
		screenData.OffsetRows = screenHeight - rows
	}
	go func(screenData *KeyValueScreenManager){
		for screenData.Active {
			update := <- screenData.CommChannel
			index := screenData.IndexOf(update)
			if index >= 0 {
				mutex.Lock()
				screenData.Elements[index] = update
				ScreenMoveCursor(index + screenData.OffsetCols + 1, screenData.OffsetRows + 1)
				var text string
				if screenData.BoldValue {
					text = ScreenColor(fmt.Sprintf("%s%s%s", StrPad(screenData.Elements[index].Name, screenData.TextLen), screenData.Separator, ScreenBold(StrPad(screenData.Elements[index].Value, screenData.MessageMaxLen))), screenData.getElementScreenColor(screenData.Elements[index]))
				} else {
					text = ScreenColor(fmt.Sprintf("%s%s%s", StrPad(screenData.Elements[index].Name, screenData.TextLen), screenData.Separator, StrPad(screenData.Elements[index].Value, screenData.MessageMaxLen)), screenData.getElementScreenColor(screenData.Elements[index]))
				}
				ScreenPrintln(text)
				ScreenFlush()
				mutex.Unlock()
			}
		}
	}(screenData)
}

func (screenData *KeyValueScreenManager) Init() {
	screenData.CommChannel = make(chan KeyValueElement)
	screenData.CtrlChannel = make(chan bool)
	if screenData.Separator == "" {
		screenData.Separator = " "
	}
	go func(screenData *KeyValueScreenManager){
		screenData.Active = <- screenData.CtrlChannel
		if ! screenData.Active {
			close(screenData.CtrlChannel)
			close(screenData.CommChannel)
		} else {
			screenData.Start()
		}
	}(screenData)
}

func (screenData *KeyValueScreenManager) IndexOf(elem KeyValueElement) int {
	for i := 0; i< len(screenData.Elements); i++  {
		if screenData.Elements[i].Id == elem.Id {
			return i
		}
	}
	return -1
}

func (screenData *KeyValueScreenManager) Remove(elem KeyValueElement) {
	index := screenData.IndexOf(elem)
	if index >= 0 {
		if len(screenData.Elements) == 1 {
			screenData.Elements = make([]KeyValueElement, 0)
		} else if index == 0 {
			screenData.Elements = screenData.Elements[1:]
		} else if index == len(screenData.Elements)-1 {
			screenData.Elements = screenData.Elements[:index]
		} else {
			screenData.Elements = screenData.Elements[0:index]
			screenData.Elements = append(screenData.Elements, screenData.Elements[(index+1):]...)
		}
	}
}

func (screenData *KeyValueScreenManager) Stop(clearScreen bool) {
	if screenData.Active {
		time.Sleep(1 * time.Second)
		screenData.CtrlChannel <- false
		if clearScreen {
			ScreenClear()
		}
		ScreenShowCursor()
	}
}

func (screenData *KeyValueScreenManager) Start() {
	if ! screenData.Active {
		screenData.CtrlChannel <- true
	} else  if ! screenData.inited {
		screenData.Active = true
		screenData.drawGrid()
		ScreenHideCursor()
	}
}
