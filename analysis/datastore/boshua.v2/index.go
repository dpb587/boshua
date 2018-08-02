package boshuaV2

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	boshuaV2 "github.com/dpb587/boshua/artifact/datastore/datastoreutil/boshua.v2"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/compilation"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	releaseversiongraphql "github.com/dpb587/boshua/releaseversion/graphql"
	"github.com/dpb587/boshua/stemcellversion"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	stemcellversiongraphql "github.com/dpb587/boshua/stemcellversion/graphql"
	"github.com/dpb587/metalink"
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type index struct {
	logger logrus.FieldLogger
	config Config
	client *boshuaV2.Client
}

var _ datastore.Index = &index{}

func New(config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		logger: logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config: config,
		client: boshuaV2.NewClient(http.DefaultClient, config.BoshuaConfig, logger),
	}
}

func (i *index) GetAnalysisArtifacts(ref analysis.Reference) ([]analysis.Artifact, error) {
	fQuery, fQueryVarsTypes, fQueryVars, err := i.getQuery(ref)
	if err != nil {
		return nil, errors.Wrap(err, "building query")
	}

	cmd := fmt.Sprintf(`query (%s, $analyzer: String!) { %s }`, fQueryVarsTypes, fQuery)
	fQueryVars["analyzer"] = ref.Analyzer

	req := graphql.NewRequest(cmd)

	for k, v := range fQueryVars {
		req.Var(k, v)
	}

	var resp filterResponse

	err = i.client.Execute(req, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "executing remote request")
	}

	var results []analysis.Artifact

	for _, a := range resp.GetAnalysis() {
		results = append(results, analysis.New(ref, a.Artifact))
	}

	return results, nil
}

func (i *index) StoreAnalysisResult(_ analysis.Reference, _ metalink.Metalink) error {
	return datastore.UnsupportedOperationErr
}

func (i *index) FlushAnalysisCache() error {
	return nil // unsupported
}

func (i *index) getQuery(ref analysis.Reference) (string, string, map[string]interface{}, error) {
	subjectRef := ref.Subject.Reference()

	switch subjectRef := subjectRef.(type) {
	case releaseversion.Reference:
		fQuery, fQueryVarsTypes, fQueryVars := releaseversiongraphql.BuildListQueryArgs(releaseversiondatastore.FilterParamsFromReference(subjectRef))

		return fmt.Sprintf(`release(%s) {
			analysis {
				results(analyzers: [$analyzer]) {
					artifact {
						name
						size
						hashes {
							type
							hash
						}
						urls {
							url
						}
						metaurls {
							url
							mediatype
						}
					}
				}
			}
		}`, fQuery), fQueryVarsTypes, fQueryVars, nil
	case compilation.Reference:
		fQuery, fQueryVarsTypes, fQueryVars := releaseversiongraphql.BuildListQueryArgs(releaseversiondatastore.FilterParamsFromReference(subjectRef.ReleaseVersion))
		fQueryVarsTypes = fmt.Sprintf("%s, $queryStemcellOS: String!, $queryStemcellVersion: String!", fQueryVarsTypes)
		fQueryVars["queryStemcellOS"] = subjectRef.OSVersion.Name
		fQueryVars["queryStemcellVersion"] = subjectRef.OSVersion.Version
		return fmt.Sprintf(`release(%s) {
				compilation(os: $queryStemcellOS, version: $queryStemcellVersion) {
					analysis {
						results(analyzers: [$analyzer]) {
							artifact {
								name
								size
								hashes {
									type
									hash
								}
								urls {
									url
								}
								metaurls {
									url
									mediatype
								}
							}
						}
					}
				}
			}`, fQuery), fQueryVarsTypes, fQueryVars, nil
	case stemcellversion.Reference:
		fQuery, fQueryVarsTypes, fQueryVars := stemcellversiongraphql.BuildListQueryArgs(stemcellversiondatastore.FilterParamsFromReference(subjectRef))

		return fmt.Sprintf(`stemcell(%s) {
				analysis {
					results(analyzers: [$analyzer]) {
						artifact {
							name
							size
							hashes {
								type
								hash
							}
							urls {
								url
							}
							metaurls {
								url
								mediatype
							}
						}
					}
				}
			}`, fQuery), fQueryVarsTypes, fQueryVars, nil
	}

	return "", "", nil, errors.New("unsupported subject")
}
