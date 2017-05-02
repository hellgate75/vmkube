package term

import (
	"strconv"
	"strings"
)

func (Screen *ScreenManager) ProgressBar(current, total, cols int) string {
	prefix := strconv.Itoa(current) + " / " + strconv.Itoa(total)
	bar_start := " ["
	bar_end := "] "

	bar_size := cols - len(prefix+bar_start+bar_end)
	amount := int(float32(current) / (float32(total) / float32(bar_size)))
	remain := bar_size - amount

	bar := strings.Repeat("X", amount) + strings.Repeat(" ", remain)
	return Screen.Bold(prefix) + bar_start + bar + bar_end
}
