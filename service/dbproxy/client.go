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
	"net/http"
)

func InsertModel(ctx uctx.IUCtx, req *InsertModelReq) (*InsertModelRsp, error) {
	var resp InsertModelRsp
	err := rpc.DoRequest(ctx, "dbproxy", "/dbproxyserver/InsertModel", http.MethodPost, req, &resp)
	if err != nil {
		return nil, gerr.Wrap(err)
	}
	return &resp, nil
}

func DelModel(ctx uctx.IUCtx, req *DelModelReq) (*DelModelRsp, error) {
	var resp DelModelRsp
	err := rpc.DoRequest(ctx, "dbproxy", "/dbproxyserver/DelModel", http.MethodPost, req, &resp)
	if err != nil {
		return nil, gerr.Wrap(err)
	}
	return &resp, nil
}

func UpdateModel(ctx uctx.IUCtx, req *UpdateModelReq) (*UpdateModelRsp, error) {
	var resp UpdateModelRsp
	err := rpc.DoRequest(ctx, "dbproxy", "/dbproxyserver/UpdateModel", http.MethodPost, req, &resp)
	if err != nil {
		return nil, gerr.Wrap(err)
	}
	return &resp, nil
}

func BatchInsertModel(ctx uctx.IUCtx, req *BatchInsertModelReq) (*BatchInsertModelRsp, error) {
	var resp BatchInsertModelRsp
	err := rpc.DoRequest(ctx, "dbproxy", "/dbproxyserver/BatchInsertModel", http.MethodPost, req, &resp)
	if err != nil {
		return nil, gerr.Wrap(err)
	}
	return &resp, nil
}

func SetModel(ctx uctx.IUCtx, req *SetModelReq) (*SetModelRsp, error) {
	var resp SetModelRsp
	err := rpc.DoRequest(ctx, "dbproxy", "/dbproxyserver/SetModel", http.MethodPost, req, &resp)
	if err != nil {
		return nil, gerr.Wrap(err)
	}
	return &resp, nil
}

func GetModelList(ctx uctx.IUCtx, req *GetModelListReq) (*GetModelListRsp, error) {
	var resp GetModelListRsp
	err := rpc.DoRequest(ctx, "dbproxy", "/dbproxy/GetModelList", http.MethodPost, req, &resp)
	if err != nil {
		return nil, gerr.Wrap(err)
	}
	return &resp, nil
}
