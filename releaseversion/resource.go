package releaseversion

func (s Subject) Resource() map[string]interface{} {
	return s.MetalinkSource
}
