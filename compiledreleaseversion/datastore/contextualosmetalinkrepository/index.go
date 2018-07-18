package contextualosmetalinkrepository

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"

	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ./{os_name}/{os_version}/*.meta4
type index struct {
	logger              logrus.FieldLogger
	config              Config
	repository          *git.Repository
	releaseVersionIndex releaseversiondatastore.Index
}

var _ datastore.Index = &index{}

func New(releaseVersionIndex releaseversiondatastore.Index, config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		logger:              logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		releaseVersionIndex: releaseVersionIndex,
		config:              config,
		repository:          git.NewRepository(logger, config.RepositoryConfig),
	}
}

func (i *index) Filter(f *datastore.FilterParams) ([]compiledreleaseversion.Artifact, error) {
	paths, err := filepath.Glob(filepath.Join(i.config.RepositoryConfig.LocalPath, "*", "*", "*.meta4"))
	if err != nil {
		return nil, errors.Wrap(err, "globbing")
	}

	var results = []compiledreleaseversion.Artifact{}

	for _, compiledReleasePath := range paths {
		compiledReleaseBytes, err := ioutil.ReadFile(compiledReleasePath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", compiledReleasePath, err)
		}

		var compiledReleaseMeta4 metalink.Metalink

		err = metalink.Unmarshal(compiledReleaseBytes, &compiledReleaseMeta4)
		if err != nil {
			// TODO warn and continue?
			return nil, fmt.Errorf("unmarshalling %s: %v", compiledReleasePath, err)
		}

		releases, err := i.releaseVersionIndex.Filter(f.Release)
		if err != nil {
			// TODO warn and continue?
			return nil, fmt.Errorf("finding release")
		}

		release, err := releaseversiondatastore.RequireSingleResult(releases)
		if err != nil {
			// TODO warn and continue?
			return nil, fmt.Errorf("finding release")
		}

		// TODO sanity checks? files = 1?

		osVersionReference := osversion.Reference{
			Name:    path.Base(path.Dir(path.Dir(compiledReleasePath))),
			Version: path.Base(path.Dir(compiledReleasePath)),
		}

		if !f.OS.NameSatisfied(osVersionReference.Name) {
			continue
		} else if !f.OS.VersionSatisfied(osVersionReference.Version) {
			continue
		}

		results = append(
			results,
			compiledreleaseversion.New(
				compiledreleaseversion.Reference{
					ReleaseVersion: release.Reference().(releaseversion.Reference),
					OSVersion:      osVersionReference,
				},
				compiledReleaseMeta4.Files[0],
			),
		)
	}

	return results, nil
}

func (i *index) Store(artifact compiledreleaseversion.Artifact) error {
	artifactRef := artifact.Reference().(compiledreleaseversion.Reference)

	path := filepath.Join(
		artifactRef.OSVersion.Name,
		artifactRef.OSVersion.Version,
		fmt.Sprintf("%s.meta4", artifactRef.ReleaseVersion.Version),
	)

	meta4 := metalink.Metalink{
		Files:     []metalink.File{artifact.MetalinkFile()},
		Generator: "boshua/contextualosmetalinkrepository",
	}

	meta4Bytes, err := metalink.MarshalXML(meta4)
	if err != nil {
		return errors.Wrap(err, "marshalling metalink")
	}

	return i.repository.Commit(
		map[string][]byte{path: meta4Bytes},
		fmt.Sprintf(
			"Compiling v%s for %s/%s",
			artifactRef.ReleaseVersion.Version,
			artifactRef.OSVersion.Name,
			artifactRef.OSVersion.Version,
		),
	)
}
