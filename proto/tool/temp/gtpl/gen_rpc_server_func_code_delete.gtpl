func {{.RpcName}}(ctx context.Context, req *{{.Client}}.{{.RpcReq}}) (*{{.Client}}.{{.RpcRsp}}, error) {
	var rsp {{.Client}}.{{.RpcRsp}}
	var err error

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		return nil, gerr.Wrap(err)
	}

	listRsp, err := Get{{.ModelName}}List(ctx, &{{.Client}}.Get{{.ModelName}}ListReq{
    		ListOption: req.ListOption.
    			SetSkipTotal().
    			AddOpt(base.DefaultListOption_DefaultListOptionSelect, {{.Client}}.FieldId_),
    	})
    if err != nil {
        return nil, gerr.Wrap(err)
    }

    if len(listRsp.List) == 0 {
        log.Infof("list is empty")
        return &rsp, nil
    }

    idList := utils.PluckUint64List(listRsp.List, {{.Client}}.FieldId)
	db := query.ModelFile.WriteDB()
	_, err = db.Where(query.ModelFile.ID.In(idList...)).Delete()
	if err != nil {
		return nil, gerr.Wrap(err)
	}

	return &rsp, err
}
