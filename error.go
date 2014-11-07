package filebase

import "fmt"

// `fault, because `error` is  type.
type fault struct {
	Err      string
	Detailed string
}

func (e *fault) Fault(details ...interface{}) *fault {
	e.Detailed = fmt.Sprintf(e.Err, details...)
	return e
}

func (e fault) Error() string {
	if e.Detailed != "" {
		return e.Detailed
	}
	return e.Err
}
