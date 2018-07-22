package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

type Env struct {
	IaaS string                 `json:"iaas"`
	Vars map[string]interface{} `json:"vars"`
	Ops  []string               `json:"ops"`
}

func main() {
	envBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Panicf("reading env: %v", err)
	}

	var env Env

	err = json.Unmarshal(envBytes, &env)
	if err != nil {
		log.Panicf("parsing env: %v", err)
	}

	var now = time.Now()

	{ // delete-env
		cmdArgs := []string{
			"delete-env",
			"bosh-deployment/bosh.yml",
			"--state", "bosh-director/state.json",
			"--vars-store", "bosh-director/vars.yml",
			"--ops-file", fmt.Sprintf("bosh-deployment/%s/cpi.yml", env.IaaS),
			"--ops-file", "bosh-deployment/external-ip-not-recommended.yml",
			"--ops-file", "bosh-deployment/bosh-lite.yml",
			"--var", fmt.Sprintf("director_name=bosh-lite-%s", now.Format("20060102T150405")),
		}

		for varKey, varVal := range env.Vars {
			cmdArgs = append(cmdArgs, "--var", fmt.Sprintf("%s=%s", varKey, varVal))
		}

		for _, opPath := range env.Ops {
			cmdArgs = append(cmdArgs, "--ops-file", opPath)
		}

		cmd := exec.Command("bosh", cmdArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Start()
		if err != nil {
			log.Panicf("executing create-env: %v", err)
		}

		err = cmd.Wait()
		if err != nil {
			log.Panicf("waiting for create-env: %v", err)
		}
	}
}
