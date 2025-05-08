/**
 * @Author: zjj
 * @Date: 2025/8/7
 * @Desc:
**/

package common

import "os"

// WriteFile 写内容进入文件
func WriteFile(file *os.File, content string) error {
	if _, err := file.Write([]byte(content)); err != nil {
		return err
	}
	return nil
}

// CreateAndWriteFile 写内容进入文件
func CreateAndWriteFile(absTargetFilePath, content string) error {
	f, err := os.Create(absTargetFilePath)
	if err != nil {
		return err
	}
	return WriteFile(f, content)
}
