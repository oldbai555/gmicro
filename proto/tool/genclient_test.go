/**
 * @Author: zjj
 * @Date: 2025/5/6
 * @Desc:
**/

package main

import (
	"fmt"
	"gmicro/pkg/exec"
	"gmicro/pkg/log"
	"gmicro/proto/tool/option"
	parse2 "gmicro/proto/tool/parse"
	"gmicro/proto/tool/vo"
	"testing"
)

func Test1(t *testing.T) {
	defer log.Flush()
	log.InitLogger(log.WithAppName("Test1"))
	GenAll("test")
}

func GenSrvAutoDb(protoName string) {
	pbCtx, err := parse2.ParsePb(fmt.Sprintf("C:\\Users\\zhangjianjun\\Desktop\\222\\gmicro\\proto\\%s.proto", protoName))
	if err != nil {
		return
	}
	params := &vo.CodeFuncParams{
		PbName:    protoName,
		PbDir:     "C:\\Users\\zhangjianjun\\Desktop\\222\\gmicro\\proto",
		OutputDir: "C:\\Users\\zhangjianjun\\Desktop\\222\\gmicro\\service\\",
		GitPath:   "",
	}
	option.WithGenServerGormGen()(pbCtx, params)
}

func GenAll(protoName string) {
	pbCtx, err := parse2.ParsePb(fmt.Sprintf("C:\\Users\\zhangjianjun\\Desktop\\222\\gmicro\\proto\\%s.proto", protoName))
	if err != nil {
		return
	}
	params := &vo.CodeFuncParams{
		PbName:    protoName,
		PbDir:     "C:\\Users\\zhangjianjun\\Desktop\\222\\gmicro\\proto",
		OutputDir: "C:\\Users\\zhangjianjun\\Desktop\\222\\gmicro\\service\\" + protoName,
		GitPath:   "",
	}
	exec.ProtocGo("C:\\Users\\zhangjianjun\\Desktop\\222", "C:\\Users\\zhangjianjun\\Desktop\\222\\gmicro\\proto", "base.proto")
	exec.ProtocGo("C:\\Users\\zhangjianjun\\Desktop\\222", "C:\\Users\\zhangjianjun\\Desktop\\222\\gmicro\\proto", protoName+".proto")
	option.WithGenClientCmdAutoGen()(pbCtx, params)
	option.WithGenClientErrcodeAutoGen()(pbCtx, params)
	option.WithGenClientFieldAutoGen()(pbCtx, params)
	option.WithGenClientRpcAutoGen()(pbCtx, params)
	option.WithGenClientTableNameAutoGen()(pbCtx, params)

	params = &vo.CodeFuncParams{
		PbName:    protoName,
		PbDir:     "C:\\Users\\zhangjianjun\\Desktop\\222\\gmicro\\proto",
		OutputDir: "C:\\Users\\zhangjianjun\\Desktop\\222\\gmicro\\service\\",
		GitPath:   "",
	}
	option.WithGenServerCmdMainAutoGen()(pbCtx, params)
	option.WithGenServerImplStateAutoGen()(pbCtx, params)
	option.WithGenServerImplRpcAutoGen()(pbCtx, params)
	option.WithGenServerCmdApplicationAutoGen()(pbCtx, params)
	option.WithGenServerCmdCmdAutoGen()(pbCtx, params)
	option.WithGenServerGormGen()(pbCtx, params)
	exec.GoFmt("C:\\Users\\zhangjianjun\\Desktop\\222\\gmicro")
}
