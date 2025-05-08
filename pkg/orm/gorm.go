/**
 * @Author: zjj
 * @Date: 2025/8/8
 * @Desc:
**/

package orm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

func NewGormEngine(dsn string) *gorm.DB {
	ormLog := NewOrmLog(time.Second * 5)
	ormLog.skipCall = 11
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,  // 是否单表，命名是否复数
			NoLowerCase:   false, // 是否关闭驼峰命名
		},

		PrepareStmt: true, // 预编译 在执行任何 SQL 时都会创建一个 prepared statement 并将其缓存，以提高后续的效率

		Logger: ormLog,
	})
	if err != nil {
		panic(err)
	}

	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := db.DB()

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}
