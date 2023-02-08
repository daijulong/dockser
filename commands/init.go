package commands

// CommandsMapping 子命令集
var CommandsMapping = NewCommands()

// 初始化子命令，注册子命令
func init() {
	CommandsMapping.Register("help", NewHelp())
	CommandsMapping.Register("make", NewMake())
	CommandsMapping.Register("init", NewInit())
}
