package concourse

import (
	"bytes"
	"io/ioutil"

	"github.com/dpb587/boshua/task"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Scheduler struct {
	config Config
	logger logrus.FieldLogger
}

var _ scheduler.Scheduler = &Scheduler{}

func New(config Config, logger logrus.FieldLogger) *Scheduler {
	return &Scheduler{
		config: config,
		logger: logger,
	}
}

func (s *Scheduler) Schedule(t task.Task) (scheduler.Task, error) {
	fly := NewFly(s.config)

	pipelineName := s.pipelineName(t)

	file, err := ioutil.TempFile("", "boshua-")
	if err != nil {
		return nil, errors.Wrap(err, "creating temp file")
	}

	defer file.Close()

	pipelineBytes, err := s.buildPipeline(t)
	if err != nil {
		return nil, errors.Wrap(err, "building pipeline")
	}

	_, err = file.Write(pipelineBytes)
	if err != nil {
		return nil, errors.Wrap(err, "writing pipeline")
	}

	_, _, err = fly.RunWithStdin(
		bytes.NewBufferString("y\n"),
		"set-pipeline",
		"--pipeline", pipelineName,
		"--config", file.Name(),
		"--load-vars-from", s.config.SecretsPath,
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

func (s *Scheduler) pipelineName(t task.Task) string {
	// TODO
	return "asdf"
}

func (s *Scheduler) buildPipeline(t task.Task) ([]byte, error) {
	return nil, nil
}
