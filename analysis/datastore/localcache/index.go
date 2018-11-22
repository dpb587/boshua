package localcache

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink/file/url/file"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type index struct {
	name    string
	config  Config
	logger  logrus.FieldLogger
	storage string
}

var _ datastore.Index = &index{}

func New(name string, config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		name:    name,
		config:  config,
		logger:  logger,
		storage: filepath.Join(os.Getenv("HOME"), ".cache", "boshua", "analysis-localcache"),
	}
}

func (i *index) GetName() string {
	return i.name
}

func (i *index) GetAnalysisArtifacts(ref analysis.Reference) ([]analysis.Artifact, error) {
	cachePath, err := i.cachePath(ref)
	if err != nil {
		return nil, errors.Wrap(err, "generating cache path")
	}

	stat, err := os.Stat(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "checking analysis file")
	}

	return []analysis.Artifact{
		analysis.New(
			i.name,
			ref,
			metalink.File{
				Name: fmt.Sprintf("%s.jsonl.gz", ref.Analyzer),
				Size: uint64(stat.Size()),
				URLs: []metalink.URL{
					{
						URL: fmt.Sprintf("file://%s", cachePath),
					},
				},
			},
		),
	}, nil
}

func (i *index) StoreAnalysisResult(ref analysis.Reference, source metalink.Metalink) error {
	cachePath, err := i.cachePath(ref)
	if err != nil {
		return errors.Wrap(err, "generating cache path")
	}

	err = os.MkdirAll(filepath.Dir(cachePath), 0750)
	if err != nil {
		return errors.Wrap(err, "creating directories")
	}

	fh, err := os.OpenFile(cachePath, os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return errors.Wrap(err, "opening file for writing")
	}

	defer fh.Close()

	local := file.NewReference(source.Files[0].URLs[0].URL)

	reader, err := local.Reader()
	if err != nil {
		return errors.Wrap(err, "open reader")
	}

	_, err = io.Copy(fh, reader)
	if err != nil {
		return errors.Wrap(err, "saving analysis")
	}

	return nil
}

func (i *index) FlushAnalysisCache() error {
	return nil
}

func (i *index) cachePath(ref analysis.Reference) (string, error) {
	subjectRef := ref.Subject.Reference()

	pieces := []string{string(ref.Analyzer)}

	switch subjectRef := subjectRef.(type) {
	case compilation.Reference:
		pieces = append(
			pieces,
			"compiled-release",
			subjectRef.OSVersion.Name,
			subjectRef.OSVersion.Version,
			subjectRef.ReleaseVersion.UniqueID(),
		)
	case releaseversion.Reference:
		pieces = append(
			pieces,
			"release",
			subjectRef.UniqueID(),
		)
	case stemcellversion.Reference:
		pieces = append(
			pieces,
			"stemcell",
			subjectRef.UniqueID(),
		)
	default:
		return "", errors.New("unsupported analysis subject")
	}

	hash := sha1.New()
	hash.Write([]byte(strings.Join(pieces, "\n")))

	key := fmt.Sprintf("%x", hash.Sum(nil))

	return filepath.Join(i.storage, key[0:2], key[2:]), nil
}
