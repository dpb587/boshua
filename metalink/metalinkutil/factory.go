package metalinkutil

import (
	"fmt"
	"path"

	"github.com/dpb587/metalink"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
	"github.com/dpb587/metalink/verification"
	"github.com/dpb587/metalink/verification/hash"
	"github.com/pkg/errors"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

func CreateFromFiles(paths ...string) (*metalink.Metalink, error) {
	meta4 := metalink.Metalink{}

	logger := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(logger)

	urlLoader := urldefaultloader.New(fs)

	for _, meta4FilePath := range paths {
		file := metalink.File{
			Name: path.Base(meta4FilePath),
			URLs: []metalink.URL{
				{
					URL: meta4FilePath,
				},
			},
		}

		origin, err := urlLoader.Load(file.URLs[0])
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

		for _, hashType := range []string{"sha-512", "sha-256", "sha-1", "md5"} {
			signer, found := hashmap[hashType]
			if !found {
				return nil, fmt.Errorf("unknown hash type: %s", hashType)
			}

			verification, err := signer.Sign(origin)
			if err != nil {
				return nil, errors.Wrap(err, "Signing hash")
			}

			err = verification.Apply(&file)
			if err != nil {
				return nil, errors.Wrap(err, "Adding verification to file")
			}
		}

		meta4.Files = append(meta4.Files, file)
	}

	return &meta4, nil
}
