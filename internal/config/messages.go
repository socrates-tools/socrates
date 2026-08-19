package config

const (
	noConfigurationDirectory = "socrates reads its configuration from your configuration directory, and it could not find one: %s. Set XDG_CONFIG_HOME and run the command again"
	unreadableFile           = "%s exists and socrates cannot read it: %s. Fix its permissions or move it aside and run the command again"
	unparsableFile           = "%s is not configuration socrates can read: %s. Fix the file or move it aside, because socrates will not quietly run on defaults while a configuration file says otherwise"
	statesSomethingUnusable  = "%s states something socrates cannot use: %s. Fix the file or move it aside and run the command again"
	unknownRole              = "%q is not a role socrates has, and its roles are %v"
	negativeMeasure          = "%s is negative, and it is a size or a duration"
)
