/**
 * @Author: zjj
 * @Date: 2024/12/11
 * @Desc:
**/

package main

import (
	"github.com/spf13/cobra"
	"gmicro/pkg/exec"
	"gmicro/pkg/log"
	"gmicro/proto/tool/option"
	parse2 "gmicro/proto/tool/parse"
	"gmicro/proto/tool/vo"

	path2 "path"
)

var genClientOps = []option.Option{
	option.WithGenClientCmdAutoGen(),
	option.WithGenClientErrcodeAutoGen(),
	option.WithGenClientFieldAutoGen(),
	option.WithGenClientRpcAutoGen(),
	option.WithGenClientTableNameAutoGen(),
}

var CmdByGenClient = &cobra.Command{
	Use:   "genClient",
	Short: "生成客户端代码",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("intPutPbRoot: %s,pbNameList: %v,output: %s", pbDir, pbNameList, outputDir)
		for _, pbName := range pbNameList {
			pbCtx, err := parse2.ParsePb(path2.Join(pbDir, parse2.GetPbNameWithSuffix(pbName)))
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
			log.Infof("gen protoc go code")
			exec.ProtocGo(goOut, pbDir, pbName)
			var req = &vo.CodeFuncParams{
				PbName:    pbName,
				PbDir:     pbDir,
				OutputDir: outputDir,
				GitPath:   gitPath,
			}
			for _, op := range genClientOps {
				op(pbCtx, req)
			}
		}
		// fmt
		exec.GoFmt(outputDir)
	},
}
