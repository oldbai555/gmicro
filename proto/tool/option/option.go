/**
 * @Author: zjj
 * @Date: 2024/12/11
 * @Desc:
**/

package option

import (
	"gmicro/proto/tool/parse"
	"gmicro/proto/tool/vo"
)

type Option func(pbCtx *parse.PbContext, req *vo.CodeFuncParams)
