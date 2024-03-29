package commands

import (
	"github.com/daijulong/dockser/v2/core"
	"github.com/daijulong/dockser/v2/lib"
	"github.com/daijulong/dockser/v2/resource"
	"os"
)

// Init 子命令 struct
type Init struct{}

// NewInit Init 子命令 constructor
func NewInit() *Init {
	return &Init{}
}

// Handle 执行命令
func (i *Init) Handle(args []string, options map[string]string) {
	// 生成 .env, .env_example 文件
	// 生成 compose 目录和 services, templates 子目录
	// 生成 compose/groups.yml 文件
	// 如果生成演示文件，则：
	//    groups.yml 中多一组演示配置
	//    services 目录下多 nginx
	//    .env 文件增加相关内容

	// 获取初始化项目目录，优先级：option > default
	initDir := lib.GetOptionWithDefault(options, core.DOCKSER_INIT_DIR, true, "dir", "d")
	initDir = lib.FilePath(initDir)
	// 是否生成演示数据
	withDemo, _ := lib.GetOption(options, "with-demo")

	// 需要创建的各目录
	dockerComposeDir := lib.FilePath(initDir, "compose")
	dockerComposeServicesDir := lib.FilePath(dockerComposeDir, "services")
	dockerComposeTemplatesDir := lib.FilePath(dockerComposeDir, "templates")
	// 需要写入的各文件路径
	envFile := lib.FilePath(initDir, ".env")
	envExampleFile := lib.FilePath(initDir, ".env.example")
	groupFile := lib.FilePath(dockerComposeDir, "groups.yml")
	defaultTemplateFile := lib.FilePath(dockerComposeTemplatesDir, "docker-compose.yml")
	// demo 需要的各文件
	demoTemplateFile := lib.FilePath(dockerComposeTemplatesDir, "docker-compose-demo.yml")
	// 服务
	serviceNginxFile := lib.FilePath(dockerComposeServicesDir, "nginx.yml")

	// 检查并创建所需目录
	dirs := make([][]string, 0)
	dirs = append(dirs,
		[]string{initDir, "project"},
		[]string{dockerComposeDir, "compose"},
		[]string{dockerComposeServicesDir, "services"},
		[]string{dockerComposeTemplatesDir, "templates"},
	)
	for _, dir := range dirs {
		if lib.IsDir(dir[0]) {
			lib.Info(dir[1], "dir [", dir[0], "] already exists.")
		} else {
			err := os.Mkdir(dir[0], 0755)
			// 创建失败退出
			lib.IfErrorExit(err != nil, "create ", dir[1], " dir [", dir[0], "] fail: ", err)
			// 创建成功提示
			lib.Info("create", dir[1], "dir ["+dir[0]+"] success.")
		}
	}

	// 写入文件内容
	InitFiles := newInitFiles()
	envFileContent := resource.InitFileEnvContent
	groupFileCotent := resource.InitFileGroupContent
	if withDemo {
		envFileContent = resource.InitFileEnvDemoContent
		groupFileCotent = resource.InitFileGroupDemoContent
	}
	InitFiles.add(newInitFile(".env", envFile, envFileContent))
	InitFiles.add(newInitFile(".env.example", envExampleFile, resource.InitFileEvnExampleContent))
	InitFiles.add(newInitFile("groups.yml", groupFile, groupFileCotent))
	InitFiles.add(newInitFile("template", defaultTemplateFile, resource.InitFileTemplateContent))
	if withDemo {
		InitFiles.add(newInitFile("demo template", demoTemplateFile, resource.InitFileTemplateDemoContent))
	}
	InitFiles.add(newInitFile("service:nginx", serviceNginxFile, resource.InitFileServiceNginxContent))
	// 检查文件是否存在，如果存在则不写入
	for _, file := range InitFiles.Files {
		if lib.FileExist(file.File) {
			lib.Warn("init file [" + file.Title + "] warning: file already exists and nothing will be written.")
			continue
		}
		outputBytes := []byte(file.Content)
		err := os.WriteFile(file.File, outputBytes, 0755)
		lib.IfErrorExit(err != nil, "init file ["+file.Title+"] failed: ", err)
		lib.Info("init file [" + file.Title + "] success")
	}
	lib.Success("init success. please open the directory [", initDir, "] to view.")
}

// Help 显示帮助信息
func (i *Init) Help() {
	doc := NewCommandHelpDocument()
	doc.Description = "init your docker-compose project."
	doc.Usage = "dockser " + lib.TextYellow("init") + " [" + lib.TextYellow("options") + "] "
	doc.Options = append(doc.Options, map[string]string{"-d, --dir": "your project directory, default is the current directory"})
	doc.Options = append(doc.Options, map[string]string{"--with-demo": "init with demo data"})
	doc.Print()
}

// 初始化文件 struct
type initFile struct {
	Title   string
	File    string
	Content string
}

// initFile constructor
func newInitFile(title string, file string, content string) *initFile {
	return &initFile{Title: title, File: file, Content: content}
}

// 初始化文件集 struct
type initFiles struct {
	Files []*initFile
}

// initFiles constructor
func newInitFiles() *initFiles {
	return &initFiles{Files: make([]*initFile, 0)}
}

// add 向文件集中添加一个文件
func (f *initFiles) add(file *initFile) {
	f.Files = append(f.Files, file)
}
