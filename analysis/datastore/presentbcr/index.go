package presentbcr

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type index struct {
	logger             logrus.FieldLogger
	metalinkRepository string
	localPath          string
	pullInterval       time.Duration

	inmemory   datastore.Index
	lastLoaded time.Time
}

var _ datastore.Index = &index{}

func New(config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		logger:             logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		metalinkRepository: config.Repository,
		localPath:          config.LocalPath,
		pullInterval:       config.PullInterval,
	}
}

func (i *index) Filter(ref analysis.Reference) ([]analysis.Artifact, error) {
	_, err := i.reloader()
	if err != nil {
		return nil, errors.Wrap(err, "reloading")
	}

	meta4Path := filepath.Join(i.localPath, ref.ArtifactStorageDir(), "artifact.meta4")

	meta4Bytes, err := ioutil.ReadFile(meta4Path)
	if err != nil {
		if os.IsNotExist(err) {
			return []analysis.Artifact{}, nil
		}

		return nil, fmt.Errorf("reading %s: %v", meta4Path, err)
	}

	var meta4 metalink.Metalink

	err = metalink.Unmarshal(meta4Bytes, &meta4)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling %s: %v", meta4Path, err)
	}

	return []analysis.Artifact{
		analysis.New(
			ref.Artifact,
			ref.Analyzer,
			meta4.Files[0],
			map[string]interface{}{
				"uri": fmt.Sprintf("%s//%s", i.metalinkRepository, strings.TrimPrefix(path.Dir(strings.TrimPrefix(meta4Path, i.localPath)), "/")),
			},
		),
	}, nil
}

func (i *index) Find(ref analysis.Reference) (analysis.Artifact, error) {
	return datastore.FilterForOne(i, ref)
}

func (i *index) reloader() (bool, error) {
	if time.Now().Sub(i.lastLoaded) < i.pullInterval {
		return false, nil
	} else if !strings.HasPrefix(i.metalinkRepository, "git+") {
		return false, nil
	}

	i.lastLoaded = time.Now()

	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = i.localPath

	outbuf := bytes.NewBuffer(nil)
	errbuf := bytes.NewBuffer(nil)

	cmd.Stdout = outbuf
	cmd.Stderr = errbuf

	err := cmd.Run()
	if err != nil {
		i.logger.WithField("error", err).Errorf("pulling repository")

		return false, errors.Wrap(err, "pulling repository")
	}

	if strings.Contains(outbuf.String(), "Already up to date.") {
		i.logger.Debugf("repository already up to date")

		return false, nil
	}

	i.logger.Debugf("repository updated")

	return true, nil
}
