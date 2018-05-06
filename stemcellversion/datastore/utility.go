package datastore

import "github.com/dpb587/boshua/stemcellversion"

func FilterForOne(index Index, ref stemcellversion.Reference) (stemcellversion.Artifact, error) {
	results, err := index.Filter(ref)
	if err != nil {
		return stemcellversion.Artifact{}, err
	} else if len(results) == 0 {
		return stemcellversion.Artifact{}, NoMatchErr
	} else if len(results) > 1 {
		return stemcellversion.Artifact{}, MultipleMatchErr
	}

	return results[0], nil
}
