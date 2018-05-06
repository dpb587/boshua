package datastore

import "github.com/dpb587/boshua/compiledreleaseversion"

func FilterForOne(index Index, ref compiledreleaseversion.Reference) (compiledreleaseversion.Artifact, error) {
	results, err := index.Filter(ref)
	if err != nil {
		return compiledreleaseversion.Artifact{}, err
	} else if len(results) == 0 {
		return compiledreleaseversion.Artifact{}, NoMatchErr
	} else if len(results) > 1 {
		return compiledreleaseversion.Artifact{}, MultipleMatchErr
	}

	return results[0], nil
}
