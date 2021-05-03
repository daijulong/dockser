package commands

import (
	"fmt"
	"github.com/daijulong/dockser/core"
	"github.com/daijulong/dockser/lib"
	"github.com/daijulong/dockser/load"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"time"
)

type Make struct{}

func NewMake() *Make {
	return &Make{}
}

func (this *Make) Handle(args []string, options map[string]string) {
	// 默认分组，优先从 .env 文件中读取 DEFAULT_GROUP
	defaultGroup := "default"
	if _, ok := core.Envs["DEFAULT_GROUP"]; ok {
		defaultGroup = core.Envs["DEFAULT_GROUP"]
	}

	// 指定的 services 分组
	servicesGroupName := defaultGroup
	if len(args) < 1 {
		lib.Info("parameter [" + lib.TextYellow("group") + "] not found, the default group [" + lib.TextYellow(defaultGroup) + "] will be taken")
	} else {
		servicesGroupName = args[0]
	}

	// 组配置
	groupFile := "./compose/groups.yml"
	lib.IfErrorExit(!lib.FileExist(groupFile), "services group file ["+groupFile+"] does not exist")

	groups := newGroups()
	groupFileBytes, err := ioutil.ReadFile(groupFile)
	lib.IfErrorExit(err != nil, "read services group file ["+groupFile+"] failed: ", err)

	err = yaml.Unmarshal(groupFileBytes, &groups.Groups)
	lib.IfErrorExit(err != nil, "read services group file ["+groupFile+"] failed: ", err)
	group := groups.Get(servicesGroupName) //组设置

	// 输出文件名
	defaultOutputFile := ""
	if group.Output != "" {
		defaultOutputFile = group.Output
	}
	t := time.Now()
	now := fmt.Sprintf("%d%d%d%d%d%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	if defaultOutputFile == "" {
		defaultOutputFile = "docker-compose-" + now + ".yml"
	}
	defaultOutputFile = lib.AutoFilenameSuffix(defaultOutputFile, "yml", "yml", "yaml")
	outputFile := lib.GetOptionWithDefault(options, defaultOutputFile, true, "output", "out", "o")
	outputFile = lib.AutoFilenameSuffix(outputFile, "yml", "yml", "yaml")
	// 自动覆盖时，
	if group.IsAutoOverride() && lib.IsFile(outputFile) {
		outputFile = lib.ForceFilenameSuffix(outputFile, now + ".yml", "yml", "yaml")
	}

	// 使用的模板
	defaultTemplateName := ""
	if group.Template != "" {
		defaultTemplateName = group.Template
	}
	if defaultTemplateName == "" {
		defaultTemplateName = "docker-compose.yml"
	}
	templateFileName := lib.GetOptionWithDefault(options, defaultTemplateName, true, "template", "tpl", "t")
	templateFileName = lib.AutoFilenameSuffix(templateFileName, "yml", "yml", "yaml")
	templateFile := "./compose/templates/" + templateFileName

	lib.IfErrorExit(!lib.FileExist(templateFile), "template file ["+templateFile+"] does not exist")

	// 读取 services
	servicesContent, err := load.Services(group.Services)
	lib.IfErrorExit(err != nil, "load services failed: ", err)

	// 读取模板内容
	templateLines, err := lib.ReadFileLines(templateFile, templateFileName)
	lib.IfErrorExit(err != nil, "load template file ["+templateFile+"] failed: ", err)

	// 将模板中的 services 占位符替换成 services 内容
	for row, line := range templateLines {
		if strings.TrimSpace(line) == "@@_SERVICES_@@" {
			templateLines[row] = servicesContent
		}
	}
	// 输出到文件
	outputBytes := []byte(strings.Join(templateLines, "\n"))
	err = ioutil.WriteFile(outputFile, outputBytes, 0755)
	lib.IfErrorExit(err != nil, "output to file ["+outputFile+"] failed: ", err)
	lib.Success("make docker-compose file [" + outputFile + "] success")
}

func (this *Make) Help() {
	doc := NewCommandHelpDocument()
	doc.Description = "make docker-compose.yml file with your group settings. "
	doc.Usage = "dockposer " + lib.TextYellow("make") + " [" + lib.TextYellow("group") + "] [" + lib.TextYellow("options") + "] "
	doc.Options = append(doc.Options, map[string]string{"-o, --out, --output": "output file name"})
	doc.Options = append(doc.Options, map[string]string{"-t, --tpl, --template": "docker-compose.yml template"})
	doc.Args = append(doc.Args, map[string]string{"group": "group name in the settings, default is \"default\", view in \"groups.yml\" file"})
	doc.Print()
}

type group struct {
	Services []string `yaml:"services"`
	Template string   `yaml:"template"`
	Output   string   `yaml:"output"`
	Override string   `yaml:"override"`
}

func (this *group) IsAutoOverride() bool {
	return this.Override != "force"
}

type groups struct {
	Groups map[string]group
}

func newGroups() *groups {
	return &groups{}
}

func (this *groups) Get(name string) group {
	if _, ok := this.Groups[name]; !ok {
		lib.ErrorExit("services group [" + name + "] is undefined")
	}
	return this.Groups[name]
}
