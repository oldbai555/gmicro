func {{.RpcName}}(ctx context.Context, req *{{.Client}}.{{.RpcReq}}) (*{{.Client}}.{{.RpcRsp}}, error) {
	var rsp {{.Client}}.{{.RpcRsp}}
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, gerr.Wrap(err)
	}

	listRsp, err := a.Get{{.ModelName}}List(ctx, &{{.Client}}.Get{{.ModelName}}ListReq{
    		ListOption: req.ListOption.
    			SetSkipTotal().
    			AddOpt(core.DefaultListOption_DefaultListOptionSelect, {{.Client}}.FieldId_),
    	})
    if err != nil {
        return nil, gerr.Wrap(err)
    }

    if len(listRsp.List) == 0 {
        log.Infof("list is empty")
        return &rsp, nil
    }

    idList := utils.PluckUint64List(listRsp.List, {{.Client}}.FieldId)
    _, err = Orm{{.ModelName}}.NewBaseScope().WhereIn({{.Client}}.FieldId_, idList).Delete(uCtx)
    if err != nil {
        return nil, gerr.Wrap(err)
    }

	return &rsp, err
}
