func {{.RpcName}}(ctx context.Context, req *{{.Client}}.{{.RpcReq}}) (*{{.Client}}.{{.RpcRsp}}, error) {
	var rsp {{.Client}}.{{.RpcRsp}}
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, gerr.Wrap(err)
	}

    data, err := Orm{{.ModelName}}.NewBaseScope().Where({{.Client}}.FieldId_, req.Id).First(uCtx)
    if err != nil {
        return nil, gerr.Wrap(err)
    }
    rsp.Data = data

	return &rsp, err
}
