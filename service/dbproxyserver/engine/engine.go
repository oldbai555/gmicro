/**
 * @Author: zjj
 * @Date: 2024/6/18
 * @Desc:
**/

package engine

import (
	"gmicro/pkg/uctx"
	"gmicro/service/gormx"
)

type IOrmEngine interface {
	GetModelList(ctx uctx.IUCtx, req *gormx.GetModelListReq) (*gormx.GetModelListRsp, error)
	InsertModel(ctx uctx.IUCtx, req *gormx.InsertModelReq) (*gormx.InsertModelRsp, error)
	DelModel(ctx uctx.IUCtx, req *gormx.DelModelReq) (*gormx.DelModelRsp, error)
	UpdateModel(ctx uctx.IUCtx, req *gormx.UpdateModelReq) (*gormx.UpdateModelRsp, error)
	BatchInsertModel(ctx uctx.IUCtx, req *gormx.BatchInsertModelReq) (*gormx.BatchInsertModelRsp, error)
	SetModel(ctx uctx.IUCtx, req *gormx.SetModelReq) (*gormx.SetModelRsp, error)
	RegObjectType(objType ...*gormx.ModelObjectType)
	Begin() (string, error)
	Rollback(trId string) error
	Commit(trId string) error
}

var ormEngine IOrmEngine

func SetOrmEngine(val IOrmEngine) {
	ormEngine = val
}

func GetOrmEngine() IOrmEngine {
	if ormEngine == nil {
		panic("orm engine is nil")
	}
	return ormEngine
}
