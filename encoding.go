package amqp

import (
	"encoding/binary"
	"io"
)

// ReadNext reads the next AMQP data element from the buffer.
// The format code is read first, and then any remaining
// reading is delegated based on type. Any error is passed
// through. ReadNext returns ErrInvalidCode for any
// unrecognized format code.
func ReadNext(r io.Reader) (interface{}, error) {
	code := make([]byte, 1, 1)
	_, err := r.Read(code)
	if err != nil {
		return nil, err
	}
	switch code[0] {
	case 0x40: // Coded null value
		return nil, nil
	case 0x41: // Coded boolean true
		return true, nil
	case 0x42: // Coded boolean false
		return false, nil
	case 0x43: // Coded uint32 0
		return uint32(0), nil
	case 0x44: // Coded uint64 0
		return uint64(0), nil
	case 0x50:
		var v uint8
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	case 0x51:
		var v int8
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	case 0x52: // uint32 in a single byte
		var v uint8
		err := binary.Read(r, binary.BigEndian, &v)
		return uint32(v), err
	case 0x53: // uint64 in a single byte
		var v uint8
		err := binary.Read(r, binary.BigEndian, &v)
		return uint64(v), err
	case 0x54: // int32 in a single byte
		var v int8
		err := binary.Read(r, binary.BigEndian, &v)
		return int32(v), err
	case 0x55: // int64 in a single byte
		var v int8
		err := binary.Read(r, binary.BigEndian, &v)
		return int64(v), err
	case 0x56:
		return readBoolean(r)
	case 0x60:
		var v uint16
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	case 0x61:
		var v int16
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	case 0x70:
		var v uint32
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	case 0x71:
		var v int32
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	case 0x72:
		var v float32
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	case 0x80:
		var v uint64
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	case 0x81:
		var v int64
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	case 0x82:
		var v float64
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	case 0x98:
		var v UUID
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	}
	return nil, ErrInvalidCode
}
