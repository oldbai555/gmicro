/**
 * @Author: zjj
 * @Date: 2025/5/6
 * @Desc:
**/

package option

import (
	"github.com/emicklei/proto"
	utils "gmicro/common"
	parse2 "gmicro/proto/tool/parse"
	"gmicro/proto/tool/temp"
	"gmicro/proto/tool/temp/gtpl"
	"gmicro/proto/tool/vo"
	"path"
	"strings"
)

func WithGenClientErrcodeAutoGen() Option {
	return func(pbCtx *parse2.PbContext, req *vo.CodeFuncParams) {
		var genVo = vo.GenClientErrCode{
			ProtoBase: vo.ProtoBase{
				ProtoName: pbCtx.PackageName,
			},
		}
		enum := pbCtx.GetEnum("ErrCode")
		if enum != nil {
			for _, node := range enum.Elements {
				if ele, ok := node.(*proto.EnumField); ok {
					var comment string
					if ele.InlineComment != nil {
						comment = strings.Join(ele.InlineComment.Lines, " ")
					}
					comment = strings.TrimSpace(comment)
					if comment == "" {
						comment = ele.Name
					}
					genVo.Items = append(genVo.Items, &vo.GenClientErrCodeItem{
						Name:    ele.Name,
						Value:   int32(ele.Integer),
						Comment: comment,
					})
				}
			}
		}
		template, err := temp.GenCodeByTemplate(gtpl.GenClientErrcode, &genVo)
		if err != nil {
			return
		}
		err = utils.CreateAndWriteFile(path.Join(req.OutputDir, "errcode_autogen.go"), template)
		if err != nil {
			return
		}
	}
}
