package releaseversion

import (
	"sort"
)

func Sort(results []Artifact) {
	x := sorter(results)
	sort.Sort(&x)
}

type sorter []Artifact

func (s sorter) Len() int {
	return len(s)
}

func (s sorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sorter) Less(i, j int) bool {
	av := s[i].Semver()
	if av == nil {
		return false
	}

	bv := s[j].Semver()
	if bv == nil {
		return true
	}

	return !av.LessThan(bv)
}
