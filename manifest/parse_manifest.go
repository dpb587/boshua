package manifest

type parseManifest struct {
	Releases []parseManifestRelease `yaml:"releases"`

	InstanceGroups []parseManifestInstanceGroup `yaml:"instance_groups"`
	Jobs           []parseManifestInstanceGroup `yaml:"jobs"`

	Stemcell *parseManifestStemcell `yaml:"stemcell"`

	// init
	ResourcePools []parseManifestResourcePool `yaml:"resource_pools"`
	CloudProvider parseManifestCloudProvider  `yaml:"cloud_provider"`
}

func (m parseManifest) InstalledReleases() []parseManifestReleaseRef {
	var releases []parseManifestReleaseRef

	for _, ig := range m.InstanceGroups {
		releases = append(releases, ig.Jobs...)
		releases = append(releases, ig.Templates...)
	}

	for _, j := range m.Jobs {
		releases = append(releases, j.Jobs...)
		releases = append(releases, j.Templates...)
	}

	return releases
}

type parseManifestRelease struct {
	Name     string                        `yaml:"name"`
	Version  string                        `yaml:"version"`
	Sha1     string                        `yaml:"sha1"`
	URL      string                        `yaml:"url"`
	Stemcell *parseManifestReleaseStemcell `yaml:"stemcell"`
}

type parseManifestReleaseStemcell struct{}

type parseManifestInstanceGroup struct {
	Jobs      []parseManifestReleaseRef `yaml:"jobs"`
	Templates []parseManifestReleaseRef `yaml:"templates"`
}

type parseManifestResourcePool struct {
	Stemcell *parseManifestStemcell `yaml:"stemcell"`
}

type parseManifestStemcell struct {
	Name    string `yaml:"name"`
	OS      string `yaml:"os"`
	Version string `yaml:"version"`
}

type parseManifestCloudProvider struct {
	Template *parseManifestReleaseRef `yaml:"template"`
}

type parseManifestReleaseRef struct {
	Release string `yaml:"release"`
}
