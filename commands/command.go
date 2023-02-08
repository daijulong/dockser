package commands

import (
	"bytes"
	"fmt"
	"github.com/daijulong/dockser/lib"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// CommandInterface 命令接口
type CommandInterface interface {
	// Handle 执行命令
	Handle(args []string, options map[string]string)
	// Help 显示帮助信息
	Help()
}

// Commands 命令集 struct
type Commands struct {
	commands map[string]CommandInterface
}

// NewCommands Commands constructor
func NewCommands() *Commands {
	return &Commands{commands: make(map[string]CommandInterface)}
}

// Get 获取一个命令
func (c *Commands) Get(command string) CommandInterface {
	if _, ok := c.commands[command]; ok {
		return c.commands[command]
	}
	return nil
}

// Register 注册一个命令
func (c *Commands) Register(name string, command CommandInterface) {
	c.commands[name] = command
}

// CommandHelpDocument 命令行帮助文档
type CommandHelpDocument struct {
	Description string
	Usage       string
	Args        []map[string]string
	Options     []map[string]string
	Commands    []map[string]string
	PrintMaxLen int
}

// NewCommandHelpDocument 命令行帮助文档 constructor
func NewCommandHelpDocument() *CommandHelpDocument {
	defaultOptions := make([]map[string]string, 0)
	defaultOptionHelp := make(map[string]string)
	defaultOptionHelp["-h, --help"] = "display help information"
	defaultOptions = append(defaultOptions, defaultOptionHelp)
	return &CommandHelpDocument{
		Description: "",
		Usage:       "",
		Args:        make([]map[string]string, 0),
		Options:     defaultOptions,
		Commands:    make([]map[string]string, 0),
		PrintMaxLen: 20,
	}
}

// Print 打印帮助文档，并自动对齐
func (h *CommandHelpDocument) Print() {
	if h.Description != "" {
		fmt.Println(h.Description)
	}
	if h.Usage != "" {
		fmt.Println()
		fmt.Println(lib.TextGreen("Usage: "), "\n", "\n     "+h.Usage)
	}
	// 计算左侧最大长度
	maxLen := h.printMaxLen(h.Args)
	if maxLen > h.PrintMaxLen {
		h.PrintMaxLen = maxLen
	}
	maxLen = h.printMaxLen(h.Options)
	if maxLen > h.PrintMaxLen {
		h.PrintMaxLen = maxLen
	}
	maxLen = h.printMaxLen(h.Commands)
	if maxLen > h.PrintMaxLen {
		h.PrintMaxLen = maxLen
	}

	h.PrintLines(h.Args, "Args")
	h.PrintLines(h.Options, "Options")
	h.PrintLines(h.Commands, "Commands")
}

// PrintLines 按行打印
func (h *CommandHelpDocument) PrintLines(lines []map[string]string, title string) {
	if len(lines) > 0 {
		fmt.Println()
		fmt.Println(lib.TextGreen(title, ": "), "\n")
		for _, items := range lines {
			if len(items) > 0 {
				for name, desc := range items {
					fmt.Printf("    %-"+strconv.Itoa(h.PrintMaxLen+4)+"s %s\n", name, desc)
				}
			}
		}
	}
}

// printMaxLen 计算一组字符串中最大长度
func (h *CommandHelpDocument) printMaxLen(lines []map[string]string) int {
	maxLen := 0
	if len(lines) > 0 {
		for _, items := range lines {
			if len(items) > 0 {
				for name, _ := range items {
					nameLen := len(name)
					if nameLen > maxLen {
						maxLen = nameLen
					}
				}
			}
		}
	}
	return maxLen
}

// Run 运行子命令
func Run() {
	args := os.Args
	command := "help"
	argsTotal := len(args)

	params := make([]string, 0)
	options := make(map[string]string)
	if argsTotal > 1 {
		for _, arg := range args[1:] {
			//以 - 开头的视为选项，否则为参数
			if strings.Index(arg, "-") == 0 {
				option := strings.Split(arg, "=")
				optionKey := strings.TrimLeft(option[0], "-")
				optionVal := ""
				if len(option) > 1 {
					optionVal = option[len(option)-1]
				}
				options[optionKey] = optionVal
			} else {
				params = append(params, arg)
			}
		}
	}
	if len(params) > 0 {
		command = params[0]
		params = params[1:]
	}
	commandHandler := CommandsMapping.Get(command)
	// 如果子命令不存在，则转发执行 docker-compsoe 命令
	// 目前只实现了同步执行的命令，交互式命令（如 docker-compose exec xxx sh）暂未实现
	if commandHandler == nil {
		warnTip := "*  no dockser sub command, will be exec: docker-compose " + strings.Join(args[1:], " ") + "  *"
		lib.Warn(strings.Repeat("*", len(warnTip)))
		lib.Warn(warnTip)
		lib.Warn(strings.Repeat("*", len(warnTip)))
		var stdout bytes.Buffer
		cmd := exec.Command("docker-compose", args[1:]...)
		cmd.Stdout = &stdout
		_ = cmd.Run()
		lib.Info(cmd.Stdout)
	} else {
		helpOption, _ := lib.GetOption(options, "help", "h")
		if helpOption == true {
			commandHandler.Help()
		} else {
			commandHandler.Handle(params, options)
		}
	}
}
