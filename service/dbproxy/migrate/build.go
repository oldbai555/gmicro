package migrate

import (
	"fmt"
	"gmicro/pkg/log"
	"gmicro/service/dbproxy/engine"
	"gmicro/service/dbproxy/mysql"
)

func BuildTables() {
	for _, table := range tables {
		table.Build()
	}
}

func Exec(sql string) {
	if len(sql) <= 0 {
		return
	}
	gormEngine := engine.GetOrmEngine().(*mysql.GormEngine)
	if gormEngine == nil {
		return
	}
	err := gormEngine.GetDB("").Exec(sql).Error
	if nil != err {
		log.Errorf("err:%v", err)
	}
}

func DropColumnSQL(tblName string, columnName string) string {
	return fmt.Sprintf("alter table %s drop column %s;", tblName, columnName)
}

func DropKeySQL(tblName string, columnName string, kt string) string {
	switch kt {
	case MUL:
		return fmt.Sprintf("alter table %s drop index %s", tblName, columnName)
	case PRI:
		return fmt.Sprintf("alter table %s drop primary key", tblName)
	case UNI:
		return fmt.Sprintf("alter table %s drop index %s", tblName, columnName)
	default:
		log.Fatalf("drop key sql error")
	}
	return ""
}
