package pivnet

import (
	"io/ioutil"
	"log"
	"path"
	"reflect"
	"strings"

	"github.com/dpb587/boshua/pivnetfile"
	"github.com/dpb587/boshua/pivnetfile/datastore"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logshim"
)

type index struct {
	name   string
	logger logrus.FieldLogger
	client pivnet.Client
}

var _ datastore.Index = &index{}

func New(name string, config Config, logger logrus.FieldLogger) datastore.Index {
	clientConfig := pivnet.ClientConfig{
		Host:      config.Host,
		Token:     config.Token,
		UserAgent: "boshua/0.0.0+dev", // TODO import app main version
	}

	// TODO wrap for logrus?
	stdoutLogger := log.New(ioutil.Discard, "", log.LstdFlags)
	stderrLogger := log.New(ioutil.Discard, "", log.LstdFlags)

	verbose := false
	pivnetLogger := logshim.NewLogShim(stdoutLogger, stderrLogger, verbose)

	return &index{
		name:       name,
		logger:     logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		client:     pivnet.NewClient(clientConfig, pivnetLogger),
	}
}

func (i *index) GetName() string {
	return i.name
}

func (i *index) GetArtifacts(f datastore.FilterParams) ([]pivnetfile.Artifact, error) {
	if !f.ProductNameExpected || !f.ReleaseIDExpected || !f.FileIDExpected {
		return nil, errors.New("product name, release id, file id are currently required")
	}

	var results = []pivnetfile.Artifact{}

	found, err := i.client.ProductFiles.GetForRelease(f.ProductName, f.ReleaseID, f.FileID)
	if err != nil {
		// TODO catch 404 not found for safe, empty artifact results
		return nil, errors.Wrap(err, "finding pivnet file")
	}

	// TODO separate pivnet api call to load release metadata?

	file := metalink.File{
		Name: path.Base(found.AWSObjectKey), // TODO weird; correct?
		Size: uint64(found.Size),
	}

	if found.MD5 != "" {
		file.Hashes = append(file.Hashes, metalink.Hash{
			Type: "md5",
			Hash: found.MD5,
		})
	}

	if found.SHA256 != "" {
		file.Hashes = append(file.Hashes, metalink.Hash{
			Type: "sha-256",
			Hash: found.SHA256,
		})
	}

	downloadURL, err := found.DownloadLink()
	if err == nil {
		file.URLs = append(file.URLs, metalink.URL{
			URL: strings.Replace(downloadURL, "https://", "pivnet://", 1),
		})
	}

	results = append(
		results,
		pivnetfile.Artifact{
			Datastore:   i.name,
			ProductName: f.ProductName,
			ReleaseID:   f.ReleaseID,
			FileID:      found.ID,
			File:        file,
		},
	)

	return results, nil
}

func (i *index) FlushCache() error {
	// TODO mandatory interface?
	return nil
}
