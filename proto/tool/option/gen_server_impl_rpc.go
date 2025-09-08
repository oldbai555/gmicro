/**
 * @Author: zjj
 * @Date: 2025/5/6
 * @Desc:
**/

package option

import (
	"fmt"
	utils "gmicro/common"
	"gmicro/pkg/log"
	parse2 "gmicro/proto/tool/parse"
	"gmicro/proto/tool/temp"
	"gmicro/proto/tool/temp/gtpl"
	"gmicro/proto/tool/vo"
	"os"
	"path"
	"strings"
)

func WithGenServerImplRpcAutoGen() Option {
	return func(pbCtx *parse2.PbContext, req *vo.CodeFuncParams) {
		for _, rpcNode := range pbCtx.RpcList {
			var defaultFileName = "impl.go"
			specFileName := strings.ToLower(rpcNode.Options["(base.CodeGenRpcFuncFileName)"])
			if specFileName != "" {
				defaultFileName = specFileName + ".go"
			}
			var absGoFilePath = path.Join(req.OutputDir, pbCtx.ServiceName+"server", "impl", defaultFileName)
			// 不存在就创建
			if !parse2.FileExists(absGoFilePath) {
				var strs strings.Builder
				strs.WriteString(fmt.Sprintf(`
package impl

import (
	"gmicro/pkg/uctx"
	"gmicro/service/%s"
)
`, pbCtx.ServiceName))
				err := utils.CreateAndWriteFile(absGoFilePath, strs.String())
				if err != nil {
					log.Errorf("err:%v", err)
					return
				}
			}

			goFile, err := parse2.ParseGoFile(absGoFilePath)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}

			var content string
			funcSvr := utils.Slice2MapKeyByStructField(goFile.Funcs, "Name").(map[string]*parse2.GoFunc)

			if rpcNode.IgnoreSvrRpc {
				continue
			}

			var defaultTemp = gtpl.GenRpcServerFuncCode
			//switch strings.ToUpper(rpcNode.Options["(base.GenCRUDSvrRpcTemp)"]) {
			//case strings.ToUpper("ADD"):
			//	defaultTemp = gtpl.GenRpcServerFuncCodeByAdd
			//case strings.ToUpper("GET"):
			//	defaultTemp = gtpl.GenRpcServerFuncCodeByGet
			//case strings.ToUpper("UPDATE"):
			//	defaultTemp = gtpl.GenRpcServerFuncCodeByUpdate
			//case strings.ToUpper("DELETE"):
			//	defaultTemp = gtpl.GenRpcServerFuncCodeByDelete
			//case strings.ToUpper("LIST"):
			//	defaultTemp = gtpl.GenRpcServerFuncCodeByList
			//}
			if _, ok := funcSvr[rpcNode.Rpc.Name]; !ok {
				content, err = temp.GenCodeByTemplate(defaultTemp, &vo.GenServerImplRpc{
					RpcName:   rpcNode.Rpc.Name,
					RpcReq:    rpcNode.Rpc.RequestType,
					RpcRsp:    rpcNode.Rpc.ReturnsType,
					Server:    pbCtx.ServiceName,
					NewSev:    utils.UpperFirst(pbCtx.ServiceName),
					ModelName: utils.UpperFirst(rpcNode.Options["(base.ModelName)"]),
					Client:    pbCtx.ServiceName,
				})
				if err != nil {
					log.Errorf("err is %v", err)
					return
				}
			}
			log.Infof("content is %s to %s", content, defaultFileName)

			// 打开文件
			f, err := os.OpenFile(absGoFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				log.Errorf("can not generate file %s,Error :%v", absGoFilePath, err)
				return
			}

			// 追加文件内容
			err = parse2.WriteFile(f, content)
			if err != nil {
				log.Errorf("err is %v", err)
				return
			}

			err = f.Close()
			if err != nil {
				log.Errorf("can not close file %s,Error :%v", absGoFilePath, err)
				return
			}
		}
	}
}
