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

func WithGenClientCmdAutoGen() Option {
	return func(pbCtx *parse2.PbContext, req *vo.CodeFuncParams) {
		var genVo = vo.GenClientCMD{
			ProtoBase: vo.ProtoBase{
				ProtoName: pbCtx.PackageName,
			},
		}
		for _, node := range pbCtx.RpcList {
			genVo.Items = append(genVo.Items, &vo.GenClientCMDItem{
				ProtoName: pbCtx.PackageName,
				RpcName:   node.Rpc.Name,
			})
		}
		template, err := temp.GenCodeByTemplate(gtpl.GenClientCmd, &genVo)
		if err != nil {
			return
		}
		err = utils.CreateAndWriteFile(path.Join(req.OutputDir, "cmd_autogen.go"), template)
		if err != nil {
			return
		}
	}
}
