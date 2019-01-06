package opts

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/dpb587/boshua/cli/args"
	"github.com/dpb587/boshua/config"
	"github.com/dpb587/boshua/config/provider"
	cmdopts "github.com/dpb587/boshua/main/boshua/cmd/opts"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/util/semverutil"
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

func (o *Opts) Artifact(cfg *provider.Config) (releaseversion.Artifact, error) {
	index, err := cfg.GetReleaseIndex(config.DefaultName)
	if err != nil {
		return releaseversion.Artifact{}, errors.Wrap(err, "loading release index")
	}

	f, l := o.ArtifactParams()
	l.MinExpected = true
	l.Min = 1
	l.LimitExpected = true
	l.Limit = 1

	results, err := index.GetArtifacts(f, l)
	if err != nil {
		return releaseversion.Artifact{}, errors.Wrap(err, "finding release")
	}

	return results[0], nil
}

func (o Opts) ArtifactParams() (datastore.FilterParams, datastore.LimitParams) {
	f := datastore.FilterParams{
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

	l := datastore.LimitParams{}

	if f.VersionExpected {
		l.MinExpected = true
		l.Min = 1

		if f.Version == "latest" {
			f.VersionExpected = false
			f.Version = ""
			l.LimitExpected = true
			l.Limit = 1
		} else if strings.HasSuffix(f.Version, ".latest") {
			f.Version = fmt.Sprintf("%s.x", strings.TrimSuffix(f.Version, ".latest"))
			l.LimitExpected = true
			l.Limit = 1
		}
	}

	if f.VersionExpected && semverutil.IsConstraint(f.Version) {
		// ignoring errors since it can fallback to literal match
		f.VersionConstraint, _ = semver.NewConstraint(f.Version)
	}

	return f, l
}
