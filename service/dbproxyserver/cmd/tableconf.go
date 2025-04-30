package main

import "gmicro/service/dbproxyserver/migrate"

var Tables = []*migrate.TableConf{
	{
		Name:    "testTable",
		Comment: "测试表",
		Fields: []*migrate.ColumnConf{
			{Name: "id", Type: "bigint", Comment: "ID", AutoIncrement: true},
			{Name: "name", Type: "varchar", Comment: "名字", Size: 128},
			{Name: "dept", Type: "int", Comment: "部门", Default: "0"},
			{Name: "role", Type: "int", Comment: "角色", Default: "0"},
		},
		Keys: []*migrate.ColumnKeySt{
			{Type: "pri", Name: "id"},
			{Type: "mul", Name: "dept"},
		},
	},
}
