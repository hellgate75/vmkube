// Provides basic bulding blocks for advanced console UI
//
// Coordinate system:
//
//  1/1---X---->
//   |
//   Y
//   |
//   v
//
// Documentation for ANSI codes: http://en.wikipedia.org/wiki/ANSI_escape_code#Colors
//
// Inspired by: http://www.darkcoding.net/software/pretty-command-line-console-output-on-unix-in-python-and-go-lang/
package term

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"runtime"
	"syscall"
	"unsafe"
)

// Reset all custom styles
const RESET = "\033[0m"

// Reset to default color
const RESET_COLOR = "\033[32m"

// Return curor to start of line and clean it
const RESET_LINE = "\r\033[K"

// List of possible colors
const (
	BLACK = iota
	RED
	GREEN
	YELLOW
	BLUE
	MAGENTA
	CYAN
	WHITE
)

var OutStream *bufio.Writer = bufio.NewWriter(os.Stdout)
var Buffer *bytes.Buffer = new(bytes.Buffer)

func getScreenColor(code int) string {
	return fmt.Sprintf("\033[3%dm", code)
}

func getScreenBgColor(code int) string {
	return fmt.Sprintf("\033[4%dm", code)
}

// Set percent flag: num | PCT
//
// Check percent flag: num & PCT
//
// Reset percent flag: num & 0xFF
const shift = uint(^uint(0)>>63) << 4
const PCT = 0x8000 << shift

type ScreenSize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// Get relative or absolute coorditantes
// To get relative, set PCT flag to number:
//
//      // Get 10% of total width to `x` and 20 to y
//      x, y = tm.GetXY(10|tm.PCT, 20)
//
func GetScreenXY(x int, y int) (int, int) {
	if y == -1 {
		y = ScreenCurrentHeight() + 1
	}
	
	if x&PCT != 0 {
		x = int((x & 0xFF) * ScreenWidth() / 100)
	}
	
	if y&PCT != 0 {
		y = int((y & 0xFF) * ScreenHeight() / 100)
	}
	
	return x, y
}

type sf func(int, string) string

// Apply given transformation func for each line in string
func applyScreenTransform(str string, transform sf) (out string) {
	out = ""
	
	for idx, line := range strings.Split(str, "\n") {
		out += transform(idx, line)
	}
	
	return
}

// Clear screen
func ScreenClear() {
	OutStream.WriteString("\033[2J")
}

// Move cursor to given position
func ScreenMoveCursor(x int, y int) {
	fmt.Fprintf(Buffer, "\033[%d;%dH", x, y)
}

// Move cursor up relative the current position
func ScreenMoveCursorUp(spaces int) {
	fmt.Fprintf(Buffer, "\033[%dA", spaces);
}

// Move cursor down relative the current position
func ScreenMoveCursorDown(spaces int) {
	fmt.Fprintf(Buffer, "\033[%dB", spaces);
}

// Move cursor forward relative the current position
func ScreenMoveCursorForward(spaces int) {
	fmt.Fprintf(Buffer, "\033[%dC", spaces);
}

// Move cursor backward relative the current position
func ScreenMoveCursorBackward(spaces int) {
	fmt.Fprintf(Buffer, "\033[%dD", spaces);
}

// Negative is Left/Top ward positive is Right/Down ward
func ScreenMoveCursorRelative(XSpaces int,YSpaces int) {
	if XSpaces > 0 {
		ScreenMoveCursorDown(XSpaces)
	} else if XSpaces < 0 {
		ScreenMoveCursorUp(-XSpaces)
	}
	if YSpaces > 0 {
		ScreenMoveCursorForward(YSpaces)
	} else if XSpaces < 0 {
		ScreenMoveCursorBackward(-YSpaces)
	}
}

// Move string to possition
func ScreenMoveTo(str string, x int, y int) (out string) {
	x, y = GetScreenXY(x, y)
	
	return applyScreenTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("\033[%d;%dH%s", x+idx, y, line)
	})
}

// Return carrier to start of line
func ScreenResetLine(str string) (out string) {
	return applyScreenTransform(str, func(idx int, line string) string {
		return fmt.Sprintf(RESET_LINE, line)
	})
}

// Make bold
func ScreenBold(str string) string {
	return applyScreenTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("\033[1m%s\033[0m", line)
	})
}

// Apply given color to string:
//
//     tm.Color("RED STRING", tm.RED)
//
func ScreenColor(str string, color int) string {
	return applyScreenTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("%s%s%s", getScreenColor(color), line, RESET)
	})
}

func ScreenHighlight(str, substr string, color int) string {
	hiSubstr := ScreenColor(substr, color)
	return strings.Replace(str, substr, hiSubstr, -1)
}

func ScreenHighlightRegion(str string, from, to, color int) string {
	return str[:from] + ScreenColor(str[from:to], color) + str[to:]
}

// Change background color of string:
//
//     tm.Background("string", tm.RED)
//
func ScreenBackground(str string, color int) string {
	return applyScreenTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("%s%s%s", getScreenBgColor(color), line, RESET)
	})
}

// Get console width
func ScreenWidth() int {
	ws, err := getScreenSize()
	
	if err != nil {
		return -1
	}
	
	return int(ws.Col)
}

func getScreenSize() (*ScreenSize, error) {
	ws := new(ScreenSize)
	
	var _TIOCGWINSZ int64
	
	switch runtime.GOOS {
	case "linux":
		_TIOCGWINSZ = 0x5413
	case "darwin":
		_TIOCGWINSZ = 1074295912
	}
	
	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(_TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)
	
	if int(r1) == -1 {
		fmt.Println("Error:", os.NewSyscallError("GetWinsize", errno))
		return nil, os.NewSyscallError("GetWinsize", errno)
	}
	return ws, nil
}
// Get console height
func ScreenHeight() int {
	ws, err := getScreenSize()
	if err != nil {
		return -1
	}
	return int(ws.Row)
}

// Get current height. Line count in Screen buffer.
func ScreenCurrentHeight() int {
	return strings.Count(Buffer.String(), "\n")
}

// Flush buffer and ensure that it will not overflow screen
func ScreenFlush() {
	for idx, str := range strings.Split(Buffer.String(), "\n") {
		if idx > ScreenHeight() {
			return
		}
		if idx > 0 {
			OutStream.WriteString("\n" + str)
		} else {
			OutStream.WriteString(str)
		}
	}
	
	OutStream.Flush()
	Buffer.Reset()
}

func ScreenPrint(a ...interface{}) {
	fmt.Fprint(Buffer, a...)
}

func ScreenPrintln(a ...interface{}) {
	fmt.Fprintln(Buffer, a...)
}

var cursorHidden bool = false

func ScreenHideCursor() {
	fmt.Print("\033[?25l")
	cursorHidden = true
}

func ScreenShowCursor() {
	if ScreenHasCursorHidden() {
		fmt.Print("\033[?25h")
		cursorHidden = false
	}
}

func ScreenHasCursorHidden() bool {
	return cursorHidden
}

func ScreenPrintf(format string, a ...interface{}) {
	fmt.Fprintf(Buffer, format, a...)
}

func ScreenContext(data string, idx, max int) string {
	var start, end int
	
	if len(data[:idx]) < (max / 2) {
		start = 0
	} else {
		start = idx - max/2
	}
	
	if len(data)-idx < (max / 2) {
		end = len(data) - 1
	} else {
		end = idx + max/2
	}
	
	return data[start:end]
}


func StrPad(instr string, capping int) string {
	strlen := len(instr)
	if strlen == capping  {
		return  instr
	} else  {
		if strlen < capping {
			padding := strings.Repeat(" ", (capping-strlen))
			return  instr + padding
		} else {
			val := instr[0:(capping-2)]
			val += ".."
			return  val
		}
	}
}
