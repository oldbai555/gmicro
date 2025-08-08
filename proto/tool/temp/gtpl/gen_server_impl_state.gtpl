package impl

import (
	"gmicro/pkg/gerr"
	"gmicro/pkg/orm"
	"gmicro/pkg/syscfg"
	"gmicro/service/{{.ProtoName}}server/impl/query"
)

func InitState() error {
	err := initGorm()
	if err != nil {
		return gerr.Wrap(err)
	}
	return nil
}

func initGorm() error {
	mysqlConf := syscfg.NewGormMysqlConf("")
	if mysqlConf == nil {
		return gerr.RecordNotFound
	}
	db := orm.NewGormEngine(mysqlConf.Dsn())
	if db == nil {
		return gerr.OrmInitFailed
	}
	query.SetDefault(db)
	return nil
}
