package migrate

const (
	TINYBLOB   = "tinyblob"
	BLOB       = "blob"
	MEDIUMBLOB = "mediumblob"

	Int      = "int(10)"
	TinyInt  = "tinyint(3)"
	SmallInt = "smallint(5)"
	BigInt   = "bigint(20)"

	NO  = ""
	MUL = "MUL" // 普通索引
	UNI = "UNI" // 唯一索引
	PRI = "PRI" // （主键）

	MySqlEngine = "INNODB" // 数据库引擎

	CreateSQLHead = "create table %s \n(\n"
	CreateSQLTail = "\n)\nENGINE=" + MySqlEngine + " DEFAULT CHARSET=utf8 COMMENT '%s'"

	ProcedureDropTemplate = "drop procedure if exists %s;"
)
