package main

import (
	"fmt"
	"gmicro/pkg/log"
	"gmicro/service/dbproxy/engine"
	"gmicro/service/dbproxy/migrate"
	"gmicro/service/dbproxy/mysql"
)

func main() {
	log.InitLogger(
		log.WithAppName("dbproxy"),
		log.WithLevel(log.DebugLevel),
		log.WithScreen(true),
	)
	dns := fmt.Sprintf("root:game@2023@tcp(192.168.61.231:3306)/u3dv1_actor_24?charset=utf8&parseTime=True&loc=Local")
	gormEngine := mysql.NewGormEngine(dns)
	engine.SetOrmEngine(gormEngine)
	migrate.Parse(Tables)
	migrate.BuildTables()
}
