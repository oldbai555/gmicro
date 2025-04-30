package migrate

import (
	"fmt"
	utils "gmicro/common"
	"strings"
)

type StringColumnSt struct {
	ColumnBaseSt
}

func find(s string) string {
	start := strings.Index(s, "(")
	end := strings.Index(s, ")")
	if start >= 0 && end > start {
		return s[start+1 : end]
	}
	return ""
}

// IsCompatible 是否兼容
func (st *StringColumnSt) IsCompatible(info *MysqlColumnSt) bool {
	size := find(info.Type)
	//r := regexp.MustCompile(`\(.*?\)`)
	//size := r.FindStringSubmatch(info.Type)
	if len(size) <= 0 {
		return false
	}
	if st.Size < utils.Atoi(size) {
		return false
	}

	return true
}

func NewStringColumn(name string, size int, comment string, def string) ColumnInterface {
	column := &StringColumnSt{}
	column.Name = name
	column.Type = fmt.Sprintf("varchar(%d)", size)
	column.Size = size
	column.Comment = comment
	column.Default = def
	return column
}
