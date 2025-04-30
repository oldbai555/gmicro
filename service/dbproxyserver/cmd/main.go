package main

import (
	"fmt"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"gmicro/common"
	"gmicro/pkg/autocmd"
	"gmicro/pkg/log"
	"gmicro/service/dbproxyserver/engine"
	"gmicro/service/dbproxyserver/mysql"
	"google.golang.org/protobuf/proto"
	"net/http"
)

func main() {
	log.InitLogger(
		log.WithAppName("dbproxyserver"),
		log.WithLevel(log.DebugLevel),
		log.WithScreen(true),
	)
	dns := fmt.Sprintf("root:game@2023@tcp(192.168.61.231:3306)/u3dv1_actor_24?charset=utf8&parseTime=True&loc=Local")
	gormEngine := mysql.NewGormEngine(dns)
	engine.SetOrmEngine(gormEngine)
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
			Formatter: common.NewLogFormatter("dbproxy"),
			Output:    log.GetWriter(),
		}),
	)

	// proto 生成的路由
	for _, cmd := range cmdList {
		registerCmd(router, cmd)
	}
	router.GET("hellp", func(context *gin.Context) {
		context.JSON(200, "你好")
	})

	// 注册自定义路由
	ginpprof.Wrap(router)

	ginSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 20001),
		Handler: router,
	}

	// 启动服务
	log.Infof("====== start gin %s server, port is %d ======", "dbproxy", 20001)
	err := ginSrv.ListenAndServe()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	return
}

func registerCmd(r *gin.Engine, cmd *autocmd.Cmd) {
	cmd.WithHandleError(func(ctx *gin.Context, err error) {

	}).WithHandleResult(func(ctx *gin.Context, result proto.Message) {

	})
	r.POST(cmd.Path, cmd.GinPost)
}
