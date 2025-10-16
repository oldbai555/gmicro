func {{.RpcName}}(ctx context.Context, req *{{.Client}}.{{.RpcReq}}) (*{{.Client}}.{{.RpcRsp}}, error) {
	var rsp {{.Client}}.{{.RpcRsp}}
	var err error

	// 通用处理器
	var newOptionsProcessor = func(db query.IModel{{.ModelName}}Do, listOption *base.ListOption) *base.OptionsProcessor {
		tableName := db.TableName()
		return base.NewOptionsProcessor(listOption).
			AddStringList(
				base.DefaultListOption_DefaultListOptionSelect,
				func(valList []string) error {
					var fields []field.Expr
					for _, val := range valList {
						fields = append(fields, field.NewField(tableName, val))
					}
					db.Select(fields...)
					return nil
				}).
			AddUint32(
				base.DefaultListOption_DefaultListOptionOrderBy,
				func(val uint32) error {
					if val == uint32(base.DefaultOrderBy_DefaultOrderByCreatedAtDesc) {
						newField := field.NewField(tableName, "created_at")
						db.Order(newField.Desc())
					} else if val == uint32(base.DefaultOrderBy_DefaultOrderByCreatedAtAcs) {
						newField := field.NewField(tableName, "created_at")
						db.Order(newField.Asc())
					} else if val == uint32(base.DefaultOrderBy_DefaultOrderByIdDesc) {
						newField := field.NewField(tableName, "id")
						db.Order(newField.Desc())
					}
					return nil
				}).
			AddStringList(
				base.DefaultListOption_DefaultListOptionGroupBy,
				func(valList []string) error {
					var fields []field.Expr
					for _, val := range valList {
						fields = append(fields, field.NewField(tableName, val))
					}
					db.Group(fields...)
					return nil
				}).
			AddBool(
				base.DefaultListOption_DefaultListOptionWithTrash,
				func(val bool) error {
					if val {
						db.Unscoped()
					}
					return nil
				}).
			AddUint64List(
				base.DefaultListOption_DefaultListOptionIdList,
				func(valList []uint64) error {
					newField := field.NewUint64(tableName, "id")
					if len(valList) == 1 {
						db.Where(newField.Eq(valList[0]))
					} else {
						db.Where(newField.In(valList...))
					}
					return nil
				}).
			AddTimeStampRange(
				base.DefaultListOption_DefaultListOptionCreatedAt,
				func(begin, end uint32) error {
					newField := field.NewUint32(tableName, "created_at")
					db.Where(newField.Between(begin, end))
					return nil
				}).
			AddUint64List(
				base.DefaultListOption_DefaultListOptionCreatorIdList,
				func(valList []uint64) error {
					newField := field.NewUint64(tableName, "creator_id")
					if len(valList) == 1 {
						db.Where(newField.Eq(valList[0]))
					} else {
						db.Where(newField.In(valList...))
					}
					return nil
				})
	}

	// 处理业务
	listOption := req.ListOption

	db := query.Model{{.ModelName}}.ReadDB()
	processor := newOptionsProcessor(db, listOption)
	defer func() {
		base.ReleaseOptionsProcessor(processor)
	}()

	err = processor. // 往这追加非通用处理的listOption
				Process()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, gerr.Wrap(err)
	}
	var result []*model.Model{{.ModelName}}
	if !listOption.SkipTotal {
		var count int64
		limit := listOption.Limit
		if limit > 2000 {
			limit = 2000
		}
		offset := listOption.Offset
		result, count, err = db.FindByPage(int(offset), int(limit))
		rsp.Paginate = &base.Paginate{
			Total:  uint32(count),
			Limit:  limit,
			Offset: offset,
		}
	} else {
		result, err = db.Find()
	}
	if err != nil {
		return nil, gerr.Wrap(err)
	}
	for _, data := range result {
		// 如果能手动赋值尽量手动，想使用copier强烈建议只拷贝单个对象影响不是很大
		var pbData {{.Client}}.Model{{.ModelName}}
		_ = copier.Copy(&pbData, data)
		rsp.List = append(rsp.List, &pbData)
	}

	return &rsp, err
}
