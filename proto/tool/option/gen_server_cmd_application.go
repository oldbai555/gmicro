/**
 * @Author: zjj
 * @Date: 2025/5/6
 * @Desc:
**/

package option

import (
	"bufio"
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

func WithGenServerCmdApplicationAutoGen() Option {
	return func(pbCtx *parse2.PbContext, req *vo.CodeFuncParams) {
		serverPath := path.Join(req.OutputDir, pbCtx.ServiceName+"server")
		if !parse2.FileExists(serverPath) {
			utils.CreateDir(serverPath)
		}
		cmdPath := path.Join(serverPath, "cmd")
		if !parse2.FileExists(cmdPath) {
			utils.CreateDir(cmdPath)
		}
		yamlPath := path.Join(cmdPath, "application.yaml")
		if utils.FileExists(yamlPath) {
			return
		}
		var maxPort uint32 = 20000
		var step uint32 = 10000
		var minErrcode = 2 * step
		var maxErrcode = 3 * step
		registerPath := path.Join(req.PbDir, "register")
		var genVo = vo.GenServerCmdApplication{
			Ip:   "0.0.0.0",
			Name: pbCtx.ServiceName,
		}
		if !parse2.FileExists(registerPath) {
			err := utils.CreateAndWriteFile(registerPath, fmt.Sprintf("%s,%d,%d,%d\n", pbCtx.ServiceName, maxPort, minErrcode, maxErrcode))
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
			return
		} else {
			bytes, err := os.ReadFile(registerPath)
			if err != nil {
				return
			}
			scanner := bufio.NewScanner(strings.NewReader(string(bytes)))
			var lastPort uint32
			var lastErrcode uint32
			var find bool
			for scanner.Scan() {
				line := scanner.Text()
				split := strings.Split(line, ",")
				if len(split) != 4 {
					continue
				}
				srvName := split[0]
				if srvName == pbCtx.ServiceName {
					genVo.Port = utils.AtoUint32(split[1])
					find = true
					break
				}
				lastPort = utils.AtoUint32(split[1])
				lastErrcode = utils.AtoUint32(split[3])
			}
			if lastPort == 0 {
				lastPort = maxPort
			} else {
				lastPort += 1
			}
			if lastErrcode == 0 {
				lastErrcode = minErrcode
			}
			genVo.Port = lastPort
			if !find {
				utils.AppendWriteFile(registerPath, fmt.Sprintf("%s,%d,%d,%d\n", pbCtx.ServiceName, lastPort, lastErrcode, lastErrcode+step))
			}
		}

		template, err := temp.GenCodeByTemplate(gtpl.GenServerCmdApplication, &genVo)
		if err != nil {
			return
		}

		err = utils.CreateAndWriteFile(yamlPath, template)
		if err != nil {
			return
		}
	}
}
