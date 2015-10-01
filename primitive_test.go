package amqp

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

var validReadTests = []struct {
	input    []byte
	expected interface{}
}{
	{[]byte{0x40}, nil},
	{[]byte{0x41}, true},
	{[]byte{0x42}, false},
	{[]byte{0x43}, uint32(0)},
	{[]byte{0x44}, uint64(0)},
	{[]byte{0x50, 15}, uint8(15)},
	{[]byte{0x50, 200}, uint8(200)},
	{[]byte{0x51, 22}, int8(22)},
	{[]byte{0x51, 0xA0}, int8(-96)},
	{[]byte{0x52, 225}, uint32(225)},
	{[]byte{0x53, 230}, uint64(230)},
	{[]byte{0x54, 0xA1}, int32(-95)},
	{[]byte{0x55, 0xA3}, int64(-93)},
	{[]byte{0x56, 0x00}, false},
	{[]byte{0x56, 0x01}, true},
	{[]byte{0x60, 0x01, 0x02}, uint16(258)},
	{[]byte{0x60, 0xFF, 0xFE}, uint16(65534)},
	{[]byte{0x61, 0x02, 0x03}, int16(515)},
	{[]byte{0x61, 0xBA, 0x21}, int16(-17887)},
	{[]byte{0x70, 0x01, 0x02, 0x03, 0x04}, uint32(16909060)},
	{[]byte{0x70, 0xCC, 0xAA, 0x88, 0x66}, uint32(3433728102)},
	{[]byte{0x71, 0x01, 0x02, 0x03, 0x04}, int32(16909060)},
	{[]byte{0x71, 0xCC, 0xAA, 0x88, 0x66}, int32(-861239194)},
	{[]byte{0x72, 0x40, 0x49, 0x0F, 0xD0}, float32(3.14159)},
	{[]byte{0x80, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, uint64(72623859790382856)},
	{[]byte{0x80, 0xFF, 0xDD, 0xBB, 0x99, 0x77, 0x55, 0x33, 0x11}, uint64(18437098717331141393)},
	{[]byte{0x81, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, int64(72623859790382856)},
	{[]byte{0x81, 0xFF, 0xDD, 0xBB, 0x99, 0x77, 0x55, 0x33, 0x11}, int64(-9645356378410223)},
	{[]byte{0x82, 0xC0, 0xA0, 0x81, 0xF9, 0xAD, 0xD3, 0xD6, 0x00}, float64(-2112.9876543234567)},
}

var validStringTests = []struct {
	input    []byte
	expected string
}{
	{[]byte{0x98, 0x17, 0xad, 0x53, 0x20, 0x67, 0xd8, 0x11, 0xe5,
		0xbc, 0xbe, 0x00, 0x02, 0xa5, 0xd5, 0xc5, 0x1b}, "17ad5320-67d8-11e5-bcbe-0002a5d5c51b"},
}

var invalidReadTests = []struct {
	input    []byte
	expected error
}{
	{nil, io.EOF},
	{[]byte{0x30}, ErrInvalidCode},
	{[]byte{0x50}, io.EOF},
	{[]byte{0x51}, io.EOF},
	{[]byte{0x52}, io.EOF},
	{[]byte{0x53}, io.EOF},
	{[]byte{0x56}, io.EOF},
	{[]byte{0x56, 0x02}, ErrInvalidFormat},
	{[]byte{0x60}, io.EOF},
	{[]byte{0x60, 0x01}, io.ErrUnexpectedEOF},
	{[]byte{0x61}, io.EOF},
	{[]byte{0x61, 0x00}, io.ErrUnexpectedEOF},
	{[]byte{0x70}, io.EOF},
	{[]byte{0x70, 0x00, 0x00, 0x00}, io.ErrUnexpectedEOF},
	{[]byte{0x71}, io.EOF},
	{[]byte{0x71, 0x00, 0x00, 0x00}, io.ErrUnexpectedEOF},
	{[]byte{0x72}, io.EOF},
	{[]byte{0x72, 0x00, 0x00, 0x00}, io.ErrUnexpectedEOF},
	{[]byte{0x80}, io.EOF},
	{[]byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, io.ErrUnexpectedEOF},
	{[]byte{0x81}, io.EOF},
	{[]byte{0x81, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, io.ErrUnexpectedEOF},
	{[]byte{0x82}, io.EOF},
	{[]byte{0x82, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, io.ErrUnexpectedEOF},
}

var validWriteTests = []struct {
	name     string
	f        func(io.Writer) error
	expected []byte
}{
	{"WriteNull",
		func(b io.Writer) error {
			return WriteNull(b)
		}, []byte{0x40}},
	{"WriteBoolean True",
		func(b io.Writer) error {
			return WriteBoolean(b, true)
		}, []byte{0x41}},
	{"WriteBoolean False",
		func(b io.Writer) error {
			return WriteBoolean(b, false)
		}, []byte{0x42}},
	{"WriteUByte Low",
		func(b io.Writer) error {
			return WriteUByte(b, 15)
		}, []byte{0x50, 15}},
}

var invalidWriteTests = []struct {
	name     string
	f        func() error
	expected error
}{
	{"WriteBoolean",
		func() error {
			b := InvalidWriter(true)
			return WriteBoolean(b, true)
		}, io.EOF},
	{"WriteUByte",
		func() error {
			b := InvalidWriter(true)
			return WriteUByte(b, 15)
		}, io.EOF},
}

type InvalidWriter bool

func (b InvalidWriter) Write(p []byte) (int, error) {
	return len(p), io.EOF
}

func (b InvalidWriter) WriteByte(c byte) error {
	return io.EOF
}

func TestValidRead(t *testing.T) {
	var b bytes.Buffer
	for _, tt := range validReadTests {
		b.Reset()
		b.Write(tt.input)
		res, err := ReadNext(&b)
		if err != nil {
			t.Errorf("Reading value failed with input %v: %v", tt.input, err)
		}
		if res != tt.expected {
			t.Errorf("Unexpected value for input %v: %v", tt.input, res)
		}
	}
}

func TestValidStringRead(t *testing.T) {
	var b bytes.Buffer
	for _, tt := range validStringTests {
		b.Reset()
		b.Write(tt.input)
		res, err := ReadNext(&b)
		if err != nil {
			t.Errorf("Reading value failed with input %v: %v", tt.input, err)
		}
		s := fmt.Sprint(res)
		if s != tt.expected {
			t.Errorf("Expected %v, got %v", tt.expected, res)
		}
	}
}

func TestInvalidRead(t *testing.T) {
	var b bytes.Buffer
	for _, tt := range invalidReadTests {
		b.Reset()
		b.Write(tt.input)
		_, err := ReadNext(&b)
		if err == nil {
			t.Errorf("For input %v: expected error %s but no error returned", tt.input, tt.expected)
		} else if err != tt.expected {
			t.Errorf("For input %v: expected error %s but received %s", tt.input, tt.expected, err)
		}
	}
}

func TestValidWrite(t *testing.T) {
	var b bytes.Buffer
	for _, tt := range validWriteTests {
		b.Reset()
		err := tt.f(&b)
		if err != nil {
			t.Errorf("%v: writing value failed: %s", tt.name, err)
		}
		res := b.Bytes()
		if !bytes.Equal(res, tt.expected) {
			t.Errorf("%v: expected %v but got %v", tt.name, tt.expected, res)
		}
	}
}

func TestInvalidWrite(t *testing.T) {
	for _, tt := range invalidWriteTests {
		err := tt.f()
		if err == nil {
			t.Errorf("%v: expected error %s but no error returned", tt.name, tt.expected)
		} else if err != tt.expected {
			t.Errorf("%v: expected error %s but received %s", tt.name, tt.expected, err)
		}
	}
}
