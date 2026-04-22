package rhombifer

import "github.com/racg0092/rhombifer/pkg/models"

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

// Returns rhombifer global configuration. It can be used to adjust default behaviors.
// It returns a pointer to the default configuration
func GetConfig() *Config {
	return config
}

// Sets global configurtion. Overwrites any default configuration set.
// If you just want to adjust just certain behaviors without overwritting the
// defaults use [GetConfig]
func SetConfig(c *Config) {
	config = c
}

// A pointer to the found flags for the current command execution
var (
	ff         *[]*models.Flag // TODO: to be decom
	foundflags []*models.Flag
)
