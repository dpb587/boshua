package concourse

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/concourse/atc"
	"github.com/cppforlife/go-patch/patch"
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/analyzer/factory"
	"github.com/dpb587/boshua/releaseversion"
	compilationtask "github.com/dpb587/boshua/releaseversion/compilation/task"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/task"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/dpb587/boshua/task/scheduler/storecommon"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type ConfigLoader func() ([]byte, error)

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

func (s Scheduler) ScheduleAnalysis(analysisRef analysis.Reference) (scheduler.ScheduledTask, error) {
	tt, err := factory.SoonToBeDeprecatedFactory.BuildTask(analysisRef.Analyzer, analysisRef.Subject)
	if err != nil {
		return nil, errors.Wrap(err, "preparing task")
	}

	tt = storecommon.AppendAnalysisStore(tt, analysisRef)

	return s.schedule(tt)
}

func (s Scheduler) ScheduleCompilation(release releaseversion.Artifact, stemcell stemcellversion.Artifact) (scheduler.ScheduledTask, error) {
	tt, err := compilationtask.New(release, stemcell)
	if err != nil {
		return nil, errors.Wrap(err, "preparing task")
	}

	tt = storecommon.AppendCompilationStore(tt, release, stemcell)

	return s.schedule(tt)
}

func (s Scheduler) schedule(tt *task.Task) (scheduler.ScheduledTask, error) {
	fly := NewFly(s.config)

	pipelineBytes, pipelineVars, pipelineOps, err := s.buildBasePipeline(tt)
	if err != nil {
		return nil, errors.Wrap(err, "building pipeline")
	}

	pipelineName := s.pipelineName(tt, pipelineBytes)

	pipelineBytes, err = s.buildFinalPipeline(pipelineBytes, pipelineOps)
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

	return newScheduledTask(fly, pipelineName), nil
}

func (s *Scheduler) pipelineName(t *task.Task, pipelineBytes []byte) string {
	return fmt.Sprintf("boshua:%s:%x", t.Type, sha1.Sum(pipelineBytes))
}

func (s *Scheduler) buildFinalPipeline(pipelineBytes []byte, opDefs []patch.OpDefinition) ([]byte, error) {
	var pipeline interface{}

	err := yaml.Unmarshal(pipelineBytes, &pipeline)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshaling internal pipeline")
	}

	ops, err := patch.NewOpsFromDefinitions(opDefs)
	if err != nil {
		return nil, errors.Wrap(err, "building ops")
	}

	pipeline, err = ops.Apply(pipeline)
	if err != nil {
		return nil, errors.Wrap(err, "applying ops")
	}

	pipelineBytes, err = yaml.Marshal(pipeline)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling internal pipeline")
	}

	return pipelineBytes, nil
}

func (s *Scheduler) buildBasePipeline(t *task.Task) ([]byte, map[string]interface{}, []patch.OpDefinition, error) {
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

		runConfig := atc.TaskRunConfig{
			Path: "boshua",
			Args: step.Args,
		}

		if step.Privileged {
			runConfig.Path = "bash"
			runConfig.Args = append([]string{"-c", fmt.Sprintf(`%s
exec boshua "$@"`, privilegedMounts), "--"}, runConfig.Args...)
		}

		plan = append(plan, atc.PlanConfig{
			Task:       step.Name,
			Privileged: step.Privileged,
			TaskConfig: &atc.TaskConfig{
				Platform:      "linux",
				ImageResource: imageResource,
				Run:           runConfig,
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
				"PATH":          "config/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
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
					"interval": "672h", // some long period which will avoid rerunning
				},
			},
		},
	}

	pipelineBytes, err := yaml.Marshal(pipeline)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "marshaling pipeline")
	}

	rawConfigBytes, err := s.boshuaConfigLoader()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "loading raw config")
	}

	pipelineVars := map[string]interface{}{
		"boshua_config": string(rawConfigBytes),
	}

	var pipelineOps []patch.OpDefinition

	for _, taskConfig := range s.config.Tasks {
		if taskConfig.Type != string(t.Type) {
			continue
		}

		for _, opsFile := range taskConfig.OpsFiles {
			opsBytes, err := ioutil.ReadFile(opsFile)
			if err != nil {
				return nil, nil, nil, errors.Wrap(err, "reading ops file")
			}

			var opDefs []patch.OpDefinition

			err = yaml.Unmarshal(opsBytes, &opDefs)
			if err != nil {
				return nil, nil, nil, errors.Wrap(err, "unmarshaling ops file")
			}

			pipelineOps = append(pipelineOps, opDefs...)
		}

		pipelineOps = append(pipelineOps, taskConfig.Ops...)

		for varKey, varVal := range taskConfig.Vars {
			pipelineVars[varKey] = varVal
		}
	}

	return pipelineBytes, pipelineVars, pipelineOps, nil
}

var privilegedMounts = `set -eu

# This is copied from https://github.com/concourse/concourse/blob/3c070db8231294e4fd51b5e5c95700c7c8519a27/jobs/baggageclaim/templates/baggageclaim_ctl.erb#L23-L54
# helps the /dev/mapper/control issue and lets us actually do scary things with the /dev mounts
# This allows us to create device maps from partition tables in image_create/apply.sh
function permit_device_control() {
	local devices_mount_info=$(cat /proc/self/cgroup | grep devices)

	local devices_subsytems=$(echo $devices_mount_info | cut -d: -f2)
	local devices_subdir=$(echo $devices_mount_info | cut -d: -f3)

	cgroup_dir=/mnt/tmp-todo-devices-cgroup

	if [ ! -e ${cgroup_dir} ]; then
		# mount our container's devices subsystem somewhere
		mkdir ${cgroup_dir}
	fi

	if ! mountpoint -q ${cgroup_dir}; then
		mount -t cgroup -o $devices_subsytems none ${cgroup_dir}
	fi

	# permit our cgroup to do everything with all devices
	# ignore failure in case something has already done this; echo appears to
	# return EINVAL, possibly because devices this affects are already in use
	echo a > ${cgroup_dir}${devices_subdir}/devices.allow || true
}

permit_device_control

# Also copied from baggageclaim_ctl.erb creates 64 loopback mappings. This fixes failures with losetup --show --find ${disk_image}
for i in $(seq 0 64); do
	if ! mknod -m 0660 /dev/loop$i b 7 $i; then
		break
	fi
done
`
