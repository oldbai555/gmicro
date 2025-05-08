package impl

import (
	"fmt"
	"gmicro/service/dbproxyserver/engine"
	"gmicro/service/dbproxyserver/mysql"
)

func InitState() error {
	dns := fmt.Sprintf("root:game@2023@tcp(192.168.61.231:3306)/u3dv1_actor_24?charset=utf8&parseTime=True&loc=Local")
	gormEngine := mysql.NewGormEngine(dns)
	engine.SetOrmEngine(gormEngine)
	return nil
}
