package commands

var CommandsMapping = NewCommands()

func init() {
	CommandsMapping.Register("help", NewHelp())
	CommandsMapping.Register("make", NewMake())
	CommandsMapping.Register("init", NewInit())
}
