package util

import "fmt"

func MetalinkHashType(algorithm string) (string, error) {
	switch algorithm {
	case "md5":
		return "md5", nil
	case "sha1":
		return "sha-1", nil
	case "sha256":
		return "sha-256", nil
	case "sha512":
		return "sha-512", nil
	}

	return "", fmt.Errorf("unrecognized hash type: %s", algorithm)
}

func FromMetalinkHashType(algorithm string) (string, error) {
	switch algorithm {
	case "md5":
		return "md5", nil
	case "sha-1":
		return "sha1", nil
	case "sha-256":
		return "sha256", nil
	case "sha-512":
		return "sha512", nil
	}

	return "", fmt.Errorf("unrecognized metalink hash type: %s", algorithm)
}
