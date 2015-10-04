/*
Package amqp implements the AMQP 1.0 specification.
*/
package amqp

import "errors"

// ErrInvalidCode indicates an unrecognized AMQP format code
var ErrInvalidCode = errors.New("Invalid Code")

// ErrInvalidFormat indicates that the retrieved data could not
// be converted to the expected data type
var ErrInvalidFormat = errors.New("Invalid Format")
