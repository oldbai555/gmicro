/**
 * @Author: zjj
 * @Date: 2024/6/18
 * @Desc:
**/

package engine

import (
	"gmicro/pkg/uctx"
	"gmicro/service/dbproxy"
)

type IOrmEngine interface {
	GetModelList(ctx uctx.IUCtx, req *dbproxy.GetModelListReq) (*dbproxy.GetModelListRsp, error)
	InsertModel(ctx uctx.IUCtx, req *dbproxy.InsertModelReq) (*dbproxy.InsertModelRsp, error)
	DelModel(ctx uctx.IUCtx, req *dbproxy.DelModelReq) (*dbproxy.DelModelRsp, error)
	UpdateModel(ctx uctx.IUCtx, req *dbproxy.UpdateModelReq) (*dbproxy.UpdateModelRsp, error)
	BatchInsertModel(ctx uctx.IUCtx, req *dbproxy.BatchInsertModelReq) (*dbproxy.BatchInsertModelRsp, error)
	SetModel(ctx uctx.IUCtx, req *dbproxy.SetModelReq) (*dbproxy.SetModelRsp, error)
	RegObjectType(objType ...*dbproxy.ModelObjectType)
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
