package exec

import (
	"bytes"
	"fmt"
	"gmicro/pkg/log"
	"io/fs"
	"os"
	"os/exec"
	path2 "path"
	"path/filepath"
	"strings"
	"sync"
)

// ProtocGo 生成静态文件
func ProtocGo(goOut, pbDir, pbName string) {
	var protocTempFile = "protoc"
	pbName = strings.TrimSuffix(pbName, ".proto")
	goOutPath := filepath.ToSlash(goOut)
	// 生成pb文件
	var args []string
	args = append(args, fmt.Sprintf("--go_out=%s", goOutPath))
	args = append(args, fmt.Sprintf("--validate_out=lang=go:%s", goOutPath))
	args = append(args, fmt.Sprintf("--proto_path=%s", filepath.ToSlash(pbDir)))

	protoDir, protoName := filepath.Split(path2.Join(pbDir, pbName+".proto"))
	args = append(args, fmt.Sprintf("-I=%s", protoDir), protoName)

	protocPath := filepath.ToSlash(protocTempFile)
	cmd := exec.Command(protocPath, args...)

	log.Infof("exec: %s", cmd.String())

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	outStr := outBuf.String()
	errStr := errBuf.String()
	if err != nil {
		log.Infof("Err: %s \nStdout: %s \nStderr: %s", err, outStr, errStr)
		return
	}
	if outStr != "" {
		log.Infof("out:%v", outStr)
	}
	if errStr != "" {
		log.Errorf("err:%v", errStr)
	}
}

// ProtocGoTag 去除标签
func ProtocGoTag(input, pbName string) {
	pbName = strings.TrimSuffix(pbName, ".proto")

	// 生成pb文件
	var args []string
	args = append(args, fmt.Sprintf("-input=%s", filepath.ToSlash(input)))
	cmd := exec.Command("protoc-go-inject-tag", args...)

	log.Infof("exec: %s", cmd.String())

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	outStr := outBuf.String()
	errStr := errBuf.String()
	if err != nil {
		log.Infof("Err: %s \nStdout: %s \nStderr: %s", err, outStr, errStr)
		return
	}
	if outStr != "" {
		log.Infof("out: %v", outStr)
	}
	if errStr != "" {
		log.Errorf("err: %v", errStr)
	}
}
func GoFmt(path string) {
	var w sync.WaitGroup
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if d == nil {
			return nil
		}

		if !d.IsDir() {
			return nil
		}

		w.Add(1)
		go func(path string) {
			defer w.Done()
			cmd := exec.Command("gofmt", "-w", ".")
			cmd.Dir = path
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			log.Infof("run %s", strings.Join(cmd.Args, " "))
			log.Infof("run on %s", path)

			err = cmd.Run()
			if err != nil {
				log.Errorf("err:%v", err)
			}
		}(path)

		return nil
	})
	if err != nil {
		log.Errorf("err:%v", err)
	}

	w.Wait()
}
