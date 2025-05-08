func {{.RpcName}}(ctx context.Context, req *{{.Client}}.{{.RpcReq}}) (*{{.Client}}.{{.RpcRsp}}, error) {
	var rsp {{.Client}}.{{.RpcRsp}}
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, gerr.Wrap(err)
	}

    db := Orm{{.ModelName}}.NewList(req.ListOption)
	err = gormx.ProcessDefaultOptions[*{{.Client}}.Model{{.ModelName}}](req.ListOption, db)
	if err != nil {
		return nil, gerr.Wrap(err)
	}

	err = core.NewOptionsProcessor(req.ListOption).
		Process()

	rsp.Paginate, err = db.FindPaginate(uCtx, &rsp.List)
	if err != nil {
		return nil, gerr.Wrap(err)
	}

	return &rsp, err
}
