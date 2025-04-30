/**
 * @Author: zjj
 * @Date: 2025/4/30
 * @Desc:
**/

package gormx

import (
	"gmicro/pkg/log"
	"gmicro/pkg/uctx"
	"testing"
)

type ModelTestTable struct {
	Id int64 `json:"id"`
}

func TestNewBaseModel(t *testing.T) {
	log.InitLogger()
	OrmTestTable := NewBaseModel[*ModelTestTable](ModelConfig{
		NotFoundErrCode: int32(1),
		Db:              "u3dv1_actor_24",
	})
	ctx := uctx.NewBaseUCtx()
	first, err := OrmTestTable.First(ctx)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	t.Log(first)
}
