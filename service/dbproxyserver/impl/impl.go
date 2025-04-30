/**
 * @Author: zjj
 * @Date: 2025/4/30
 * @Desc:
**/

package impl

import (
	"context"
	"gmicro/pkg/log"
	"gmicro/pkg/uctx"
	"gmicro/service/dbproxy"
	"gmicro/service/dbproxyserver/engine"
)

func InsertModel(ctx context.Context, req *dbproxy.InsertModelReq) (*dbproxy.InsertModelRsp, error) {
	rsp, err := engine.GetOrmEngine().InsertModel(uctx.NewBaseUCtx(), req)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rsp, err
}

func DelModel(ctx context.Context, req *dbproxy.DelModelReq) (*dbproxy.DelModelRsp, error) {
	rsp, err := engine.GetOrmEngine().DelModel(uctx.NewBaseUCtx(), req)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rsp, err
}

func UpdateModel(ctx context.Context, req *dbproxy.UpdateModelReq) (*dbproxy.UpdateModelRsp, error) {
	rsp, err := engine.GetOrmEngine().UpdateModel(uctx.NewBaseUCtx(), req)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rsp, err
}

func BatchInsertModel(ctx context.Context, req *dbproxy.BatchInsertModelReq) (*dbproxy.BatchInsertModelRsp, error) {
	rsp, err := engine.GetOrmEngine().BatchInsertModel(uctx.NewBaseUCtx(), req)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rsp, err
}

func SetModel(ctx context.Context, req *dbproxy.SetModelReq) (*dbproxy.SetModelRsp, error) {
	rsp, err := engine.GetOrmEngine().SetModel(uctx.NewBaseUCtx(), req)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rsp, err
}

func GetModelList(ctx context.Context, req *dbproxy.GetModelListReq) (*dbproxy.GetModelListRsp, error) {
	rsp, err := engine.GetOrmEngine().GetModelList(uctx.NewBaseUCtx(), req)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return rsp, err
}
