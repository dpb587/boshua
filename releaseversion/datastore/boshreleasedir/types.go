package boshreleasedir

type boshReleaseIndex struct {
	Builds map[string]boshReleaseIndexBuild `yaml:"builds"`
}

type boshReleaseIndexBuild struct {
	Version string `yaml:"version"`
}

type boshRelease struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}
