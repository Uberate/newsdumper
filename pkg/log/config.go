package log

// Config about the config of log for app.
type Config struct {
	Level string `json:"level" yaml:"level"`

	// DisableColor will disable color of output. Default false.
	// If not disable color, will auto check tty, if colorable, output will with color.
	DisableColor bool `json:"disable-color" yaml:"disable-color"`

	// Override coloring based on CLICOLOR and CLICOLOR_FORCE. - https://bixense.com/clicolors/
	EnvironmentOverrideColors bool `json:"environment-override-colors,omitempty" yaml:"environment-override-colors"`

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool `json:"disable-timestamp,omitempty" yaml:"disable-timestamp"`

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool `json:"full-timestamp,omitempty" yaml:"full-timestamp"`

	// TimestampFormat to use for display when a full timestamp is printed.
	// The format to use is the same than for time.Format or time.Parse from the standard
	// library.
	// The standard Library already provides a set of predefined format.
	TimestampFormat string `json:"timestamp-format,omitempty" yaml:"timestamp-format"`
}
