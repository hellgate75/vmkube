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
	"syscall"
	"unsafe"
)

// Reset all custom styles
const RESET = "\033[0m"

// Reset to default color
const RESET_COLOR = "\033[32m"

// Return cursor to start of line and clean it
const RESET_LINE = "\r\033[K"

// Define color selector base code
const COLOR_SELECTOR = "\033[3%dm"

// Define background color selector base code
const BG_COLOR_SELECTOR = "\033[4%dm"

// Reset Screen base position after last row, scrolling down screen
const CLEAR_SCREEN = "\033[2J"

// Move Cursor on Screen since first visible row of x rows and y columns
const MOVE_CURSOR_TO_COORD = "\033[%d;%dH"

// Make a transformation to move Cursor on Screen since first visible row of x rows and y columns since last position
const MOVE_CURSOR_RELATIVE_OF = "\033[%d;%dH%s"

// Move Cursor up on Screen since first visible row of x rows
const MOVE_CURSOR_UP_ROWS = "\033[%dA"

// Move Cursor down on Screen since first visible row of x rows
const MOVE_CURSOR_DOWN_ROWS = "\033[%dB"

// Move Cursor forward on Screen since first visible row of y columns
const MOVE_CURSOR_FORWARD_COLUMNS = "\033[%dC"

// Move Cursor backward on Screen since first visible row of y columns
const MOVE_CURSOR_BACKWARD_COLUMNS = "\033[%dD"

// Apply bold effect to string
const APPLY_BOLD_EFFECT = "\033[1m%s\033[0m"

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

type winsize struct {
	rows    uint16
	cols    uint16
	xpixels uint16
	ypixels uint16
}

type ScreenManager struct {
	OutStream *bufio.Writer
	Buffer    *bytes.Buffer
	AutoFlush bool
	X         int
	Y         int
}

var Screen ScreenManager = ScreenManager{
	OutStream: bufio.NewWriter(os.Stdout),
	Buffer:    new(bytes.Buffer),
	AutoFlush: false,
}

func getScreenColor(code int) string {
	return fmt.Sprintf(COLOR_SELECTOR, code)
}

func getScreenBgColor(code int) string {
	return fmt.Sprintf(BG_COLOR_SELECTOR, code)
}

// Get relative or absolute coorditantes
// To get relative, set PCT flag to number:
//
//      // Get 10% of total width to `x` and 20 to y
//      x, y = tm.GetXY(10|tm.PCT, 20)
//
func (Screen *ScreenManager) GetScreenXY(x int, y int) (int, int) {
	if y == -1 {
		y = Screen.CurrentHeight() + 1
	}

	if x&PCT != 0 {
		x = int((x & 0xFF) * Screen.Width() / 100)
	}

	if y&PCT != 0 {
		y = int((y & 0xFF) * Screen.Height() / 100)
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
func (Screen *ScreenManager) Clear() {
	Screen.OutStream.WriteString(CLEAR_SCREEN)
}

// Move cursor to given position
func (Screen *ScreenManager) MoveCursor(x int, y int) {
	fmt.Fprintf(Screen.Buffer, MOVE_CURSOR_TO_COORD, x, y)
	if Screen.AutoFlush {
		Screen.Flush()
	}
}

// Move cursor up relative the current position
func (Screen *ScreenManager) MoveCursorUp(spaces int) {
	fmt.Fprintf(Screen.Buffer, MOVE_CURSOR_UP_ROWS, spaces)
	if Screen.AutoFlush {
		Screen.Flush()
	}
}

// Move cursor down relative the current position
func (Screen *ScreenManager) MoveCursorDown(spaces int) {
	fmt.Fprintf(Screen.Buffer, MOVE_CURSOR_DOWN_ROWS, spaces)
	if Screen.AutoFlush {
		Screen.Flush()
	}
}

// Move cursor forward relative the current position
func (Screen *ScreenManager) MoveCursorForward(spaces int) {
	fmt.Fprintf(Screen.Buffer, MOVE_CURSOR_FORWARD_COLUMNS, spaces)
	if Screen.AutoFlush {
		Screen.Flush()
	}
}

// Move cursor backward relative the current position
func (Screen *ScreenManager) MoveCursorBackward(spaces int) {
	fmt.Fprintf(Screen.Buffer, MOVE_CURSOR_BACKWARD_COLUMNS, spaces)
	if Screen.AutoFlush {
		Screen.Flush()
	}
}

// Negative is Left/Top ward positive is Right/Down ward
func (Screen *ScreenManager) MoveCursorRelative(XSpaces int, YSpaces int) {
	if XSpaces > 0 {
		Screen.MoveCursorDown(XSpaces)
	} else if XSpaces < 0 {
		Screen.MoveCursorUp(-XSpaces)
	}
	if YSpaces > 0 {
		Screen.MoveCursorForward(YSpaces)
	} else if XSpaces < 0 {
		Screen.MoveCursorBackward(-YSpaces)
	}
}

// Move string to position
func (Screen *ScreenManager) MoveTo(str string, x int, y int) (out string) {
	x, y = Screen.GetScreenXY(x, y)

	return applyScreenTransform(str, func(idx int, line string) string {
		return fmt.Sprintf(MOVE_CURSOR_RELATIVE_OF, x+idx, y, line)
	})
}

// Return carrier to start of line
func (Screen *ScreenManager) ResetLine(str string) (out string) {
	return applyScreenTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("%s"+RESET_LINE, line)
	})
}

// Make bold
func (Screen *ScreenManager) Bold(str string) string {
	return applyScreenTransform(str, func(idx int, line string) string {
		return fmt.Sprintf(APPLY_BOLD_EFFECT, line)
	})
}

// Apply given color to string:
//
//     tm.Color("RED STRING", tm.RED)
//
func (Screen *ScreenManager) Color(str string, color int) string {
	return applyScreenTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("%s%s%s", getScreenColor(color), line, RESET)
	})
}

func (Screen *ScreenManager) Highlight(str, substr string, color int) string {
	hiSubstr := Screen.Color(substr, color)
	return strings.Replace(str, substr, hiSubstr, -1)
}

func (Screen *ScreenManager) HighlightRegion(str string, from, to, color int) string {
	return str[:from] + Screen.Color(str[from:to], color) + str[to:]
}

// Change background color of string:
//
//     tm.Background("string", tm.RED)
//
func (Screen *ScreenManager) Background(str string, color int) string {
	return applyScreenTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("%s%s%s", getScreenBgColor(color), line, RESET)
	})
}

// Get console width
func (Screen *ScreenManager) Width() int {
	ws, err := Screen.getScreenSize()

	if err != nil {
		return -1
	}

	return int(ws.cols)
}

func (Screen *ScreenManager) getTermSize(fd uintptr) (*winsize, error) {
	var sz winsize
	r1, _, errNo := syscall.Syscall(syscall.SYS_IOCTL,
		fd, uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&sz)))
	if int(r1) == -1 {
		fmt.Println("Error:", os.NewSyscallError("GetWinsize", errNo))
		return nil, os.NewSyscallError("GetWinsize", errNo)
	}
	return &sz, nil
}

func (Screen *ScreenManager) getScreenSize() (*winsize, error) {
	//ws := new(ScreenSize)
	//
	//var _TIOCGWINSZ int64
	//
	//switch runtime.GOOS {
	//case "linux":
	//	_TIOCGWINSZ = 0x5413
	//case "darwin":
	//	_TIOCGWINSZ = 1074295912
	//case "windows":
	//default:
	//	_TIOCGWINSZ = syscall.TIOCGWINSZ
	//}
	//
	//r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
	//	uintptr(syscall.Stdin),
	//	uintptr(_TIOCGWINSZ),
	//	uintptr(unsafe.Pointer(ws)),
	//)
	//
	//if int(r1) == -1 {
	//	fmt.Println("Error:", os.NewSyscallError("GetWinsize", errno))
	//	return nil, os.NewSyscallError("GetWinsize", errno)
	//}
	//return ws, nil
	return Screen.getTermSize(os.Stdout.Fd())
}

// Get console height
func (Screen *ScreenManager) Height() int {
	ws, err := Screen.getScreenSize()
	if err != nil {
		return -1
	}
	return int(ws.rows)
}

// Get current height. Line count in Screen buffer.
func (Screen *ScreenManager) CurrentHeight() int {
	return strings.Count(Screen.Buffer.String(), "\n")
}

// Flush buffer and ensure that it will not overflow screen
func (Screen *ScreenManager) Flush() {
	for idx, str := range strings.Split(Screen.Buffer.String(), "\n") {
		if idx > Screen.Height() {
			return
		}
		if idx > 0 {
			Screen.OutStream.WriteString("\n" + str)
		} else {
			Screen.OutStream.WriteString(str)
		}
	}

	Screen.OutStream.Flush()
	Screen.Buffer.Reset()
}

func (Screen *ScreenManager) Print(a ...interface{}) {
	fmt.Fprint(Screen.Buffer, a...)
	if Screen.AutoFlush {
		Screen.Flush()
	}
}

func (Screen *ScreenManager) Println(a ...interface{}) {
	fmt.Fprintln(Screen.Buffer, a...)
	if Screen.AutoFlush {
		Screen.Flush()
	}
}

var cursorHidden bool = false

func (Screen *ScreenManager) HideCursor() {
	Screen.OutStream.WriteString("\033[?25l")
	cursorHidden = true
}

func (Screen *ScreenManager) ShowCursor() {
	Screen.OutStream.WriteString("\033[?25h")
	cursorHidden = false
}

func (Screen *ScreenManager) HasCursorHidden() bool {
	return cursorHidden
}

func (Screen *ScreenManager) Printf(format string, a ...interface{}) {
	fmt.Fprintf(Screen.Buffer, format, a...)
}

func (Screen *ScreenManager) Context(data string, idx, max int) string {
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
	if strlen == capping {
		return instr
	} else {
		if strlen < capping {
			padding := strings.Repeat(" ", (capping - strlen))
			return instr + padding
		} else {
			val := instr[0:(capping - 2)]
			val += ".."
			return val
		}
	}
}
