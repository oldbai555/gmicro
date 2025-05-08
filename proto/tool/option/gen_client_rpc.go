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

func WithGenClientRpcAutoGen() Option {
	return func(pbCtx *parse2.PbContext, req *vo.CodeFuncParams) {
		var genVo = vo.GenClientRpc{
			ProtoBase: vo.ProtoBase{
				ProtoName: pbCtx.PackageName,
			},
		}
		for _, node := range pbCtx.RpcList {
			genVo.Items = append(genVo.Items, &vo.GenClientRpcItem{
				RpcName:   node.Rpc.Name,
				ProtoName: pbCtx.PackageName,
			})
		}
		template, err := temp.GenCodeByTemplate(gtpl.GenClientRpc, &genVo)
		if err != nil {
			return
		}
		err = utils.CreateAndWriteFile(path.Join(req.OutputDir, "client.go"), template)
		if err != nil {
			return
		}
	}
}
