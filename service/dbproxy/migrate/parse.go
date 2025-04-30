package migrate

import (
	"fmt"
	"gmicro/pkg/log"
	"strings"
)

func Parse(tableCfgs []*TableConf) {
	parseTable := func(cfg *TableConf, tableName string) {
		table := NewTableSt(tableName, cfg.Comment, len(cfg.Fields))
		for _, field := range cfg.Fields {
			var column ColumnInterface
			if field.Type == "varchar" {
				column = NewStringColumn(field.Name, field.Size, field.Comment, field.Default)
			} else if field.Type == "datetime" {
				column = NewDateTimeColumn(field.Name, field.Comment)
			} else if strings.Contains(field.Type, "blob") {
				column = NewBlobColumn(field.Name, field.Type, field.Comment)
			} else if strings.Contains(field.Type, "int") {
				column = NewIntColumn(field.Name, field.Type, field.Unsigned, field.Comment, field.AutoIncrement, field.Default)
			} else {
				log.Fatalf("字段类型定义错误。 table:%s, field:%s", tableName, field.Name)
			}
			table.AddColumn(column)
		}
		for _, key := range cfg.Keys {
			table.AddKey(KeyFuncMap[key.Type](key.Name))
		}
		tables[tableName] = table
	}
	for _, line := range tableCfgs {
		if line.Scaled {
			for i := line.ScaleMinSeq; i <= line.ScaleMaxSeq; i++ {
				parseTable(line, fmt.Sprintf("%s_%d", line.Name, i))
			}
			continue
		}
		parseTable(line, line.Name)
	}
}
