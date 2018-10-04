package display

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSON is a convenience function for printing JSON to STDOUT.
func JSON(data interface{}) (int, error) {
	msg, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	return fmt.Println(string(msg))
}

// Error prints to os.Stderr.
func Error(a ...interface{}) (int, error) {
	return fmt.Fprint(os.Stderr, a...)
}
