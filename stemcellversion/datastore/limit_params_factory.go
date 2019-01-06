package datastore

func LimitParamsFromMap(args map[string]interface{}) (LimitParams, error) {
	l := LimitParams{}

	l.Limit, l.LimitExpected = args["limitFirst"].(int)
	l.Offset, l.OffsetExpected = args["limitOffset"].(int)
	l.Min, l.MinExpected = args["limitMin"].(int)
	l.Max, l.MaxExpected = args["limitMax"].(int)

	return l, nil
}
