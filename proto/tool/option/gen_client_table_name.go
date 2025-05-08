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

func WithGenClientTableNameAutoGen() Option {
	return func(pbCtx *parse2.PbContext, req *vo.CodeFuncParams) {
		var genVo = vo.GenClientTableName{
			ProtoBase: vo.ProtoBase{
				ProtoName: pbCtx.PackageName,
			},
		}
		for _, message := range pbCtx.MsgList {
			if !parse2.IsModel(message) {
				continue
			}
			genVo.Items = append(genVo.Items, &vo.GenClientTableNameItem{
				MessageName: message.Name,
				TableName:   parse2.GenTableName(pbCtx.ServiceName, message.Name),
			})
		}
		template, err := temp.GenCodeByTemplate(gtpl.GenClientTableName, &genVo)
		if err != nil {
			return
		}
		err = utils.CreateAndWriteFile(path.Join(req.OutputDir, "table_name_autogen.go"), template)
		if err != nil {
			return
		}
	}
}
