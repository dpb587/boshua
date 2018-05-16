package boshreleasedpb

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/cheggaaa/pb"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/metalink"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
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

func (i *index) Filter(ref analysis.Reference) ([]analysis.Artifact, error) {
	return nil, errors.New("TODO")
}

func (i *index) Find(ref analysis.Reference) (analysis.Artifact, error) {
	return datastore.FilterForOne(i, ref)
}

func (i *index) storagePath(ref analysis.Reference) (string, error) {
	subjectRef := ref.Subject.Reference()

	switch subjectRef := subjectRef.(type) {
	case compiledreleaseversion.Reference:
		return filepath.Join(
			"compiled-release",
			i.config.Channel,
			subjectRef.OSVersion.Name,
			subjectRef.OSVersion.Version,
			"analysis",
			string(ref.Analyzer),
			fmt.Sprintf("%s.meta4", subjectRef.ReleaseVersion.Version),
		), nil
	case releaseversion.Reference:
		return filepath.Join(
			"release",
			i.config.Channel,
			"analysis",
			string(ref.Analyzer),
			fmt.Sprintf("%s.meta4", subjectRef.Version),
		), nil
	}

	return "", datastore.UnsupportedOperationErr
}

func (i *index) Store(ref analysis.Reference, artifactMeta4 metalink.Metalink) error {
	config, err := i.loadConfig()
	if err != nil {
		return errors.Wrap(err, "loading release config")
	}

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
	os.Setenv("AWS_ACCESS_KEY_ID", config.Blobstore.Options.SecretAccessKey)
	defer os.Setenv("AWS_ACCESS_KEY_ID", priorAccessKey)

	priorSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	os.Setenv("AWS_SECRET_ACCESS_KEY", config.Blobstore.Options.AccessKeyID)
	defer os.Setenv("AWS_SECRET_ACCESS_KEY", priorSecretKey)

	remote, err := urlLoader.Load(metalink.URL{URL: fmt.Sprintf(
		"s3://%s/%s/analysis/%s/%s",
		config.Blobstore.Options.Host,
		config.Blobstore.Options.BucketName,
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
			},
		},
		Generator: "boshua/boshreleasedpb",
	}

	commitMeta4Bytes, err := metalink.Marshal(commitMeta4)
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

func (i *index) loadConfig() (releaseConfig, error) {
	var finalConfig, privateConfig releaseConfig

	{ // final.yml
		finalBytes, err := ioutil.ReadFile(filepath.Join(i.config.LocalPath, "config", "final.yml"))
		if err != nil {
			return releaseConfig{}, errors.Wrap(err, "reading final.yml")
		}

		err = yaml.Unmarshal(finalBytes, &finalConfig)
		if err != nil {
			return releaseConfig{}, errors.Wrap(err, "unmarshalling final.yml")
		}
	}

	{ // private.yml
		privateBytes, err := ioutil.ReadFile(filepath.Join(i.config.LocalPath, "config", "private.yml"))
		if err != nil {
			return releaseConfig{}, errors.Wrap(err, "reading private.yml")
		}

		err = yaml.Unmarshal(privateBytes, &privateConfig)
		if err != nil {
			return releaseConfig{}, errors.Wrap(err, "unmarshalling private.yml")
		}
	}

	finalConfig.Merge(privateConfig)

	return finalConfig, nil
}
