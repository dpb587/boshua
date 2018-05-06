package datastore

import "github.com/dpb587/boshua/osversion"

func FilterForOne(index Index, ref osversion.Reference) (osversion.Artifact, error) {
	results, err := index.Filter(ref)
	if err != nil {
		return osversion.Artifact{}, err
	} else if len(results) == 0 {
		return osversion.Artifact{}, NoMatchErr
	} else if len(results) > 1 {
		return osversion.Artifact{}, MultipleMatchErr
	}

	return results[0], nil
}
