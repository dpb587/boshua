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
	var results = []pivnetfile.Artifact{}

	products, err := i.getProducts(f)
	if err != nil {
		return nil, errors.Wrap(err, "getting products")
	}

	for _, product := range products {
		releases, err := i.getReleases(product, f)
		if err != nil {
			return nil, errors.Wrap(err, "getting releases")
		}

		for _, release := range releases {
			files, err := i.getProductFiles(product, release, f)
			if err != nil {
				return nil, errors.Wrap(err, "getting product files")
			}

			for _, file := range files {
				metalinkFile := metalink.File{
					Name: path.Base(file.AWSObjectKey), // TODO weird; correct?
					Size: uint64(file.Size),
				}

				if file.MD5 != "" {
					metalinkFile .Hashes = append(metalinkFile.Hashes, metalink.Hash{
						Type: "md5",
						Hash: file.MD5,
					})
				}

				if file.SHA256 != "" {
					metalinkFile.Hashes = append(metalinkFile.Hashes, metalink.Hash{
						Type: "sha-256",
						Hash: file.SHA256,
					})
				}

				downloadURL, err := file.DownloadLink()
				if err == nil {
					metalinkFile.URLs = append(metalinkFile.URLs, metalink.URL{
						URL: strings.Replace(downloadURL, "https://", "pivnet://", 1),
					})
				}

				results = append(
					results,
					pivnetfile.Artifact{
						Datastore:      i.name,
						ProductName:    product.Slug,
						ReleaseID:      release.ID,
						ReleaseVersion: release.Version,
						FileID:         file.ID,
						File:           metalinkFile,
					},
				)
			}
		}
	}

	return results, nil
}

func (i *index) getProducts(f datastore.FilterParams) ([]pivnet.Product, error) {
	if f.ProductNameExpected {
		return []pivnet.Product{
			{
				Slug: f.ProductName,
			},
		}, nil
	}

	return nil, errors.New("product slug is required")
}

func (i *index) getReleases(product pivnet.Product, f datastore.FilterParams) ([]pivnet.Release, error) {
	if f.ReleaseIDExpected {
		return []pivnet.Release{
			{
				ID: f.ReleaseID,
			},
		}, nil
	}

	return i.client.Releases.List(product.Slug)
}

func (i *index) getProductFiles(product pivnet.Product, release pivnet.Release, f datastore.FilterParams) ([]pivnet.ProductFile, error) {
	if f.FileIDExpected {
		file, err := i.client.ProductFiles.GetForRelease(product.Slug, release.ID, f.FileID)
		if err != nil {
			return nil, err
		}

		return []pivnet.ProductFile{file}, nil
	}

	return i.client.ProductFiles.ListForRelease(product.Slug, release.ID)
}

func (i *index) FlushCache() error {
	// TODO mandatory interface?
	return nil
}
