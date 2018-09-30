package pivnet

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/cheggaaa/pb"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"
	"github.com/dpb587/metalink/template"
	"github.com/dpb587/boshua/pivnetfile"
	"github.com/dpb587/metalink"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type index struct {
	name       string
	logger     logrus.FieldLogger
	config     Config
	repository *repository.Repository
}

var _ datastore.Index = &index{}

func New(name string, config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		name:       name,
		logger:     logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config:     config,
		repository: repository.NewRepository(logger, config.RepositoryConfig),
	}
}

func (i *index) GetName() string {
	return i.name
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
		i.logger.Errorf("failed to unmarshal %s", analysisPath)

		return nil, nil
	}

	return []analysis.Artifact{
		analysis.New(i.name, ref, analysisMeta4.Files[0]),
	}, nil
}

func (i *index) FlushAnalysisCache() error {
	// TODO defer reload?
	return i.repository.ForceReload()
}

func (i *index) storagePath(ref analysis.Reference) (string, error) {
	subjectRef := ref.Subject.Reference()

	switch subjectRef := subjectRef.(type) {
	case pivnetfile.Reference:
		return filepath.Join(
			i.config.RepositoryPath,
			subjectRef.ProductName,
			strconv.Itoa(subjectRef.ReleaseID),
			fmt.Sprintf("%s.%s.meta4", strconv.Itoa(subjectRef.FileID), ref.Analyzer),
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

	mirroredFile := metalink.File{
		Name:   fmt.Sprintf("%s.jsonl", ref.Analyzer),
		Size:   file.Size,
		Hashes: file.Hashes,
	}

	for _, mirror := range i.config.StorageConfig {
		tmpl, err := template.New(mirror.URI)
		if err != nil {
			return errors.Wrap(err, "parsing mirror destination")
		}

		mirrorWriterURI, err := tmpl.ExecuteString(mirroredFile)
		if err != nil {
			return errors.Wrap(err, "generating mirror uri")
		}

		for k, v := range mirror.Options {
			// TODO unset/revert after?
			os.Setenv(k, v)
		}

		remote, err := urlLoader.Load(metalink.URL{URL: mirrorWriterURI})
		if err != nil {
			return errors.Wrap(err, "loading upload destination")
		}

		progress := pb.New64(int64(file.Size)).Set(pb.Bytes, true).SetRefreshRate(time.Second).SetWidth(80)
		progress.Start()

		err = remote.WriteFrom(local, progress)
		if err != nil {
			return errors.Wrap(err, "copying blob")
		}

		progress.Finish()

		mirroredFile.URLs = append(mirroredFile.URLs, metalink.URL{
			URL:      remote.ReaderURI(),
			Location: mirror.Location,
			Priority: mirror.Priority,
		})
	}

	commitMeta4 := metalink.Metalink{
		Files:     []metalink.File{mirroredFile},
		Generator: "boshua/pivnet",
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
