package evaluator

import (
	"fmt"
	"monkey/object"
)

func newError(format string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}

// func newUnknownOperatorError(operator string, right)
