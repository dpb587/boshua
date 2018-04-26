package compilation

import (
	"encoding/json"
	"fmt"

	"github.com/concourse/atc"
	"github.com/dpb587/boshua"
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/scheduler/task"
)

type Task struct {
	artifact       compiledreleaseversion.Artifact
	releaseVersion releaseversion.Artifact
	osVersion      osversion.Artifact
}

var _ task.Task = &Task{}

func (t Task) Type() string {
	return "compilation"
}

func (t Task) ArtifactReference() boshua.Reference {
	return t.artifact.ArtifactReference()
}

func (t Task) Config() (atc.Config, error) {
	contextBytes, err := json.MarshalIndent(map[string]interface{}{
		"release": map[string]interface{}{
			"name":      t.artifact.ReleaseVersion.Name,
			"version":   t.artifact.ReleaseVersion.Version,
			"checksums": t.artifact.ReleaseVersion.Checksums,
		},
		"os": map[string]interface{}{
			"name":    t.artifact.OSVersion.Name,
			"version": t.artifact.OSVersion.Version,
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
				Source:     atc.Source(t.osVersion.ArtifactMetalinkStorage()),
			},
			{
				Name:       "release",
				CheckEvery: "24h",
				Type:       "metalink-repository",
				Source:     atc.Source(t.releaseVersion.ArtifactMetalinkStorage()),
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
									"storage":         t.artifact.StoragePath(),
									"context":         string(contextBytes),
									"release_version": t.artifact.ReleaseVersion.Version,
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
