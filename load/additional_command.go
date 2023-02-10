package load

import (
	"errors"
	"github.com/daijulong/dockser/v2/lib"
	"io"
	"os"
	"reflect"
	"strings"
)

// AdditionalInstructionInterface 附加指令接口
type AdditionalInstructionInterface interface {
	// Load 解析、加载
	Load(args interface{}) error
	// Apply 执行
	Apply() error
}

//------------------------------------------------
//  附加指令
//------------------------------------------------

// Instructions 附加指令集
type Instructions struct {
	Instructions []AdditionalInstructionInterface
}

// NewInstructions Instructions constructor
func NewInstructions() *Instructions {
	return &Instructions{}
}

// 向指令集中追加一个指令
func (i *Instructions) append(instruction AdditionalInstructionInterface) {
	i.Instructions = append(i.Instructions, instruction)
}

// Load 加载指令集
func (i *Instructions) Load(name string, args interface{}) ([]AdditionalInstructionInterface, error) {
	//处理支持的附加指令列表
	switch name {
	case "copy":
		files, ok := args.([]interface{})
		if !ok {
			return nil, errors.New("load additional command[" + name + "] fail: invalid arguments")
		}
		for _, file := range files {
			if reflect.TypeOf(file).String() != "string" {
				continue
			}
			var cmd = NewAdditionalCopyInstruction()
			err := cmd.Load(file)
			if err != nil {
				return nil, err
			}
			i.append(&cmd)
		}
	case "remove":
		files, ok := args.([]interface{})
		if !ok {
			return nil, errors.New("load additional command[" + name + "] fail: invalid arguments")
		}
		for _, file := range files {
			if reflect.TypeOf(file).String() != "string" {
				continue
			}
			var cmd = NewAdditionalRemoveInstruction()
			err := cmd.Load(file)
			if err != nil {
				return nil, err
			}
			i.append(cmd)
		}
	default:
		return nil, errors.New("don't supported addition instruction [" + name + "]")
	}
	return i.Instructions, nil
}

//------------------------------------------------
//  附加指令：复制文件
//------------------------------------------------

// AdditionalCopyInstruction 附加指令：复制文件
type AdditionalCopyInstruction struct {
	Source      string
	Destination string
	Override    bool
}

// NewAdditionalCopyInstruction AdditionalCopyInstruction constructor
func NewAdditionalCopyInstruction() AdditionalCopyInstruction {
	return AdditionalCopyInstruction{}
}

// Load 解析、加载复制文件附加指令
func (ac *AdditionalCopyInstruction) Load(args interface{}) error {
	file, ok := args.(string)
	if !ok {
		return errors.New("load additional command[copy] fail: invalid arguments")
	}
	opts := strings.Split(file, ":")
	optsTotal := len(opts)
	if optsTotal < 2 {
		return errors.New("load additional command[copy] fail: invalid arguments")
	}
	src := strings.Trim(opts[0], " ")
	dst := strings.Trim(opts[1], " ")
	if src == "" || dst == "" {
		return errors.New("load additional command[copy] fail: invalid arguments")
	}

	ac.Source = src
	ac.Destination = dst
	if optsTotal >= 3 {
		ac.Override = opts[2] == "override"
	}

	return nil
}

// Apply 执行复制文件附加指令
func (ac *AdditionalCopyInstruction) Apply() error {
	lib.Warn("copy file ["+ac.Source+"] to ["+ac.Destination+"], override: ", ac.Override)
	//检查源文件
	if ok := lib.FileExist(ac.Source); !ok {
		return errors.New("source file not exist")
	}
	//目标文件是否存在，如果存在且不覆盖，直接结束
	dstExist := lib.FileExist(ac.Destination)
	if dstExist && ac.Override == false {
		return nil
	}
	//复制源文件到目标位置
	source, err := os.Open(ac.Source)
	if err != nil {
		return err
	}
	defer source.Close()
	destination, err := os.Create(ac.Destination)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)

	return err
}

//------------------------------------------------
//  附加指令：移除文件
//------------------------------------------------

// AdditionalRemoveInstruction 附加指令：移除文件
type AdditionalRemoveInstruction struct {
	File string
}

// NewAdditionalRemoveInstruction AdditionalRemoveInstruction constructor
func NewAdditionalRemoveInstruction() *AdditionalRemoveInstruction {
	return &AdditionalRemoveInstruction{}
}

// Load 解析、加载移除文件附加指令
func (ac *AdditionalRemoveInstruction) Load(args interface{}) error {
	file, ok := args.(string)
	if !ok {
		return errors.New("load additional command[remove] fail: invalid arguments")
	}
	file = strings.Trim(file, " ")
	if file == "" {
		return errors.New("load additional command[remove] fail: invalid arguments")
	}
	ac.File = file
	return nil
}

// Apply 执行移除文件附加指令
func (ac *AdditionalRemoveInstruction) Apply() error {
	lib.Warn("delete file [" + ac.File + "]")
	//检查文件是否存在
	if ok := lib.FileExist(ac.File); !ok {
		return nil
	}
	err := os.Remove(ac.File)
	if err != nil {
		return err
	}

	return nil
}
