package term

import (
	"strconv"
	"strings"
	"fmt"
	"time"
)

func (Screen *ScreenManager) ProgressBar(description string, currentValue, totalValue, numberOfColumns int, failure bool) string {
	var maxSpace int = len(strconv.Itoa(totalValue))
	prefix := StrFill(strconv.Itoa(currentValue),maxSpace) + " / " + StrFill(strconv.Itoa(totalValue),maxSpace)
	bar_start := " ["
	bar_end := "] "

	bar_size := numberOfColumns - len(prefix+bar_start+bar_end)
	amount := int(float32(currentValue) / (float32(totalValue) / float32(bar_size)))
	remain := bar_size - amount

	percString := int(((currentValue * 100) / totalValue))

	percentage_pre := " ("
	percentage := StrFill(fmt.Sprintf("%d%s", percString, "%"), 4)
	percentage_suf := ")"
	if failure {
		percentage = Screen.Color(percentage, RED)
	} else  if percString == 100 {
		percentage = Screen.Color(percentage, GREEN)
	} else  {
		percentage = Screen.Color(percentage, YELLOW)
	}

var  bar string = ""

	if ! failure {
		bar = strings.Repeat(Screen.Background(" ", YELLOW), amount) + strings.Repeat(" ", remain)
	} else {
		bar = strings.Repeat(Screen.Background(" ", WHITE), amount) + strings.Repeat(Screen.Bold(Screen.Background(" ", RED)), remain)
	}

	return description + " " + prefix + percentage_pre + Screen.Bold(percentage) + percentage_suf + bar_start + bar + bar_end
}

func (Screen *ScreenManager) WriteProgressBar(description string, currentValue, totalValue, numberOfColumns, positionColumn, positionRow int, failure bool) {
	Screen.ApplyText(Screen.ProgressBar(description, currentValue, totalValue, numberOfColumns, failure), positionRow, positionColumn)
}

func (Screen *ScreenManager) ApplyText(text string, line, col int) {
	Screen.MoveCursor(col+1, line+1)
	Screen.Println(text)
	Screen.Flush()
}

const PRGRESS_BAR_MAX_CHANNEL_TIMEOUT = 900

type ProgressBar struct {
	Running					bool
	Errors					bool
	MaxValues				int
	BarSteps				int
	Current					int
	ScreenRow				int
	ScreenCol				int
	PostReset				bool
	ResetRow				int
	ResetCol				int
	Prefix					string
	ClearScreen			bool
	HideCursor			bool
	HasCallBack			bool
	FinalCallBack		func()
}

func (Bar *ProgressBar) Start(IncreaseChannel	chan int, FailureChannel	chan bool ) bool {
	Bar.Running = true
	Bar.Errors = false
	if Bar.ClearScreen {
		Screen.Clear()
	}
	if Bar.HideCursor {
		Screen.HideCursor()
	}
	Screen.WriteProgressBar(Bar.Prefix, Bar.Current, Bar.MaxValues, Bar.BarSteps, Bar.ScreenCol, Bar.ScreenRow, Bar.Errors)
	go func() {
		for Bar.Current < Bar.MaxValues && Bar.Running {
			select {
			case newValue := <- IncreaseChannel :
				go func(newValue int) {
					mutex.Lock()
					Bar.Current += newValue
					Screen.WriteProgressBar(Bar.Prefix, Bar.Current, Bar.MaxValues, Bar.BarSteps, Bar.ScreenCol, Bar.ScreenRow, Bar.Errors)
					if Bar.PostReset {
						Screen.MoveCursor(Bar.ResetCol, Bar.ResetRow)
						Screen.Flush()
					}
					mutex.Unlock()
				}(newValue)
			case <-time.After(PRGRESS_BAR_MAX_CHANNEL_TIMEOUT * time.Second):
			}
		}
		Bar.Running = false
	}()
	go func() {
		for Bar.Current < Bar.MaxValues && Bar.Running {
			select {
			case newValue := <- FailureChannel :
				go func(newValue bool) {
					mutex.Lock()
					Bar.Errors = newValue
					Screen.WriteProgressBar(Bar.Prefix, Bar.Current, Bar.MaxValues, Bar.BarSteps, Bar.ScreenCol, Bar.ScreenRow, Bar.Errors)
					if Bar.PostReset {
						Screen.MoveCursor(Bar.ResetCol, Bar.ResetRow)
						Screen.Flush()
					}
					mutex.Unlock()
				}(newValue)
			case <-time.After(PRGRESS_BAR_MAX_CHANNEL_TIMEOUT * time.Second):
			}
		}
		Bar.Running = false
	}()
	return  true
}

func (Bar *ProgressBar) Stop() {
	if Bar.HideCursor {
		Screen.ShowCursor()
	}
	Bar.Running = false
	if Bar.HasCallBack {
		go Bar.FinalCallBack()
	}
	if Bar.PostReset {
		Screen.MoveCursor(Bar.ResetCol, Bar.ResetRow)
		Screen.Flush()
	}
}