package rhombifer

// Global rhombifer configuration
type Config struct {
	// Determines if the root command allows for flags. By Default is true
	AllowFlagsInRoot bool

	// Determines if the help command should be run if no subcommand or flags are found. Defaults to false
	RunHelpIfNoInput bool
}

var config = &Config{
	AllowFlagsInRoot: true,
}

// GetConfig returns rhombifer global configuration. It can be used to adjust default behaviors.
func GetConfig() *Config {
	return config
}

// SetConfig global configurtion. Overwrites any default configuration set.
// If you just want to adjust just certain behaviors without overwritting the
// defaults use [GetConfig]
func SetConfig(c *Config) {
	config = c
}

// A pointer to the found flags for the current command execution
var (
	foundflags []*Flag
)
