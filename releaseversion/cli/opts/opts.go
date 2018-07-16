package opts

import (
	"github.com/dpb587/boshua/cli/args"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
)

type Opts struct {
	Name     string         `long:"release-name" description:"The release name"`
	Version  string         `long:"release-version" description:"The release version"`
	Checksum *args.Checksum `long:"release-checksum" description:"The release checksum"`
}

func (o Opts) Reference() releaseversion.Reference {
	ref := releaseversion.Reference{
		Name:    o.Name,
		Version: o.Version,
	}

	if o.Checksum != nil {
		ref.Checksums = append(ref.Checksums, o.Checksum.ImmutableChecksum)
	}

	return ref
}

func (o Opts) FilterParams() *datastore.FilterParams {
	f := &datastore.FilterParams{}

	if o.Name != "" {
		f.NameExpected = true
		f.Name = o.Name
	}

	if o.Version != "" {
		f.VersionExpected = true
		f.Version = o.Version
	}

	if o.Checksum != nil {
		f.ChecksumExpected = true
		f.Checksum = o.Checksum.ImmutableChecksum.String()
	}

	return f
}
