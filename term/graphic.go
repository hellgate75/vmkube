package term

import (
	"strconv"
	"strings"
	"fmt"
)

func (Screen *ScreenManager) ProgressBar(description string, current, total, cols int, failure bool) string {
	var maxSpace int = len(strconv.Itoa(total))
	prefix := StrPad(strconv.Itoa(current),maxSpace) + " / " + StrPad(strconv.Itoa(total),maxSpace)
	bar_start := " ["
	bar_end := "] "

	bar_size := cols - len(prefix+bar_start+bar_end)
	amount := int(float32(current) / (float32(total) / float32(bar_size)))
	remain := bar_size - amount

	perc := int((current/total)*100)

	percentage_pre := " ("
	percentage := fmt.Sprintf("%d%s", perc, "%")
	percentage_suf := ")"
	if failure {
		percentage = Screen.Color(percentage, RED)
	} else  if perc == 100 {
		percentage = Screen.Color(percentage, GREEN)
	} else  {
		percentage = Screen.Color(percentage, YELLOW)
	}

var  bar string = ""

	if ! failure {
		bar = strings.Repeat(Screen.Background(" ", YELLOW), amount) + strings.Repeat(" ", remain)
	} else {
		bar = strings.Repeat(Screen.Background(" ", YELLOW), amount) + strings.Repeat(Screen.Bold(Screen.Background(" ", RED)), remain)
	}

	return description + " " + prefix + percentage_pre + Screen.Bold(percentage) + percentage_suf + bar_start + bar + bar_end
}

func (Screen *ScreenManager) ApplyText(text string, line, col int) {
	Screen.MoveCursor(col, line)
	Screen.Print(text)
}
