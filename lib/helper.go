package lib

import (
	"fmt"
	"github.com/daijulong/dockser/v2/core"
	"github.com/gookit/color"
	"os"
	"strings"
)

// ReadFile 读取文件，并替换环境变量
func ReadFile(file string, fileTitle string) (string, error) {
	if !FileExist(file) {
		return "", fmt.Errorf("%s is not exist", fileTitle)
	}

	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	fileContent := string(fileBytes)
	for env, value := range core.Envs {
		fileContent = strings.ReplaceAll(fileContent, "#@_"+env+"_@#", value)
	}
	return fileContent, nil
}

// FileExist 文件是否存在
func FileExist(file string) bool {
	fileStat, e := os.Stat(file)
	if e != nil {
		return false
	}
	return !fileStat.IsDir()
}

// GetOption 按顺序检查选项并获取选项值
func GetOption(options map[string]string, getOptions ...string) (exist bool, value string) {
	for _, option := range getOptions {
		if _, ok := options[option]; ok {
			return true, options[option]
		}
	}
	return false, ""
}

// GetOptionWithDefault 按顺序检查选项并获取选项值，如无则取默认值
func GetOptionWithDefault(options map[string]string, defaultValue string, setDefaultValueIfEmpty bool, getOptions ...string) string {
	exist, value := GetOption(options, getOptions...)
	if !exist || (value == "" && setDefaultValueIfEmpty == true) {
		return defaultValue
	}
	return value
}

// AutoFilenameSuffix 自动检测并添加文件名后缀，不区分大小写
func AutoFilenameSuffix(filename string, defaultSuffix string, suffixes ...string) string {
	filenameLen := len(filename)
	for _, suffix := range suffixes {
		suffixLen := len(strings.TrimLeft(suffix, ".")) + 1
		if suffixLen < filenameLen && strings.ToLower(filename[filenameLen-suffixLen:]) == "."+strings.ToLower(suffix) {
			return filename
		}
	}
	return filename + "." + strings.TrimLeft(defaultSuffix, ".")
}

// ForceFilenameSuffix 自动检测并强制添加文件名后缀，不区分大小写
func ForceFilenameSuffix(filename string, suffix string, suffixes ...string) string {
	filenameLen := len(filename)
	for _, _suffix := range suffixes {
		suffixLen := len(strings.TrimLeft(_suffix, ".")) + 1
		if suffixLen < filenameLen && strings.ToLower(filename[filenameLen-suffixLen:]) == "."+strings.ToLower(_suffix) {
			filename = filename[0 : filenameLen-suffixLen]
		}
	}
	return filename + "." + strings.TrimLeft(suffix, ".")
}

// FilePath 拼接并格式化文件路径
func FilePath(path ...string) string {
	pathList := make([]string, 0)
	sep := string(os.PathSeparator)
	for _, p := range path {
		_path := strings.Trim(p, "\\/")
		_path = strings.ReplaceAll(_path, "\\", sep)
		_path = strings.ReplaceAll(_path, "/", sep)
		pathList = append(pathList, _path)
	}
	return strings.Join(pathList, sep)
}

// IsDir 是否是目录
func IsDir(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil || os.IsExist(err)
}

// IsFile 是否是文件
func IsFile(dir string) bool {
	fi, e := os.Stat(dir)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}

// ErrorExit 输出错误信息并退出
func ErrorExit(a ...interface{}) {
	color.Error.Println(a...)
	os.Exit(0)
}

// IfErrorExit 如果出错则输出错误信息并退出
func IfErrorExit(expression bool, a ...interface{}) {
	if expression {
		ErrorExit(a...)
	}
}

// Error 输出错误信息
func Error(a ...interface{}) {
	color.Error.Println(a...)
}

// IfError 如果出错则输出错误信息
func IfError(expression bool, a ...interface{}) {
	if expression {
		Error(a...)
	}
}

// Success 输出成功信息
func Success(a ...interface{}) {
	color.Success.Println(a...)
}

// Warn 输出警告信息
func Warn(a ...interface{}) {
	color.Warn.Println(a...)
}

// Info 输出 info 信息
func Info(a ...interface{}) {
	color.Info.Println(a...)
}

// TextYellow 黄色文本
func TextYellow(a ...interface{}) string {
	return color.Yellow.Sprint(a...)
}

// TextGreen 黄色文本
func TextGreen(a ...interface{}) string {
	return color.Green.Sprint(a...)
}

// TextRed 红文本
func TextRed(a ...interface{}) string {
	return color.Red.Sprint(a...)
}
