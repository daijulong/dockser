package commands

import (
	"github.com/daijulong/dockser/core"
	"github.com/daijulong/dockser/lib"
	"github.com/gookit/color"
)

type Help struct{}

func NewHelp() *Help {
	return &Help{}
}

func (this *Help) Handle(args []string, options map[string]string) {
	versionOption, _ := lib.GetOption(options, "version", "v")
	if versionOption == true {
		doc := NewCommandHelpDocument()
		doc.Description = core.DOCKSER_NAME + " version: " + lib.TextYellow(core.DOCKSER_VERSION)
		doc.Options = make([]map[string]string, 0)
		doc.Print()
	} else {
		this.Help()
	}
}

func (this *Help) Help() {
	doc := NewCommandHelpDocument()
	doc.Description = `     ________                  ___
    /  ___   \                /  / ___ 
   /  /   |  /_____  ______  /  /_/  /_____________  ______
  /  /   /  /  __  \/  ____\/     __/ ______/ ___  \/  ___/
 /  /___/  /  /__/ /  /__ _   /\  \/_____  \  _____/  /
/_________/\______/\______/__/  \_/________/\_____/__/
` + "\nManage your docker-compose.yml more flexibly. version: " + lib.TextYellow(core.DOCKSER_VERSION)
	doc.Usage = core.DOCKSER_NAME + " " + lib.TextYellow("command") + " [" + lib.TextYellow("options") + "] [" + color.Yellow.Sprint("arguments") + "]"
	doc.Options = append(doc.Options, map[string]string{"-v, --version": "display version"})
	doc.Args = append(doc.Args, map[string]string{"command": "sub command"})
	doc.Commands = append(doc.Commands, map[string]string{"make": "make docker-compose.yml file with your group configs"})
	doc.Commands = append(doc.Commands, map[string]string{"init": "init your docker-compose project"})
	doc.Print()

	// 联系方式
	contacts := make([]map[string]string, 0)
	contacts = append(contacts, map[string]string{"Email: ": "daijulong@qq.com"})
	contacts = append(contacts, map[string]string{"Wechat: ": "julongdai"})
	contacts = append(contacts, map[string]string{"QQ: ": "88622090"})
	doc.PrintLines(contacts, "Contact me")
}
