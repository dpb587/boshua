package contextualosmetalinkrepository

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"time"

	"github.com/cheggaaa/pb"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/git"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/metalink"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
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

func (i *index) GetCompilationArtifacts(f datastore.FilterParams) ([]compilation.Artifact, error) {
	if !f.Release.NameSatisfied(i.config.Release) {
		return nil, nil
	}

	err := i.repository.Reload()
	if err != nil {
		return nil, errors.Wrap(err, "reloading repository")
	}

	paths, err := filepath.Glob(i.repository.Path(i.config.Prefix, "*", "*", "*.meta4"))
	if err != nil {
		return nil, errors.Wrap(err, "globbing")
	}

	var results = []compilation.Artifact{}

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

		if !f.Release.VersionSatisfied(compiledReleaseMeta4.Files[0].Version) {
			continue
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

		releases, err := i.releaseVersionIndex.GetArtifacts(f.Release)
		if err != nil {
			// TODO warn and continue?
			return nil, fmt.Errorf("finding release")
		}

		release, err := releaseversiondatastore.RequireSingleResult(releases)
		if err != nil {
			// TODO warn and continue?
			return nil, fmt.Errorf("finding release")
		}

		results = append(
			results,
			compilation.New(
				compilation.Reference{
					ReleaseVersion: release.Reference().(releaseversion.Reference),
					OSVersion:      osVersionReference,
				},
				compiledReleaseMeta4.Files[0],
			),
		)
	}

	return results, nil
}

func (i *index) StoreCompilationArtifact(artifact compilation.Artifact) error {
	artifactRef := artifact.Reference().(compilation.Reference)

	logger := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(logger)

	urlLoader := urldefaultloader.New(fs)

	file := artifact.MetalinkFile()

	local, err := urlLoader.Load(metalink.URL{URL: file.URLs[0].URL})
	if err != nil {
		return errors.Wrap(err, "parsing origin destination")
	}

	var sha1 string

	for _, hash := range file.Hashes {
		if hash.Type == "sha-1" {
			sha1 = hash.Hash

			break
		}
	}

	if sha1 == "" {
		return errors.New("sha-1 hash not found")
	}

	// not a good way to inject configs
	priorAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	os.Setenv("AWS_ACCESS_KEY_ID", i.config.BlobstoreConfig.S3.AccessKey)
	defer os.Setenv("AWS_ACCESS_KEY_ID", priorAccessKey)

	priorSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	os.Setenv("AWS_SECRET_ACCESS_KEY", i.config.BlobstoreConfig.S3.SecretKey)
	defer os.Setenv("AWS_SECRET_ACCESS_KEY", priorSecretKey)

	// TODO configurable? e.g. metalink-repository-resource?
	remote, err := urlLoader.Load(metalink.URL{URL: fmt.Sprintf(
		"s3://%s/%s/%s%s/%s",
		i.config.BlobstoreConfig.S3.Host,
		i.config.BlobstoreConfig.S3.Bucket,
		i.config.BlobstoreConfig.S3.Prefix,
		sha1[0:2],
		sha1[2:],
	)})
	if err != nil {
		return errors.Wrap(err, "Parsing source blob")
	}

	progress := pb.New64(int64(file.Size)).Set(pb.Bytes, true).SetRefreshRate(time.Second).SetWidth(80)
	progress.Start()

	err = remote.WriteFrom(local, progress)
	if err != nil {
		return errors.Wrap(err, "Copying blob")
	}

	progress.Finish()

	commitMeta4 := metalink.Metalink{
		Files: []metalink.File{
			{
				Name:    file.Name, // TODO rebuild name
				Version: artifactRef.ReleaseVersion.Version,
				Size:    file.Size,
				URLs: []metalink.URL{
					{
						URL: remote.ReaderURI(),
					},
				},
				Hashes: file.Hashes,
			},
		},
		Generator: "boshua/boshreleasedpb",
	}

	commitMeta4Bytes, err := metalink.MarshalXML(commitMeta4)
	if err != nil {
		return errors.Wrap(err, "marshalling metalink")
	}

	path := filepath.Join(
		i.config.Prefix,
		artifactRef.OSVersion.Name,
		artifactRef.OSVersion.Version,
		fmt.Sprintf("v%s.meta4", artifactRef.ReleaseVersion.Version),
	)

	return i.repository.Commit(
		map[string][]byte{path: commitMeta4Bytes},
		fmt.Sprintf(
			"Compiling v%s for %s/%s",
			artifactRef.ReleaseVersion.Version,
			artifactRef.OSVersion.Name,
			artifactRef.OSVersion.Version,
		),
	)
}

func (i *index) FlushCompilationCache() error {
	// TODO defer reload?
	return i.repository.ForceReload()
}
