package metalinkutil

import (
	"fmt"

	"github.com/dpb587/metalink"
)

func ToMetalinkHashType(algorithm string) (metalink.HashType, error) {
	switch algorithm {
	case "md5":
		return metalink.HashTypeMD5, nil
	case "sha1":
		return metalink.HashTypeSHA1, nil
	case "sha256":
		return metalink.HashTypeSHA256, nil
	case "sha512":
		return metalink.HashTypeSHA512, nil
	}

	return "", fmt.Errorf("unrecognized hash type: %s", algorithm)
}

func FromMetalinkHashType(algorithm metalink.HashType) (string, error) {
	switch algorithm {
	case metalink.HashTypeMD5:
		return "md5", nil
	case metalink.HashTypeSHA1:
		return "sha1", nil
	case metalink.HashTypeSHA256:
		return "sha256", nil
	case metalink.HashTypeSHA512:
		return "sha512", nil
	}

	return "", fmt.Errorf("unrecognized metalink hash type: %s", algorithm)
}
