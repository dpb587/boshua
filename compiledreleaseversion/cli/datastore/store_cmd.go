package datastore

import (
	"fmt"
	"path/filepath"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/metalink"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
	"github.com/dpb587/metalink/verification"
	"github.com/dpb587/metalink/verification/hash"
	"github.com/pkg/errors"
)

type StoreCmd struct {
	*CmdOpts `no-flag:"true"`

	Version string `long:"version" description:"A specific version to use" default:"0.0.0"`

	Args StoreCmdArgs `positional-args:"true" required:"true"`
}

type StoreCmdArgs struct {
	Local string `positional-arg-name:"PATH" description:"Path to the artifact"`
}

func (c *StoreCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("compiledrelease/datastore/filter")

	index, err := c.getDatastore()
	if err != nil {
		return fmt.Errorf("loading datastore: %v", err)
	}

	rawCompiledReleaseRef := c.CompiledReleaseOpts.Reference()

	releaseVersionIndex, err := c.AppOpts.GetReleaseIndex("default")
	if err != nil {
		return fmt.Errorf("loading release index: %v", err)
	}

	releaseVersion, err := releaseVersionIndex.Find(rawCompiledReleaseRef.ReleaseVersion)
	if err != nil {
		return fmt.Errorf("finding release: %v", err)
	}

	osVersionIndex, err := c.AppOpts.GetOSIndex("default")
	if err != nil {
		return fmt.Errorf("loading os index: %v", err)
	}

	osVersion, err := osVersionIndex.Find(rawCompiledReleaseRef.OSVersion)
	if err != nil {
		return fmt.Errorf("finding os: %v", err)
	}

	meta4File, err := c.createMetalinkFile()
	if err != nil {
		return fmt.Errorf("building metalink: %v", err)
	}

	return index.Store(compiledreleaseversion.New(
		releaseVersion.Reference,
		osVersion.Reference,
		*meta4File,
		map[string]interface{}{},
	))
}

func (c *StoreCmd) createMetalinkFile() (*metalink.File, error) {
	logger := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(logger)

	urlLoader := urldefaultloader.New(fs)

	file := metalink.File{
		Name:    filepath.Base(c.Args.Local),
		Version: c.Version,
		Hashes:  []metalink.Hash{},
	}

	origin, err := urlLoader.Load(metalink.URL{URL: c.Args.Local})
	if err != nil {
		return nil, errors.Wrap(err, "Loading origin")
	}

	file.Size, err = origin.Size()
	if err != nil {
		return nil, errors.Wrap(err, "Loading size")
	}

	hashmap := map[string]verification.Signer{
		"sha-512": hash.SHA512Verification,
		"sha-256": hash.SHA256Verification,
		"sha-1":   hash.SHA1Verification,
		"md5":     hash.MD5Verification,
	}

	for _, signer := range hashmap {
		verification, err := signer.Sign(origin)
		if err != nil {
			return nil, errors.Wrap(err, "Signing hash")
		}

		err = verification.Apply(&file)
		if err != nil {
			return nil, errors.Wrap(err, "Adding verification to file")
		}
	}

	return &file, nil
}
