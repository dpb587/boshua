package task

import (
	"fmt"

	"github.com/concourse/atc"
	"github.com/dpb587/boshua"
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/scheduler/task"
)

type Task struct {
	subject  analysis.Subject
	analyzer string
}

var _ task.Task = &Task{}

func (t Task) Type() string {
	return t.analyzer
}

func (t Task) ArtifactReference() boshua.Reference {
	return t.subject.ArtifactReference()
}

func (t Task) Config() (atc.Config, error) {
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
				Name:       "bosh-compiled-releases",
				CheckEvery: "24h",
				Type:       "git",
				Source: atc.Source{
					"uri":         "git@github.com:dpb587/bosh-compiled-releases-v2.git",
					"private_key": "((bcr_private_key))",
				},
			},
			{
				Name:       "artifact",
				CheckEvery: "24h",
				Type:       "metalink-repository",
				Source:     atc.Source(t.subject.ArtifactMetalinkStorage()),
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
				Name:                 "analyze",
				DisableManualTrigger: true,
				Serial:               true,
				RawMaxInFlight:       1,
				Plan: atc.PlanSequence{
					{
						Aggregate: &atc.PlanSequence{
							{
								Get:     "artifact",
								Trigger: true,
							},
							{
								Get: "bosh-compiled-releases",
							},
							{
								Get: "index",
							},
						},
					},
					{
						Task:           "analyze",
						TaskConfigPath: "bosh-compiled-releases/ci/tasks/generate-analysis/task.yml",
						Params: atc.Params{
							"analyzer": t.analyzer,
						},
					},
					{
						Task:           "publish",
						TaskConfigPath: "bosh-compiled-releases/ci/tasks/publish-analysis/task.yml",
						Params: atc.Params{
							"storage": fmt.Sprintf(
								"%s/%s/%s",
								t.subject.ArtifactReference().Context,
								t.subject.ArtifactReference().ID,
								t.analyzer,
							),
							"analyzer":      t.analyzer,
							"s3_bucket":     "((s3_bucket))",
							"s3_host":       "((s3_host))",
							"s3_access_key": "((s3_access_key))",
							"s3_secret_key": "((s3_secret_key))",
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
			},
		},
	}, nil
}
