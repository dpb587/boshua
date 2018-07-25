package releaseversion

import (
	"sort"
	"strings"
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
	an, bn := s[i].Name, s[j].Name
	if an == bn {
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

	return strings.Compare(an, bn) < 0
}
