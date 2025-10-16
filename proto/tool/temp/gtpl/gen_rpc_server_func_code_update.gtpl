func {{.RpcName}}(ctx context.Context, req *{{.Client}}.{{.RpcReq}}) (*{{.Client}}.{{.RpcRsp}}, error) {
	var rsp {{.Client}}.{{.RpcRsp}}
	var err error

    db := query.ModelFile.ReadDB()
    data, err := db.GetById(req.Data.Id)
    if err != nil {
    	return nil, gerr.Wrap(err)
    }

    db = query.Model{{.ModelName}}.WriteDB()
    // other logic ...
    _, err = db.Where(query.Model{{.ModelName}}.ID.Eq(data.ID)).UpdateColumnSimple()
    if err != nil {
		return nil, gerr.Wrap(err)
	}

	return &rsp, err
}
