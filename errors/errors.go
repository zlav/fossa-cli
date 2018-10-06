package errors

import (
	"strings"

	"github.com/fatih/color"
)

// Type classifies different errors.
type Type int

// A list of different error types. These are grouped by classes of failures
// that we commonly see. They're not used for any logic yet, but will be in
// logged during errors.
const (
	_ Type = iota
	UserInput
	Build
	VCS
	System
	Developer
)

// Error is our domain-specific error implementation for providing user-friendly
// error messages.
type Error struct {
	ExitCode        int
	Code            string
	Type            Type
	Cause           error
	Message         string
	Troubleshooting string
}

// Error implements interface error, but renders _unfinished_ error messages
// that may still have format verbs within them.
//
// See Render for a list of supported verbs.
func (e *Error) Error() string {
	instructions := "| " + strings.Join(strings.Split(e.Troubleshooting, "\n"), "\n| ")
	return color.HiRedString("ERROR:") + " " + e.Message + "\n" + instructions + `

` + color.HiWhiteString("REPORTING A BUG:") + `

Please try troubleshooting before filing a bug. If none of the suggestions work,
you can file a bug at ` + color.HiBlueString("https://github.com/fossas/fossa-cli/issues/new") + `.

For additional support, ask the ` + color.MagentaString("#cli") + ` channel at ` + color.HiBlueString("https://slack.fossa.io") + `.

Before creating an issue, please search GitHub issues for similar problems. When
creating the issue, please attach the debug log located at:

    $log
`
}

// MessageArgs are substituted within unfinished error messages.
type MessageArgs struct {
	Invocation string
	Endpoint   string
	LogFile    string
}

// Render renders _finished_ error messages with variables substituted for
// verbs.
//
// The supported verbs are:
//
//   - $$: a literal $
//   - $command: the command the binary was invoked with
//   - $log: the name of this run's debug log file
//
func Render(msg string, a MessageArgs) string {
	// Handle escaped "$"s.
	sections := strings.Split(msg, "$$")

	for i := range sections {
		// Handle substitutions.
		// TODO: ideally, we would also be able to _modify_ commands. For example:
		//
		//   $command.WithFlag("--output").WithoutFlag("--endpoint").Cmd("test")
		//
		sections[i] = strings.Replace(sections[i], "$command", a.Invocation, -1)
		sections[i] = strings.Replace(sections[i], "$endpoint", a.Endpoint, -1)
		sections[i] = strings.Replace(sections[i], "$log", a.LogFile, -1)
	}
	return strings.Join(sections, "$")
}
