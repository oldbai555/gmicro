package parse

import (
	utils "gmicro/common"
	"os"
)

// FileExists 文件是否存在
func FileExists(name string) bool {
	return utils.FileExists(name)
}

// WriteFile 写内容进入文件
func WriteFile(file *os.File, content string) error {
	return utils.WriteFile(file, content)
}
