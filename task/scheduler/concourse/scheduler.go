package concourse

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/concourse/atc"
	"github.com/cppforlife/go-patch/patch"
	"github.com/dpb587/boshua/config"
	"github.com/dpb587/boshua/task"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type ConfigLoader func() (*config.Config, error)

type Scheduler struct {
	config             Config
	boshuaConfigLoader ConfigLoader
	logger             logrus.FieldLogger
}

var _ scheduler.Scheduler = &Scheduler{}

func New(config Config, boshuaConfigLoader ConfigLoader, logger logrus.FieldLogger) *Scheduler {
	return &Scheduler{
		config:             config,
		boshuaConfigLoader: boshuaConfigLoader,
		logger:             logger,
	}
}

func (s *Scheduler) Schedule(t task.Task) (scheduler.Task, error) {
	fly := NewFly(s.config)

	pipelineBytes, pipelineVars, pipelineOpsFiles, err := s.buildBasePipeline(t)
	if err != nil {
		return nil, errors.Wrap(err, "building pipeline")
	}

	pipelineName := s.pipelineName(t, pipelineBytes)

	pipelineBytes, err = s.buildFinalPipeline(pipelineBytes, pipelineOpsFiles)
	if err != nil {
		return nil, errors.Wrap(err, "apply pipeline ops")
	}

	pipelineFile, err := ioutil.TempFile("", "boshua-pipeline-")
	if err != nil {
		return nil, errors.Wrap(err, "creating temp pipeline")
	}

	defer os.RemoveAll(pipelineFile.Name())

	_, err = pipelineFile.Write(pipelineBytes)
	if err != nil {
		return nil, errors.Wrap(err, "writing pipeline")
	}

	err = pipelineFile.Close()
	if err != nil {
		return nil, errors.Wrap(err, "closing temp pipeline")
	}

	pipelineVarsFile, err := ioutil.TempFile("", "boshua-vars-")
	if err != nil {
		return nil, errors.Wrap(err, "creating temp pipeline vars")
	}

	defer os.RemoveAll(pipelineVarsFile.Name())

	pipelineVarsBytes, err := yaml.Marshal(pipelineVars)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling pipeline vars")
	}

	_, err = pipelineVarsFile.Write(pipelineVarsBytes)
	if err != nil {
		return nil, errors.Wrap(err, "writing pipeline vars")
	}

	err = pipelineVarsFile.Close()
	if err != nil {
		return nil, errors.Wrap(err, "closing temp pipeline vars")
	}

	_, _, err = fly.RunWithStdin(
		bytes.NewBufferString("y\n"),
		"set-pipeline",
		"--pipeline", pipelineName,
		"--config", pipelineFile.Name(),
		"--load-vars-from", pipelineVarsFile.Name(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "setting pipeline")
	}

	_, _, err = fly.Run(
		"unpause-pipeline",
		"--pipeline", pipelineName,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unpausing pipeline")
	}

	return NewTask(fly, pipelineName), nil
}

func (s *Scheduler) pipelineName(t task.Task, pipelineBytes []byte) string {
	return fmt.Sprintf("boshua:%s:%x", t.Type, sha1.Sum(pipelineBytes))
}

func (s *Scheduler) buildFinalPipeline(pipelineBytes []byte, opsFiles []string) ([]byte, error) {
	var pipeline interface{}

	err := yaml.Unmarshal(pipelineBytes, &pipeline)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshaling internal pipeline")
	}

	for _, opsFile := range opsFiles {
		opsBytes, err := ioutil.ReadFile(opsFile)
		if err != nil {
			return nil, errors.Wrap(err, "reading ops file")
		}

		var ops patch.Ops

		err = yaml.Unmarshal(opsBytes, &ops)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshaling ops file")
		}

		pipeline, err = ops.Apply(pipeline)
		if err != nil {
			return nil, errors.Wrap(err, "applying ops file") // TODO include file paths?
		}
	}

	pipelineBytes, err = yaml.Marshal(pipeline)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling internal pipeline")
	}

	return pipelineBytes, nil
}

func (s *Scheduler) buildBasePipeline(t task.Task) ([]byte, map[string]interface{}, []string, error) {
	var plan atc.PlanSequence

	configPath := "config/boshua.yml"
	imageResource := &atc.ImageResource{
		Type: "docker-image",
		Source: atc.Source{
			"repository": "dpb587/boshua",
		},
	}

	plan = append(plan, atc.PlanConfig{
		RawName: "trigger",
		Get:     "trigger",
		Trigger: true,
	})

	plan = append(plan, atc.PlanConfig{
		Task: "configure",
		TaskConfig: &atc.TaskConfig{
			Platform:      "linux",
			ImageResource: imageResource,
			Run: atc.TaskRunConfig{
				Path: "bash",
				Args: []string{
					"-c",
					fmt.Sprintf(`echo "$config" > %s`, configPath),
				},
			},
			Outputs: []atc.TaskOutputConfig{
				{
					Name: "config",
				},
			},
		},
		Params: atc.Params{
			"config": "((boshua_config))",
		},
	})

	for _, step := range t.Steps {
		if len(step.Input) > 0 {
			if len(step.Input) > 1 {
				panic("TODO support prior inputs, multiple files")
			}

			for fileName, fileData := range step.Input {
				plan = append(plan, atc.PlanConfig{
					Task: fmt.Sprintf("%s-input", step.Name),
					TaskConfig: &atc.TaskConfig{
						Platform:      "linux",
						ImageResource: imageResource,
						Run: atc.TaskRunConfig{
							Path: "bash",
							Args: []string{
								"-c",
								fmt.Sprintf(`echo "$input" > output/%s`, fileName),
							},
						},
						Outputs: []atc.TaskOutputConfig{
							{
								Name: "input",
								Path: "output",
							},
						},
					},
					Params: atc.Params{
						"input": string(fileData),
					},
				})
			}
		}

		plan = append(plan, atc.PlanConfig{
			Task:       step.Name,
			Privileged: step.Privileged,
			TaskConfig: &atc.TaskConfig{
				Platform:      "linux",
				ImageResource: imageResource,
				Run: atc.TaskRunConfig{
					Path: "boshua",
					Args: step.Args,
				},
				Inputs: []atc.TaskInputConfig{
					{
						Name: "config",
					},
					{
						Name: "input",
					},
				},
				Outputs: []atc.TaskOutputConfig{
					{
						Name: "input",
						Path: "output",
					},
				},
			},
			Params: atc.Params{
				"PATH":          "$PWD/config/bin:$PATH",
				"BOSHUA_CONFIG": configPath,
			},
		})
	}

	var pipeline = atc.Config{
		Jobs: []atc.JobConfig{
			{
				Name: "task",
				Plan: plan,
			},
		},
		Resources: []atc.ResourceConfig{
			{
				Name: "trigger",
				Type: "time",
				Source: atc.Source{
					"interval": "6h",
				},
			},
		},
	}

	pipelineBytes, err := yaml.Marshal(pipeline)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "marshaling pipeline")
	}

	rawConfig, err := s.boshuaConfigLoader()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "loading boshua config")
	}

	rawConfigBytes, err := yaml.Marshal(rawConfig)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "marshaling boshua config")
	}

	pipelineVars := map[string]interface{}{
		"boshua_config": string(rawConfigBytes),
	}

	pipelineOpsFiles := []string{
		// TODO lookup from config[type]
	}

	return pipelineBytes, pipelineVars, pipelineOpsFiles, nil
}
