package migrate

import (
	"fmt"
	"gmicro/pkg/log"
)

type ColumnKeySt struct {
	Type string
	Name string
}

func (st *ColumnKeySt) CreateKeySQL(table *TableSt) string {
	if st.Type == MUL {
		return fmt.Sprintf("KEY %s(%s)", st.Name, st.Name)
	} else if st.Type == PRI {
		return table.CreatePRIKeySQL()
	} else if st.Type == UNI {
		return fmt.Sprintf("UNIQUE KEY %s(%s)", st.Name, st.Name)
	}
	log.Fatalf("CreateKeySQL error!!! name:%s, type:%s, table:%s", st.Name, st.Type, table.Name)
	return ""
}

func (st *ColumnKeySt) AddKeySQL(table *TableSt) string {
	switch st.Type {
	case MUL:
		return fmt.Sprintf("ALTER TABLE %s ADD INDEX %s(%s)", table.Name, st.Name, st.Name)
	case PRI:
		return table.AddPRIKeySQL()
	case UNI:
		return fmt.Sprintf("ALTER TABLE %s ADD UNIQUE (%s)", table.Name, st.Name)
	}
	log.Fatalf("%s %s add key sql error", table.Name, st.Name)
	return ""
}

func (st *ColumnKeySt) IsEqual(info *MysqlColumnSt) bool {
	return st.Type == info.Key
}

var KeyFuncMap = map[string]func(name string) *ColumnKeySt{}

func init() {
	KeyFuncMap["pri"] = func(name string) *ColumnKeySt {
		return &ColumnKeySt{Type: PRI, Name: name}
	}
	KeyFuncMap["mul"] = func(name string) *ColumnKeySt {
		return &ColumnKeySt{Type: MUL, Name: name}
	}
	KeyFuncMap["uni"] = func(name string) *ColumnKeySt {
		return &ColumnKeySt{Type: UNI, Name: name}
	}
}
