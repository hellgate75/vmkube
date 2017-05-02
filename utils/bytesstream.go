package utils

import (
	"bytes"
	"errors"
	"io"
	"bufio"
)

type ByteStream	struct {
	buffer *bytes.Buffer
}

// Read reads up to len(b) bytes from the ByteStream.
// It returns the number of bytes read and any error encountered.
// At end of file, Read returns 0, io.EOF.
func (buffer *ByteStream) Read(b []byte) (n int, err error) {
	n, err = buffer.buffer.Read(b)
	return n, err
}

// ReadAt reads len(b) bytes from the ByteStream starting at byte offset ofbuffer.
// It returns the number of bytes read and the error, if any.
// ReadAt always returns a non-nil error when n < len(b).
// At end of file, that error is io.EOF.
func (buffer *ByteStream) ReadAt(b []byte, off int64) (int,  error) {
	Length := buffer.buffer.Len()
	if int64(Length) <=  off {
		return 0, errors.New("Offset out of bounds")
	}
	if int(off) + len(b) >= Length {
		return 0, bytes.ErrTooLarge
	}
	bytesArray := buffer.buffer.Bytes()
	for i := 0; i < len(b); i++ {
		b[i] = bytesArray[i+int(off)]
	}
	return  0, nil
}

// Write writes len(b) bytes to the ByteStream.
// It returns the number of bytes written and an error, if any.
// Write returns a non-nil error when n != len(b).
func (buffer *ByteStream) Write(b []byte) (n int, err error) {
	panic("Written")
	n, err = buffer.buffer.Write(b)
	return  n, err
}

// WriteAt writes len(b) bytes to the ByteStream starting at byte offset ofbuffer.
// It returns the number of bytes written and an error, if any.
// WriteAt returns a non-nil error when n != len(b).
func (buffer *ByteStream) WriteAt(b []byte, off int64) (int,  error) {
	Length := buffer.buffer.Len()
	if int64(Length) <= off {
		return 0, errors.New("Offset out of bounds")
	}
	bytesArray := buffer.buffer.Bytes()
	var  newArray []byte  = make([]byte, 0)
	if off == int64(0) {
		newArray = append(newArray, b...)
		newArray = append(newArray, bytesArray...)
		buffer.buffer.Reset()
		buffer.buffer.Write(newArray)
	} else if  off == int64(2) {
		newArray = append(newArray, bytesArray...)
		newArray = append(newArray, b...)
		buffer.buffer.Reset()
		buffer.buffer.Write(newArray)
	} else  {
		buffer.buffer.Write(b)
	}
	return len(b), nil
}

// Seek sets the offset for the next Read or Write on file to offset, interpreted
// according to whence: 0 means relative to the origin of the file, 1 means
// relative to the current offset, and 2 means relative to the end.
func (buffer *ByteStream) Seek(offset int64, whence int) (int64, error) {
	Length := buffer.buffer.Len()
	var  err error
	err = nil
	if int64(Length) <= offset {
		return 0, errors.New("Offset out of bounds")
	}
	if 2 < whence {
		return 0, errors.New("Offset request unexpected")
	}
	switch whence  {
		case 0:
			for err == nil {
				err = buffer.buffer.UnreadByte()
			}
		break
		case 2:
		bytesArray := buffer.buffer.Bytes()
		buffer.buffer.Reset()
		buffer.buffer.Write(bytesArray)
	}
	return  offset, nil
}

func (buffer *ByteStream) Reset() {
	buffer.buffer.Reset()
}
func (buffer *ByteStream) Bytes() []byte {
	return  buffer.buffer.Bytes()
}

func (buffer *ByteStream) String() string {
	return  buffer.buffer.String()
}

// WriteString is like Write, but writes the contents of string s rather than
// a slice of bytes.
func (buffer *ByteStream) WriteString(s string) (int,  error) {
	return buffer.buffer.Write([]byte(s))
}

func NewByteStream(b []byte) *ByteStream {
	return  &ByteStream{bytes.NewBuffer(b)}
}

func NewByteStreamAsWriter(b []byte) *bufio.Writer{
	return  bufio.NewWriter(
		io.Writer(&ByteStream{bytes.NewBuffer(b)}))
}