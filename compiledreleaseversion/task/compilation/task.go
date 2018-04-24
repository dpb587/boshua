package compilation

import (
	"encoding/json"
	"fmt"

	"github.com/concourse/atc"
	"github.com/dpb587/boshua"
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/scheduler/task"
)

type Task struct {
	subject compiledreleaseversion.ResolvedSubject
}

var _ task.Task = &Task{}

func (t Task) Type() string {
	return "compilation"
}

func (t Task) SubjectReference() boshua.Reference {
	return t.subject.SubjectReference()
}

func (t Task) Config() (atc.Config, error) {
	contextBytes, err := json.MarshalIndent(map[string]interface{}{
		"release": map[string]interface{}{
			"name":      t.subject.ResolvedReleaseVersion.Name,
			"version":   t.subject.ResolvedReleaseVersion.Version,
			"checksums": t.subject.ResolvedReleaseVersion.Checksums,
		},
		"stemcell": map[string]interface{}{
			"os":      t.subject.ResolvedStemcellVersion.OS,
			"version": t.subject.ResolvedStemcellVersion.Version,
		},
	}, "", "  ")
	if err != nil {
		return atc.Config{}, fmt.Errorf("marshalling context: %v", err)
	}

	return atc.Config{
		ResourceTypes: atc.ResourceTypes{
			{
				Name: "metalink-repository",
				Type: "docker-image",
				Source: atc.Source{
					"repository": "dpb587/metalink-repository-resource",
				},
			},
		},
		Resources: atc.ResourceConfigs{
			{
				Name:       "bosh-deployment",
				CheckEvery: "24h",
				Type:       "git",
				Source: atc.Source{
					"uri": "https://github.com/cloudfoundry/bosh-deployment.git",
				},
			},
			{
				Name:       "bosh-compiled-releases",
				CheckEvery: "24h",
				Type:       "git",
				Source: atc.Source{
					"uri":         "git@github.com:dpb587/bosh-compiled-releases-v2.git",
					"private_key": "((bcr_private_key))",
				},
			},
			{
				Name:       "stemcell",
				CheckEvery: "24h",
				Type:       "metalink-repository",
				Source:     atc.Source(t.subject.ResolvedStemcellVersion.MetalinkSource),
			},
			{
				Name:       "release",
				CheckEvery: "24h",
				Type:       "metalink-repository",
				Source:     atc.Source(t.subject.ResolvedReleaseVersion.MetalinkSource),
			},
			{
				Name:       "env",
				CheckEvery: "24h",
				Type:       "pool",
				Source: atc.Source{
					"uri":         "git@github.com:dpb587/bosh-compiled-releases-envs.git",
					"branch":      "master",
					"private_key": "((envs_private_key))",
					"pool":        "dpb587-gcp",
				},
			},
			{
				Name:       "index",
				CheckEvery: "24h",
				Type:       "git",
				Source: atc.Source{
					"branch":      "master",
					"uri":         "git@github.com:dpb587/bosh-compiled-releases-index.git",
					"private_key": "((index_private_key))",
				},
			},
		},
		Jobs: atc.JobConfigs{
			{
				Name:                 "compile",
				DisableManualTrigger: true,
				Serial:               true,
				RawMaxInFlight:       1,
				Plan: atc.PlanSequence{
					{
						Aggregate: &atc.PlanSequence{
							{
								Get:     "stemcell",
								Trigger: true,
							},
							{
								Get:     "release",
								Trigger: true,
							},
							{
								Get: "bosh-compiled-releases",
							},
							{
								Get: "index",
							},
							{
								Get: "bosh-deployment",
							},
						},
					},
					{
						Put: "env",
						Params: atc.Params{
							"acquire": true,
						},
					},
					{
						Do: &atc.PlanSequence{
							{
								Task:           "create-bosh-director",
								TaskConfigPath: "bosh-compiled-releases/ci/tasks/create-bosh-director/task.yml",
							},
							{
								Task:           "upload-stemcell",
								TaskConfigPath: "bosh-compiled-releases/ci/tasks/upload-stemcell/task.yml",
							},
							{
								Task:           "upload-release",
								TaskConfigPath: "bosh-compiled-releases/ci/tasks/upload-release/task.yml",
							},
							{
								Task:           "compile-release",
								TaskConfigPath: "bosh-compiled-releases/ci/tasks/compile-release/task.yml",
							},
							{
								Task:           "publish-compiled-release",
								TaskConfigPath: "bosh-compiled-releases/ci/tasks/publish-compiled-release/task.yml",
								Params: atc.Params{
									"storage":         t.subject.StoragePath(),
									"context":         string(contextBytes),
									"release_version": t.subject.ResolvedReleaseVersion.Version,
									"s3_bucket":       "((s3_bucket))",
									"s3_host":         "((s3_host))",
									"s3_access_key":   "((s3_access_key))",
									"s3_secret_key":   "((s3_secret_key))",
								},
							},
							{
								Put: "index",
								Params: atc.Params{
									"repository": "index",
									"rebase":     true,
								},
							},
						},
						Ensure: &atc.PlanConfig{
							Do: &atc.PlanSequence{
								{
									Task:           "destroy-bosh-director",
									TaskConfigPath: "bosh-compiled-releases/ci/tasks/destroy-bosh-director/task.yml",
								},
								{
									Put: "env",
									Params: atc.Params{
										"release": "env",
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil
}
