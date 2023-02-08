package load

import (
	"errors"
	"github.com/daijulong/dockser/lib"
	"io"
	"os"
)

// AdditionalInstruction 附加指令接口
type AdditionalInstruction interface {
	Apply() error
}

// AdditionalCopyInstruction 附加指令：复制文件
type AdditionalCopyInstruction struct {
	Source      string
	Destination string
	Override    bool
}

// Apply 执行复制文件附加指令
func (ac *AdditionalCopyInstruction) Apply() error {
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
	destination.Close()
	_, err = io.Copy(destination, source)

	return err
}

// AdditionalRemoveInstruction 附加指令：移除文件
type AdditionalRemoveInstruction struct {
	File string
}

// Apply 执行移除文件附加指令
func (ac *AdditionalRemoveInstruction) Apply() error {

	return nil
}
