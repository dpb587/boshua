package pivnetfile

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
	an, bn := s[i].ProductName, s[j].ProductName
	if an == bn {
		if s[i].ReleaseVersion == s[j].ReleaseVersion {
			return strings.Compare(s[i].File.Name, s[j].File.Name) < 0
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

	return strings.Compare(an, bn) < 0
}
