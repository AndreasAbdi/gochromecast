package media

import (
	"fmt"
)

type InvalidArgumentError struct {
	arg    interface{}
	reason string
}

func (er *InvalidArgumentError) Error() string {
	return fmt.Sprintf("%v is not valid argument: %v", er.arg, er.reason)
}
