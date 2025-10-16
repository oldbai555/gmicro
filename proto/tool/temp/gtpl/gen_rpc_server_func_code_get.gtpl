func {{.RpcName}}(ctx context.Context, req *{{.Client}}.{{.RpcReq}}) (*{{.Client}}.{{.RpcRsp}}, error) {
	var rsp {{.Client}}.{{.RpcRsp}}
	var err error

	db := query.Model{{.ModelName}}.ReadDB()
	// other logic ...
    data, err := db.GetById(req.Id)
	if err != nil {
		return nil, gerr.Wrap(err)
	}

  	var pbData {{.Client}}.Model{{.ModelName}}
   	_ = copier.Copy(&pbData, data)
	rsp.Data = &pbData

	return &rsp, err
}
