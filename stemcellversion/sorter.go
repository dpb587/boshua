package stemcellversion

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
	ao, bo := s[i].OS, s[j].OS
	if ao == bo {
		if s[i].Version == s[j].Version {
			return strings.Compare(s[i].Flavor, s[j].Flavor) < 0
		}

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

	return strings.Compare(ao, bo) < 0
}
