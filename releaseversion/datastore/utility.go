package datastore

import "github.com/dpb587/boshua/releaseversion"

func FilterForOne(index Index, ref releaseversion.Reference) (releaseversion.Artifact, error) {
	results, err := index.Filter(ref)
	if err != nil {
		return releaseversion.Artifact{}, err
	} else if len(results) == 0 {
		return releaseversion.Artifact{}, NoMatchErr
	} else if len(results) > 1 {
		return releaseversion.Artifact{}, MultipleMatchErr
	}

	return results[0], nil
}
