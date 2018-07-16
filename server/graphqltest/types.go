package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Masterminds/semver"
	"github.com/dpb587/metalink"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

type release struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Checksums []string
	URIs      []string
	Tarball   *metalink.File `json:"tarball"`
}

var data []release = []release{
	{
		Name:    "openvpn",
		Version: "4.2.1",
		Checksums: []string{
			"sha1:def",
			"sha256:ghijkl",
		},
		URIs: []string{
			"https://github.com/dpb587/openvpn-bosh-release.git",
		},
		Tarball: &metalink.File{
			Name: "openvpn-4.2.1.tar.gz",
			Size: 12345,
			Hashes: []metalink.Hash{
				{
					Type: "sha-1",
					Hash: "asdfasdfsaf",
				},
			},
			URLs: []metalink.URL{
				{
					URL: "asdfasdf",
				},
			},
		},
	},
	{
		Name:    "openvpn",
		Version: "5.0.0",
		Checksums: []string{
			"sha1:abc",
			"sha256:defghi",
		},
		URIs: []string{
			"https://github.com/dpb587/openvpn-bosh-release.git",
		},
		Tarball: &metalink.File{
			Name: "openvpn-5.0.0.tar.gz",
			Size: 612345,
			URLs: []metalink.URL{
				{
					URL: "6asdfasdf",
				},
			},
		},
	},
}

var assetType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Asset",
		Description: "A specific downloadable asset.",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"hashes": &graphql.Field{
				Type: graphql.NewList(assetHashType),
			},
			"size": &graphql.Field{
				Type: graphql.Int,
			},
			"urls": &graphql.Field{
				Type: graphql.NewList(assetURLType),
			},
		},
	},
)

var assetHashType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Hash",
		Description: "A verifiable checksum.",
		Fields: graphql.Fields{
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"hash": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var assetURLType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "URL",
		Description: "A URL for downloading the asset.",
		Fields: graphql.Fields{
			"url": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var assetMetaURLType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "MetaURL",
		Description: "A Meta URL for downloading the asset.",
		Fields: graphql.Fields{
			"url": &graphql.Field{
				Type: graphql.String,
			},
			"mimetype": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var releaseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Release",
		Description: "A specific version of a release.",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"version": &graphql.Field{
				Type: graphql.String,
			},
			"tarball": &graphql.Field{
				Type: assetType,
			},
		},
	},
)

var rootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"releases": &graphql.Field{
				Type: graphql.NewList(releaseType),
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"version": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"checksum": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"uri": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var resolved = []release{}
					var expectedParsedVersion *semver.Constraints
					var err error

					expectedName, expectingName := p.Args["name"].(string)
					expectedVersion, expectingVersion := p.Args["version"].(string)
					expectedChecksum, expectingChecksum := p.Args["checksum"].(string)
					expectedURI, expectingURI := p.Args["uri"].(string)

					if expectingVersion {
						expectedParsedVersion, err = semver.NewConstraint(expectedVersion)
						if err != nil {
							return nil, errors.Wrap(err, "parsing version")
						}
					}

					for _, subject := range data {
						if expectingName && expectedName != subject.Name {
							continue
						} else if expectingVersion {
							if expectedVersion == subject.Version {
								// great
							} else {
								actualParsedVersion, err := semver.NewVersion(subject.Version)
								if err != nil {
									continue
								}

								if !expectedParsedVersion.Check(actualParsedVersion) {
									continue
								}
							}
						}

						if expectingChecksum {
							var found bool

							for _, actual := range subject.Checksums {
								if expectedChecksum == actual {
									found = true

									break
								}
							}

							if !found {
								continue
							}
						}

						if expectingURI {
							var found bool

							for _, actual := range subject.URIs {
								if expectedURI == actual {
									found = true

									break
								}
							}

							if !found {
								continue
							}
						}

						resolved = append(resolved, subject)
					}

					return resolved, nil
				},
			},
		},
	})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: rootQuery,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.Encode(result)
	})

	http.HandleFunc(pattern, func(arg1 http.ResponseWriter, arg2 *http.Request) {

	})
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	// Display some basic instructions
	fmt.Println("Now server is running on port 8080")
	fmt.Println("Get single todo: curl -g 'http://localhost:8080/graphql?query={release_list(name:\"openvpn\"){version}}'")
	fmt.Println("Access the web app via browser at 'http://localhost:8080'")

	http.ListenAndServe(":8080", nil)
}
