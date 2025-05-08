/**
 * @Author: zjj
 * @Date: 2025/5/7
 * @Desc:
**/

package rpc

import (
	"context"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"gmicro/common"
	"gmicro/pkg/autocmd"
	"gmicro/pkg/gerr"
	"gmicro/pkg/log"
	"gmicro/pkg/signal"
	"gmicro/pkg/uctx"
	"google.golang.org/protobuf/proto"
	"net/http"
	"os"
)

func ServerRun(serverName, srvAddr string, cmdList []*autocmd.Cmd) error {
	if len(serverName) == 0 || len(srvAddr) == 0 {
		return gerr.NewCustomErr("server name or srv addr is null")
	}
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.GetWriter()
	gin.DefaultErrorWriter = log.GetWriter()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-25s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	router := gin.New()

	router.Use(
		gin.Recovery(),
		gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: common.NewLogFormatter(serverName),
			Output:    log.GetWriter(),
		}),
	)

	// proto 生成的路由
	for _, cmd := range cmdList {
		registerCmd(router, cmd)
	}

	ginpprof.Wrap(router)

	ginSrv := &http.Server{
		Addr:    srvAddr,
		Handler: router,
	}

	signal.RegV2(func(signal os.Signal) error {
		log.Infof("exit: close ginSrv, signal[%v]", signal)
		if ginSrv != nil {
			err := ginSrv.Shutdown(context.Background())
			if err != nil {
				log.Errorf("err:%v", err)
			}
		}
		return nil
	})

	// 服务注册

	// 启动服务
	log.Infof("====== start gin %s server, srvAddr: %s ======", serverName, srvAddr)
	err := ginSrv.ListenAndServe()
	if err != nil {
		log.Warnf("err is %v", err)
		return err
	}
	return nil
}

func registerCmd(r *gin.Engine, cmd *autocmd.Cmd) {
	cmd.WithGenIUCtx(func(ctx *gin.Context) uctx.IUCtx {
		return uctx.NewBaseUCtx()
	}).WithCheckAuthF(func(nCtx uctx.IUCtx) (extInfo interface{}, err error) {
		return nil, nil
	}).WithHandleError(func(ctx *gin.Context, err error) {
		var rsp = map[string]interface{}{
			"data":    "",
			"errcode": gerr.GetErrCode(err),
			"errmsg":  err.Error(),
			"hint":    "",
		}
		ctx.JSON(http.StatusOK, rsp)
	}).WithHandleResult(func(ctx *gin.Context, result proto.Message) {
		var rsp = map[string]interface{}{
			"data":    result,
			"errcode": 0,
			"errmsg":  "",
			"hint":    "",
		}
		ctx.JSON(http.StatusOK, rsp)
	})
	r.POST(cmd.Path, cmd.GinPost)
}
