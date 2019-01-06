package datastore

type LimitParams struct {
	LimitExpected bool
	Limit         int

	OffsetExpected bool
	Offset         int

	MinExpected bool
	Min         int

	MaxExpected bool
	Max         int
}

var SingleArtifactLimitParams = LimitParams{
	MinExpected:   true,
	Min:           1,
	LimitExpected: true,
	Limit:         1,
}
