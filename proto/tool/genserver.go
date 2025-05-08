/**
 * @Author: zjj
 * @Date: 2024/12/11
 * @Desc:
**/

package main

import (
	"github.com/spf13/cobra"
	utils "gmicro/common"
	"gmicro/pkg/exec"
	"gmicro/pkg/log"
	"gmicro/proto/tool/option"
	parse2 "gmicro/proto/tool/parse"
	"gmicro/proto/tool/vo"
	"path"
)

var genServerOps = []option.Option{
	option.WithGenServerCmdMainAutoGen(),
	option.WithGenServerImplStateAutoGen(),
	option.WithGenServerImplRpcAutoGen(),
	option.WithGenServerCmdApplicationAutoGen(),
	option.WithGenServerCmdCmdAutoGen(),
	option.WithGenServerGormGen(),
}

var CmdByGenServer = &cobra.Command{
	Use:   "genServer",
	Short: "生成服务端代码",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("intPutPbRoot: %s,pbNameList: %v,output: %s", pbDir, pbNameList, outputDir)
		for _, pbName := range pbNameList {
			pbCtx, err := parse2.ParsePb(path.Join(pbDir, parse2.GetPbNameWithSuffix(pbName)))
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
			var req = &vo.CodeFuncParams{
				PbName:    pbName,
				PbDir:     pbDir,
				OutputDir: outputDir,
				GitPath:   gitPath,
			}
			serverPath := path.Join(outputDir, pbCtx.ServiceName+"server")
			if !parse2.FileExists(serverPath) {
				utils.CreateDir(serverPath)
			}
			cmdPath := path.Join(serverPath, "cmd")
			if !parse2.FileExists(cmdPath) {
				utils.CreateDir(cmdPath)
			}
			implPath := path.Join(serverPath, "impl")
			if !parse2.FileExists(implPath) {
				utils.CreateDir(implPath)
			}
			for _, op := range genServerOps {
				op(pbCtx, req)
			}
		}
		// fmt
		exec.GoFmt(outputDir)
	},
}
