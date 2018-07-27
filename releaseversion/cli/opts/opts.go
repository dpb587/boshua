package opts

import (
	"fmt"

	"github.com/dpb587/boshua/cli/args"
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/pkg/errors"
)

type Opts struct {
	AppOpts *cmdopts.Opts `no-flag:"true"`

	NameVersion *Release       `long:"release" description:"The release in name/version format"`
	Name        string         `long:"release-name" description:"The release name"`
	Version     string         `long:"release-version" description:"The release version"`
	Checksum    *args.Checksum `long:"release-checksum" description:"The release checksum"`
	URI         string         `long:"release-url" description:"The release source URL"`

	Labels []string `long:"release-label" description:"The label(s) to filter releases by"`
}

func ArgsFromFilterParams(f *datastore.FilterParams) []string {
	args := []string{}

	if f.NameExpected {
		args = append(args, fmt.Sprintf("--release-name=%s", f.Name))
	}

	if f.VersionExpected {
		args = append(args, fmt.Sprintf("--release-version=%s", f.Version))
	}

	if f.ChecksumExpected {
		args = append(args, fmt.Sprintf("--release-checksum=%s", f.Checksum))
	}

	if f.URIExpected {
		args = append(args, fmt.Sprintf("--release-url=%s", f.URI))
	}

	if f.LabelsExpected {
		for _, label := range f.Labels {
			args = append(args, fmt.Sprintf("--release-label=%s", label))
		}
	}

	return args
}

func (o *Opts) Artifact() (releaseversion.Artifact, error) {
	index, err := o.AppOpts.GetReleaseIndex("default")
	if err != nil {
		return releaseversion.Artifact{}, errors.Wrap(err, "loading release index")
	}

	results, err := index.Filter(o.FilterParams())
	if err != nil {
		return releaseversion.Artifact{}, errors.Wrap(err, "finding release")
	}

	result, err := datastore.RequireSingleResult(results)
	if err != nil {
		return releaseversion.Artifact{}, errors.Wrap(err, "finding release")
	}

	return result, err
}

func (o Opts) FilterParams() *datastore.FilterParams {
	f := &datastore.FilterParams{
		LabelsExpected: len(o.Labels) > 0,
		Labels:         o.Labels,
	}

	if o.NameVersion != nil {
		if o.Name != "" || o.Version != "" {
			// TODO not panic
			panic("cannot specify both --release and one of --release-name or --release-version")
		}

		f.NameExpected = true
		f.Name = o.NameVersion.Name

		f.VersionExpected = true
		f.Version = o.NameVersion.Version
	} else {
		f.NameExpected = o.Name != ""
		f.Name = o.Name

		f.VersionExpected = o.Version != ""
		f.Version = o.Version
	}

	f.URIExpected = o.URI != ""
	f.URI = o.URI

	if o.Checksum != nil {
		f.ChecksumExpected = true
		f.Checksum = o.Checksum.ImmutableChecksum.String()
	}

	return f
}
