package metalinkutil

import (
	"fmt"
	"path"

	"github.com/dpb587/metalink"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
	"github.com/dpb587/metalink/verification"
	"github.com/dpb587/metalink/verification/hash"
	"github.com/pkg/errors"
)

func CreateFromFiles(paths ...string) (*metalink.Metalink, error) {
	meta4 := metalink.Metalink{}

	urlLoader := urldefaultloader.New()

	for _, meta4FilePath := range paths {
		var err error

		meta4file := metalink.File{
			Name: path.Base(meta4FilePath),
			URLs: []metalink.URL{
				{
					URL: meta4FilePath,
				},
			},
		}

		origin, err := urlLoader.Load(meta4file.URLs[0])
		if err != nil {
			return nil, errors.Wrap(err, "parsing origin destination")
		}

		meta4file.Size, err = origin.Size()
		if err != nil {
			return nil, errors.Wrap(err, "Loading size")
		}

		hashmap := map[string]verification.Signer{
			"sha-512": hash.SHA512SignerVerifier,
			"sha-256": hash.SHA256SignerVerifier,
			"sha-1":   hash.SHA1SignerVerifier,
			"md5":     hash.MD5SignerVerifier,
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

			err = verification.Apply(&meta4file)
			if err != nil {
				return nil, errors.Wrap(err, "Adding verification to file")
			}
		}

		meta4.Files = append(meta4.Files, meta4file)
	}

	return &meta4, nil
}
