package errors

import (
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// General errors.
var (
	ErrNotImplemented = errors.New("not yet implemented")
)

type Type int

const (
	_ Type = iota
	UserInput
	Build
	VCS
	Developer
)

type Error struct {
	ExitCode        int
	Code            string
	Type            Type
	Causes          []error
	Message         string
	Troubleshooting string
}

func (e *Error) Error() string {
	return color.HiRedString("ERROR: ") + e.Message + "\n\n" + strings.TrimSpace(e.Troubleshooting) + `

` + color.WhiteString("REPORTING A BUG:") + `

Please try troubleshooting before filing a bug. If none of the suggestions work,
you can file a bug at <https://github.com/fossas/fossa-cli/issues/new>.

For additional support, ask the #cli channel at <https://slack.fossa.io>.

Before creating an issue, please search GitHub issues for similar problems. When
creating the issue, please attach the debug log located at:

    $log
`
}

// func Wrap(cause error, msg string) error {
// 	return errors.Wrap(cause, msg)
// }

// func Wrapf(cause error, format string, args ...interface{}) error {
// 	return errors.Wrapf(cause, format, args...)
// }

// func WrapError(cause error, err Error) Error {
// 	switch e := cause.(type) {
// 	case *Error:
// 		return Error{
// 			Cause:           e,
// 			Common:          err.Common,
// 			Message:         err.Message,
// 			Troubleshooting: err.Troubleshooting,
// 		}
// 	default:
// 	}
// 	return err
// }

// func New(msg string) error {
// 	return errors.New(msg)
// }
