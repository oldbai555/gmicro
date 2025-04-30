package migrate

import (
	"fmt"
	"gmicro/pkg/log"
)

type ProcedureSt struct {
	Name string
	SQL  string
}

func (st *ProcedureSt) Build() {
	Exec(st.CreateDropSQL())
	Exec(st.CreateProcedureSQL())
}

func (st *ProcedureSt) CreateDropSQL() string {
	return fmt.Sprintf(ProcedureDropTemplate, st.Name)
}

func (st *ProcedureSt) CreateProcedureSQL() string {
	if len(st.SQL) <= 0 {
		log.Fatalf("%s 储存过程SQL语句为空", st.Name)
	}
	return st.SQL
}

func BuildProcedures(procedures []*ProcedureSt) {
	for _, procedure := range procedures {
		procedure.Build()
	}
}
