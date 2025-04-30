/**
 * @Author: zjj
 * @Date: 2025/4/30
 * @Desc:
**/

package dbproxy

import (
	"gmicro/pkg/gerr"
	"gmicro/pkg/rpc"
	"gmicro/pkg/uctx"
	"gmicro/service/gormx"
	"net/http"
)

func InsertModel(ctx uctx.IUCtx, req *gormx.InsertModelReq) (*gormx.InsertModelRsp, error) {
	var resp gormx.InsertModelRsp
	err := rpc.DoRequest(ctx, "dbproxy", "/dbproxyserver/InsertModel", http.MethodPost, req, &resp)
	if err != nil {
		return nil, gerr.Wrap(err)
	}
	return &resp, nil
}

func DelModel(ctx uctx.IUCtx, req *gormx.DelModelReq) (*gormx.DelModelRsp, error) {
	var resp gormx.DelModelRsp
	err := rpc.DoRequest(ctx, "dbproxy", "/dbproxyserver/DelModel", http.MethodPost, req, &resp)
	if err != nil {
		return nil, gerr.Wrap(err)
	}
	return &resp, nil
}

func UpdateModel(ctx uctx.IUCtx, req *gormx.UpdateModelReq) (*gormx.UpdateModelRsp, error) {
	var resp gormx.UpdateModelRsp
	err := rpc.DoRequest(ctx, "dbproxy", "/dbproxyserver/UpdateModel", http.MethodPost, req, &resp)
	if err != nil {
		return nil, gerr.Wrap(err)
	}
	return &resp, nil
}

func BatchInsertModel(ctx uctx.IUCtx, req *gormx.BatchInsertModelReq) (*gormx.BatchInsertModelRsp, error) {
	var resp gormx.BatchInsertModelRsp
	err := rpc.DoRequest(ctx, "dbproxy", "/dbproxyserver/BatchInsertModel", http.MethodPost, req, &resp)
	if err != nil {
		return nil, gerr.Wrap(err)
	}
	return &resp, nil
}

func SetModel(ctx uctx.IUCtx, req *gormx.SetModelReq) (*gormx.SetModelRsp, error) {
	var resp gormx.SetModelRsp
	err := rpc.DoRequest(ctx, "dbproxy", "/dbproxyserver/SetModel", http.MethodPost, req, &resp)
	if err != nil {
		return nil, gerr.Wrap(err)
	}
	return &resp, nil
}

func GetModelList(ctx uctx.IUCtx, req *gormx.GetModelListReq) (*gormx.GetModelListRsp, error) {
	var resp gormx.GetModelListRsp
	err := rpc.DoRequest(ctx, "dbproxy", "/dbproxyserver/GetModelList", http.MethodPost, req, &resp)
	if err != nil {
		return nil, gerr.Wrap(err)
	}
	return &resp, nil
}
