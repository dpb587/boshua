package scheduler

import (
	"bcr-server/releaseversions"
	"bcr-server/stemcellversions"
	"encoding/json"
	"fmt"
)

func Plan(release releaseversions.ReleaseVersion, stemcell stemcellversions.StemcellVersion) {
	bytes, _ := json.Marshal(map[string]interface{}{
		"release_source":  release.MetalinkSource,
		"stemcell_source": stemcell.MetalinkSource,
	})

	fmt.Printf("%s\n", bytes)
}
