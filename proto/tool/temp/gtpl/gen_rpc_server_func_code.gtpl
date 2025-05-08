func {{.RpcName}}(ctx uctx.IUCtx, req *{{.Client}}.{{.RpcReq}}) (*{{.Client}}.{{.RpcRsp}}, error) {
	var rsp {{.Client}}.{{.RpcRsp}}
	var err error

	return &rsp, err
}
