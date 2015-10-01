package amqp

import "io"

// readBoolean reads a single byte, returning true for 0x01,
// false for 0x00, and returning an InvalidFormat error for any other value.
// Errors are passed through.
func readBoolean(r io.Reader) (bool, error) {
	c := make([]byte, 1, 1)
	_, err := r.Read(c)
	if err != nil {
		return false, err
	}
	if c[0] == 0x00 {
		return false, nil
	} else if c[0] == 0x01 {
		return true, nil
	} else {
		return false, ErrInvalidFormat
	}
}

// WriteNull writes a single byte representing a null value
// per the AMQP specification. Any error is passed through.
func WriteNull(r io.Writer) error {
	_, err := r.Write([]byte{0x40})
	return err
}

// WriteBoolean writes a boolean according to the AMQP spec.
// The compact form is used (boolean value embedded in the format code).
// Any error is passed through.
func WriteBoolean(r io.Writer, c bool) error {
	var err error
	if c {
		_, err = r.Write([]byte{0x41})
	} else {
		_, err = r.Write([]byte{0x42})
	}
	return err
}

// WriteUByte writes a uint8 together with its format code.
// Two total bytes are added. Errors are passed through.
func WriteUByte(r io.Writer, c uint8) error {
	_, err := r.Write([]byte{0x50, c})
	return err
}
