/**
 * @Author: zjj
 * @Date: 2025/5/6
 * @Desc:
**/

package option

import (
	utils "gmicro/common"
	parse2 "gmicro/proto/tool/parse"
	"gmicro/proto/tool/temp"
	"gmicro/proto/tool/temp/gtpl"
	"gmicro/proto/tool/vo"
	"path"
)

func WithGenServerCmdCmdAutoGen() Option {
	return func(pbCtx *parse2.PbContext, req *vo.CodeFuncParams) {
		serverPath := path.Join(req.OutputDir, pbCtx.ServiceName+"server")
		if !parse2.FileExists(serverPath) {
			utils.CreateDir(serverPath)
		}
		cmdPath := path.Join(serverPath, "cmd")
		if !parse2.FileExists(cmdPath) {
			utils.CreateDir(cmdPath)
		}
		mainPath := path.Join(cmdPath, "cmd_autogen.go")
		var genVo = vo.GenServerCmdCmdAutoGen{
			ProtoBase: vo.ProtoBase{
				ProtoName: pbCtx.ServiceName,
			},
		}
		for _, node := range pbCtx.RpcList {
			genVo.Items = append(genVo.Items, &vo.GenServerCmdCmdAutoGenItem{
				ProtoName: pbCtx.ServiceName,
				RpcName:   node.Rpc.Name,
			})
		}
		template, err := temp.GenCodeByTemplate(gtpl.GenServerCmdCmdAutogen, &genVo)
		if err != nil {
			return
		}
		err = utils.CreateAndWriteFile(mainPath, template)
		if err != nil {
			return
		}
	}
}
