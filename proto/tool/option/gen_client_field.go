/**
 * @Author: zjj
 * @Date: 2025/5/6
 * @Desc:
**/

package option

import (
	"github.com/emicklei/proto"
	utils "gmicro/common"
	"gmicro/pkg/pie"
	parse2 "gmicro/proto/tool/parse"
	"gmicro/proto/tool/temp"
	"gmicro/proto/tool/temp/gtpl"
	"gmicro/proto/tool/vo"
	"path"
	"sort"
)

func WithGenClientFieldAutoGen() Option {
	return func(pbCtx *parse2.PbContext, req *vo.CodeFuncParams) {
		fields := pie.Strings{}
		for _, message := range pbCtx.MsgList {
			// step 3.2: 拿到对应的字段
			for _, ele := range message.Elements {
				switch ele.(type) {
				case *proto.NormalField:
					fields = fields.Append(ele.(*proto.NormalField).Name)
				case *proto.MapField:
					fields = fields.Append(ele.(*proto.MapField).Name)
				}
			}
		}

		//  字段去重
		fields = fields.Unique()
		sort.Strings(fields)
		var genVo = vo.GenClientField{
			ProtoBase: vo.ProtoBase{
				ProtoName: pbCtx.PackageName,
			},
		}
		fields.Each(func(v string) {
			genVo.Items = append(genVo.Items, &vo.GenClientFieldItem{
				Name: utils.UnderScore2Camel(v),
				CV:   utils.UnderScore2Camel(v),
				UV:   utils.Camel2UnderScore(v),
			})
		})
		template, err := temp.GenCodeByTemplate(gtpl.GenClientField, &genVo)
		if err != nil {
			return
		}
		err = utils.CreateAndWriteFile(path.Join(req.OutputDir, "field_autogen.go"), template)
		if err != nil {
			return
		}
	}
}
