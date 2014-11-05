package filebase

import "fmt"

type Error struct {
	Err      string
	Detailed string
}

func (e *Error) Fault(details ...interface{}) *Error {
	e.Detailed = fmt.Sprintf(e.Err, details...)
	return e
}

func (e Error) Error() string {
	if e.Detailed != "" {
		return e.Detailed
	}
	return e.Err
}
