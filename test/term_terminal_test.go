package test

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"github.com/hellgate75/vmkube/term"
	"github.com/hellgate75/vmkube/utils"
)

var TermBuffer *utils.ByteStream = utils.NewByteStream([]byte{})
var writer *bufio.Writer = utils.NewByteStreamAsWriter([]byte{})

func TestBeforeAllTests(t *testing.T) {
	term.Screen.OutStream = writer
}

func TestTerminalDefineColorString(t *testing.T) {
	text := "Test"
	outText := term.Screen.Color(text, term.BLUE)
	expected := fmt.Sprintf("%s%s%s", fmt.Sprintf(term.COLOR_SELECTOR, int(term.BLUE)), text, term.RESET)
	assert.Equal(t, expected, outText, "Color must be stored correctly")
}

func TestTerminalDefineBoldString(t *testing.T) {
	text := "Test"
	outText := term.Screen.Bold(text)
	expected := fmt.Sprintf(term.APPLY_BOLD_EFFECT, text)
	assert.Equal(t, expected, outText, "Bold effect must be stored correctly")
}

func TestTerminalDefineBgColorString(t *testing.T) {
	text := "Test"
	outText := term.Screen.Background(text, term.BLUE)
	expected := fmt.Sprintf("%s%s%s", fmt.Sprintf(term.BG_COLOR_SELECTOR, int(term.BLUE)), text, term.RESET)
	assert.Equal(t, expected, outText, "Background color must be stored correctly")
}

func TestTerminalWriteBufferStyle(t *testing.T) {
	text := "Test"
	outText := term.Screen.Color(text, term.BLUE)
	term.Screen.Print(outText)
	expected := fmt.Sprintf("%s%s%s", fmt.Sprintf(term.COLOR_SELECTOR, int(term.BLUE)), text, term.RESET)
	assert.Equal(t, expected, term.Screen.Buffer.String(), "Color must be written on buffer correctly")
	term.Screen.Buffer.Reset()
	term.Screen.Flush()
	TermBuffer.Reset()
	outText = term.Screen.Background(text, term.BLUE)
	term.Screen.Print(outText)
	expected = fmt.Sprintf("%s%s%s", fmt.Sprintf(term.BG_COLOR_SELECTOR, int(term.BLUE)), text, term.RESET)
	assert.Equal(t, expected, term.Screen.Buffer.String(), "Background color must be written on buffer correctly")
	term.Screen.Buffer.Reset()
	term.Screen.Flush()
	TermBuffer.Reset()
	outText = term.Screen.Bold(text)
	term.Screen.Print(outText)
	expected = fmt.Sprintf(term.APPLY_BOLD_EFFECT, text)
	assert.Equal(t, expected, term.Screen.Buffer.String(), "Bold effect must be written on buffer correctly")
	term.Screen.Buffer.Reset()
	term.Screen.Flush()
	TermBuffer.Reset()
}

//func TestTerminalWriteFlushStyle(t *testing.T) {
//	text := "Test"
//	outText := term.Screen.Color(text, term.BLUE)
//	term.Screen.Print(outText)
//	term.Screen.Flush()
//	expected := fmt.Sprintf("%s%s%s", fmt.Sprintf(term.COLOR_SELECTOR, int(term.BLUE)), text, term.RESET)
//	assert.Equal(t, expected, TermBuffer.String(), "Color must be ritten on device correctly")
//	term.Screen.Buffer.Reset()
//
//}

func TestAfterAllTests(t *testing.T) {
	term.Screen.OutStream = bufio.NewWriter(os.Stdout)
}
