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
	arn, brn := s[i].reference.ReleaseVersion.Name, s[j].reference.ReleaseVersion.Name
	if arn == brn {
		arv, brv := s[i].reference.ReleaseVersion.Version, s[j].reference.ReleaseVersion.Version
		if arv == brv {
			aon, bon := s[i].reference.OSVersion.Name, s[j].reference.OSVersion.Name
			if aon == bon {
				aov := s[i].reference.OSVersion.Semver()
				if aov == nil {
					return false
				}

				bov := s[j].reference.OSVersion.Semver()
				if bov == nil {
					return true
				}

				return !bov.LessThan(bov)
			}

			return strings.Compare(aon, bon) < 0
		}

		av := s[i].reference.ReleaseVersion.Semver()
		if av == nil {
			return false
		}

		bv := s[j].reference.ReleaseVersion.Semver()
		if bv == nil {
			return true
		}

		return !av.LessThan(bv)
	}

	return strings.Compare(arn, brn) < 0
}
