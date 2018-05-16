package localcache

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/metalink"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
	"github.com/pkg/errors"
)

type index struct {
	storage string
}

var _ datastore.Index = &index{}

func New() datastore.Index {
	return &index{
		storage: filepath.Join(os.Getenv("HOME"), ".cache", "boshua", "analysis-localcache"),
	}
}

func (i *index) Filter(ref analysis.Reference) ([]analysis.Artifact, error) {
	cachePath, err := i.cachePath(ref)
	if err != nil {
		return nil, errors.Wrap(err, "generating cache path")
	}

	stat, err := os.Stat(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, datastore.NoMatchErr
		}

		return nil, errors.Wrap(err, "checking analysis file")
	}

	return []analysis.Artifact{
		analysis.New(
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

func (i *index) Store(ref analysis.Reference, source metalink.Metalink) error {
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

	logger := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(logger)

	urlLoader := urldefaultloader.New(fs)

	file := source.Files[0]

	local, err := urlLoader.Load(metalink.URL{URL: file.URLs[0].URL})
	if err != nil {
		return errors.Wrap(err, "parsing origin destination")
	}

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

func (i *index) cachePath(ref analysis.Reference) (string, error) {
	subjectRef := ref.Subject.Reference()

	pieces := []string{string(ref.Analyzer)}

	switch subjectRef := subjectRef.(type) {
	case compiledreleaseversion.Reference:
		pieces = append(
			pieces,
			"compiled-release",
			subjectRef.OSVersion.Name,
			subjectRef.OSVersion.Version,
			subjectRef.ReleaseVersion.Name,
			subjectRef.ReleaseVersion.Version,
			subjectRef.ReleaseVersion.Checksums.Preferred().String(),
		)
	case releaseversion.Reference:
		pieces = append(
			pieces,
			"release",
			subjectRef.Name,
			subjectRef.Version,
			subjectRef.Checksums.Preferred().String(),
		)
	default:
		return "", errors.New("unsupported analysis subject")
	}

	hash := sha1.New()
	hash.Write([]byte(strings.Join(pieces, "\n")))

	key := fmt.Sprintf("%x", hash.Sum(nil))

	return filepath.Join(i.storage, key[0:2], key[2:]), nil
}
