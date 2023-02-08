package commands

import (
	"fmt"
	"github.com/daijulong/dockser/core"
	"github.com/daijulong/dockser/lib"
	"github.com/daijulong/dockser/load"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

// Make 子命令 struct
type Make struct{}

// NewMake Make 子命令 constructor
func NewMake() *Make {
	return &Make{}
}

// Handle 执行命令
func (m *Make) Handle(args []string, options map[string]string) {
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
	groupFileBytes, err := os.ReadFile(groupFile)
	lib.IfErrorExit(err != nil, "read services group file ["+groupFile+"] failed: ", err)

	err = yaml.Unmarshal(groupFileBytes, &groups.Groups)
	lib.IfErrorExit(err != nil, "read services group file ["+groupFile+"] failed: ", err)
	group := groups.get(servicesGroupName) //组设置

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
	if group.isAutoOverride() && lib.IsFile(outputFile) {
		outputFile = lib.ForceFilenameSuffix(outputFile, now+".yml", "yml", "yaml")
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
	services := load.NewServices()
	lib.IfErrorExit(services.Load(group.Services) != nil, "load services fail: ", err)

	// 读取模板内容
	templateContent, err := lib.ReadFile(templateFile, templateFileName)
	lib.IfErrorExit(err != nil, "load template file ["+templateFile+"] failed: ", err)

	//解析模板
	templateContentBytes := []byte(templateContent)
	templateContentMap := make(map[string]interface{})
	if err = yaml.Unmarshal(templateContentBytes, &templateContentMap); err != nil {
		lib.ErrorExit("parse template content fail: ", err)
	}

	//压入 services
	templateServices := make(map[string]interface{})
	if sv, ok := templateContentMap["services"]; ok {
		if _, isService := sv.(map[interface{}]interface{}); isService { //只要符合 map[string]interface{} 即可
			for k, v := range sv.(map[interface{}]interface{}) {
				if _, isValidName := k.(string); isValidName {
					templateServices[k.(string)] = v
				} else {
					lib.ErrorExit("services content in template [" + templateFile + "] is invalid2")
				}
			}
		} else {
			lib.ErrorExit("services content in template [" + templateFile + "] is invalid")
		}
	} else {
		templateContentMap["services"] = make(map[string]interface{})
	}
	for _, service := range services.Services {
		if len(service.Services) > 0 {
			for k, v := range service.Services {
				templateServices[k] = v
			}
		}
		//执行添加/生成时的附加命令
		err = service.ApplyAddCommands()
		lib.IfErrorExit(err != nil, "service ["+service.Name+"] exec command fail: ", err)
	}
	templateContentMap["services"] = templateServices

	// 输出到文件
	outputBytes, err := yaml.Marshal(templateContentMap)
	lib.IfErrorExit(err != nil, "make docker-compose file ["+outputFile+"] fail: ", err)
	err = os.WriteFile(outputFile, outputBytes, 0755)
	lib.IfErrorExit(err != nil, "output to file ["+outputFile+"] failed: ", err)
	lib.Success("make docker-compose file [" + outputFile + "] success")
}

// Help 显示帮助信息
func (m *Make) Help() {
	doc := NewCommandHelpDocument()
	doc.Description = "make docker-compose.yml file with your group settings. "
	doc.Usage = "dockser " + lib.TextYellow("make") + " [" + lib.TextYellow("group") + "] [" + lib.TextYellow("options") + "] "
	doc.Options = append(doc.Options, map[string]string{"-o, --out, --output": "output file name"})
	doc.Options = append(doc.Options, map[string]string{"-t, --tpl, --template": "docker-compose.yml template"})
	doc.Args = append(doc.Args, map[string]string{"group": "group name in the settings, default is \"default\", view in \"groups.yml\" file"})
	doc.Print()
}

// 分组 struct
type group struct {
	Services []string `yaml:"services"`
	Template string   `yaml:"template"`
	Output   string   `yaml:"output"`
	Override string   `yaml:"override"`
}

// isAutoOverride 是否自动覆盖
func (g *group) isAutoOverride() bool {
	return g.Override != "force"
}

// 分组集 struct
type groups struct {
	Groups map[string]group
}

// groups 分组集 constructor
func newGroups() *groups {
	return &groups{}
}

// get 按名称获取分组
func (gs *groups) get(name string) group {
	if _, ok := gs.Groups[name]; !ok {
		lib.ErrorExit("services group [" + name + "] is undefined")
	}
	return gs.Groups[name]
}
