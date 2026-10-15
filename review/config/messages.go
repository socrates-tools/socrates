package config

const (
	followsTheHost          = "follows the host a review begins from"
	unreadableFile          = "%s exists and socrates review cannot read it: %s. Fix its permissions or move it aside and run the command again"
	unparsableFile          = "%s is not configuration socrates review can read: %s. Fix the file or move it aside, because socrates review will not quietly run on defaults while a configuration file says otherwise"
	statesSomethingUnusable = "%s states something socrates review cannot use: %s. Fix the file or move it aside and run the command again"
	agorasBlockStated       = "%s states a %s block, and that block is socrates agora's now, read from %s. Move it there as it stands and run the command again"
	unknownRole             = "%q is not a role socrates review has, and its roles are %v"
	negativeMeasure         = "%s is negative, and it is a size or a duration"
)
