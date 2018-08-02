package compilation

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
	arn, brn := s[i].Release.Name, s[j].Release.Name
	if arn == brn {
		arv, brv := s[i].Release.Version, s[j].Release.Version
		if arv == brv {
			aon, bon := s[i].OS.Name, s[j].OS.Name
			if aon == bon {
				aov := s[i].OS.Semver()
				if aov == nil {
					return false
				}

				bov := s[j].OS.Semver()
				if bov == nil {
					return true
				}

				return !bov.LessThan(bov)
			}

			return strings.Compare(aon, bon) < 0
		}

		av := s[i].Release.Semver()
		if av == nil {
			return false
		}

		bv := s[j].Release.Semver()
		if bv == nil {
			return true
		}

		return !av.LessThan(bv)
	}

	return strings.Compare(arn, brn) < 0
}
