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

func WithGenServerImplStateAutoGen() Option {
	return func(pbCtx *parse2.PbContext, req *vo.CodeFuncParams) {
		serverPath := path.Join(req.OutputDir, pbCtx.ServiceName+"server")
		if !parse2.FileExists(serverPath) {
			utils.CreateDir(serverPath)
		}
		cmdPath := path.Join(serverPath, "impl")
		if !parse2.FileExists(cmdPath) {
			utils.CreateDir(cmdPath)
		}
		statePath := path.Join(cmdPath, "state.go")
		if parse2.FileExists(statePath) {
			return
		}
		var genVo = vo.ProtoBase{ProtoName: pbCtx.ServiceName}
		template, err := temp.GenCodeByTemplate(gtpl.GenServerImplState, &genVo)
		if err != nil {
			return
		}
		err = utils.CreateAndWriteFile(statePath, template)
		if err != nil {
			return
		}
	}
}
