func {{.RpcName}}(ctx context.Context, req *{{.Client}}.{{.RpcReq}}) (*{{.Client}}.{{.RpcRsp}}, error) {
	var rsp {{.Client}}.{{.RpcRsp}}
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, gerr.Wrap(err)
	}

	db := query.ModelFile.WriteDB()
	err = db.Create(req.Data)
	if err != nil {
		return nil, gerr.Wrap(err)
	}
    rsp.Data = req.Data

	return &rsp, err
}
