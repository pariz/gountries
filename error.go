package gountries

import (
	"fmt"
)

// Error returns a formatted error
func makeError(errMsg, errType string) error {
	return fmt.Errorf("gountries error. %s: %s", errMsg, errType)
}
