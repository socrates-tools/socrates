package state

const (
	noHomeDirectory      = "socrates keeps everything it reads and writes under one root directory, and it could not find your own home directory to put it under: %s. Set %s to an absolute path and run the command again"
	statedHomeIsRelative = "%s names %q, which is not an absolute path, and socrates would keep its root wherever the command happened to be run. Set it to an absolute path and run the command again"
)
