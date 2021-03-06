package boshreleasedir

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"path"
	"path/filepath"
	"reflect"
	"regexp"

	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"
	"github.com/dpb587/boshua/metalink/file/metaurl/boshreleasesource"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore/inmemory"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
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

func (i *index) GetArtifacts(f datastore.FilterParams, l datastore.LimitParams) ([]releaseversion.Artifact, error) {
	if f.ChecksumExpected {
		return nil, nil
	} else if !f.LabelsSatisfied(i.config.Labels) {
		return nil, nil
	}

	err := i.repository.Reload()
	if err != nil {
		return nil, errors.Wrap(err, "reloading repository")
	}

	var globReleaseName = "*"
	if f.NameExpected {
		globReleaseName = filepath.Base(f.Name)
	}

	indices, err := filepath.Glob(i.repository.Path("releases", globReleaseName, "index.yml"))
	if err != nil {
		return nil, errors.Wrap(err, "globbing")
	}

	var results = []releaseversion.Artifact{}

	for _, indexPath := range indices {
		indexBytes, err := ioutil.ReadFile(indexPath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", indexPath, err)
		}

		var index boshReleaseIndex

		err = yaml.Unmarshal(indexBytes, &index)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling %s: %v", indexPath, err)
		}

		for _, build := range index.Builds {
			releaseName := path.Base(path.Dir(indexPath))
			releaseSubPath := filepath.Join("releases", releaseName, fmt.Sprintf("%s-%s.yml", releaseName, build.Version))
			releasePath := i.repository.Path(releaseSubPath)

			// TODO track with checksum of release manifest somewhere?
			releaseBytes, err := ioutil.ReadFile(releasePath)
			if err != nil {
				return nil, fmt.Errorf("reading %s: %v", releasePath, err)
			}

			var release boshRelease

			err = yaml.Unmarshal(releaseBytes, &release)
			if err != nil {
				return nil, fmt.Errorf("unmarshalling %s: %v", indexPath, err)
			}

			if !f.NameSatisfied(release.Name) {
				continue
			} else if !f.VersionSatisfied(release.Version) {
				continue
			}

			metaurls := []metalink.MetaURL{
				{
					URL:       fmt.Sprintf("%s//%s", i.config.RepositoryConfig.URI, releaseSubPath),
					MediaType: boshreleasesource.DefaultMediaType,
				},
			}

			if !f.URISatisfied(nil, metaurls) {
				continue
			}

			results = append(results, releaseversion.Artifact{
				Datastore: i.name,
				Name:      release.Name,
				Version:   release.Version,
				SourceTarball: metalink.File{
					Name:     fmt.Sprintf("%s-%s.tgz", release.Name, release.Version),
					MetaURLs: metaurls,
				},
				Labels: i.config.Labels,
			})
		}
	}

	if len(results) == 0 && i.config.DevReleases {
		// not really practical to support enumerating every commit of the repo as a version
		results = i.attemptDevReleaseFallback(f)
	}

	return inmemory.LimitArtifacts(results, l)
}

func (i *index) GetLabels() ([]string, error) {
	all, err := i.GetArtifacts(datastore.FilterParams{}, datastore.LimitParams{})
	if err != nil {
		return nil, errors.Wrap(err, "getting artifacts")
	}

	labelsMap := map[string]struct{}{}

	for _, one := range all {
		for _, label := range one.Labels {
			labelsMap[label] = struct{}{}
		}
	}

	var labels []string

	for label := range labelsMap {
		labels = append(labels, label)
	}

	return labels, nil
}

func (i *index) FlushCache() error {
	// TODO defer reload?
	return i.repository.ForceReload()
}

var devReleaseVersionRegex = regexp.MustCompile(`[^\-]+-.*commit\.([a-f0-9]{7,})`)

func (i *index) attemptDevReleaseFallback(f datastore.FilterParams) []releaseversion.Artifact {
	if !f.NameExpected || !f.VersionExpected {
		// name, version are expected, required
		return nil
	}

	match := devReleaseVersionRegex.FindStringSubmatch(f.Version)
	if match == nil {
		return nil
	}

	return []releaseversion.Artifact{
		{
			Datastore: i.name,
			Name:      f.Name,
			Version:   f.Version,
			Labels:    append(i.config.Labels, i.config.DevLabels...),
			SourceTarball: metalink.File{
				Name:    fmt.Sprintf("%s-%s.tgz", f.Name, f.Version),
				Version: f.Version,
				MetaURLs: []metalink.MetaURL{
					{
						URL: fmt.Sprintf(
							"%s?dev_release=true&name=%s&version=%s&checkout=%s",
							i.config.RepositoryConfig.URI,
							url.QueryEscape(f.Name),
							url.QueryEscape(f.Version),
							url.QueryEscape(match[1]),
						),
						MediaType: boshreleasesource.DefaultMediaType,
					},
				},
			},
		},
	}
}
