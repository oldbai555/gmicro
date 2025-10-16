func {{.RpcName}}(ctx context.Context, req *{{.Client}}.{{.RpcReq}}) (*{{.Client}}.{{.RpcRsp}}, error) {
	var rsp {{.Client}}.{{.RpcRsp}}
	var err error

	db := query.Model{{.ModelName}}.WriteDB()
	var data model.Model{{.ModelName}}
	_ = copier.Copy(&data, req.Data)
	err = db.Create(&data)
	if err != nil {
		return nil, gerr.Wrap(err)
	}
	_ = copier.Copy(req.Data, &data)
	rsp.Data = req.Data

	return &rsp, err
}
