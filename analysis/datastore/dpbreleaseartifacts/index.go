package boshreleasedpb

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/git"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/metalink"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type index struct {
	logger     logrus.FieldLogger
	config     Config
	repository *git.Repository
}

var _ datastore.Index = &index{}

func New(config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		logger:     logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config:     config,
		repository: git.NewRepository(logger, config.RepositoryConfig),
	}
}

func (i *index) GetAnalysisArtifacts(ref analysis.Reference) ([]analysis.Artifact, error) {
	err := i.repository.Reload()
	if err != nil {
		return nil, errors.Wrap(err, "reloading repository")
	}

	analysisPath, err := i.storagePath(ref)
	if err != nil {
		return nil, errors.Wrap(err, "finding analysis path")
	}

	analysisBytes, err := ioutil.ReadFile(i.repository.Path(analysisPath))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "reading analysis meta4")
	}

	var analysisMeta4 metalink.Metalink

	err = metalink.Unmarshal(analysisBytes, &analysisMeta4)
	if err != nil {
		// TODO warn and continue?
		return nil, fmt.Errorf("unmarshalling %s: %v", analysisPath, err)
	}

	return []analysis.Artifact{
		analysis.New(ref, analysisMeta4.Files[0]),
	}, nil
}

func (i *index) FlushAnalysisCache() error {
	// TODO defer reload?
	return i.repository.ForceReload()
}

func (i *index) storagePath(ref analysis.Reference) (string, error) {
	subjectRef := ref.Subject.Reference()

	switch subjectRef := subjectRef.(type) {
	case compilation.Reference:
		return filepath.Join(
			i.config.CompiledReleasePrefix,
			subjectRef.OSVersion.Name,
			subjectRef.OSVersion.Version,
			"analysis",
			string(ref.Analyzer),
			fmt.Sprintf("%s.meta4", subjectRef.ReleaseVersion.Version),
		), nil
	case releaseversion.Reference:
		return filepath.Join(
			i.config.ReleasePrefix,
			"analysis",
			string(ref.Analyzer),
			fmt.Sprintf("%s.meta4", subjectRef.Version),
		), nil
	case stemcellversion.Reference:
		return filepath.Join(
			i.config.StemcellPrefix,
			"analysis",
			subjectRef.OS,
			subjectRef.Version,
			strings.TrimSuffix(fmt.Sprintf("%s-%s-%s", subjectRef.IaaS, subjectRef.Hypervisor, subjectRef.DiskFormat), "-"),
			fmt.Sprintf("%s.%s.meta4", ref.Analyzer, subjectRef.Flavor),
		), nil
	}

	return "", datastore.UnsupportedOperationErr
}

func (i *index) StoreAnalysisResult(ref analysis.Reference, artifactMeta4 metalink.Metalink) error {
	logger := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(logger)

	urlLoader := urldefaultloader.New(fs)

	path, err := i.storagePath(ref)
	if err != nil {
		return errors.Wrap(err, "generating path")
	}

	file := artifactMeta4.Files[0]

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

	// TODO configurable?
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
				Name: fmt.Sprintf("%s.jsonl", ref.Analyzer),
				Size: file.Size,
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

	return i.repository.Commit(
		map[string][]byte{path: commitMeta4Bytes},
		fmt.Sprintf(
			"Add %s analysis",
			ref.Analyzer,
		),
	)
}
