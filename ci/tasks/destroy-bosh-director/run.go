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

type Lock struct {
	Sponsor map[string]string      `json:"sponsor"`
	IaaS    string                 `json:"iaas"`
	Vars    map[string]interface{} `json:"vars"`
	Ops     []string               `json:"ops"`
}

func main() {
	lockBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Panicf("reading lock: %v", err)
	}

	var lock Lock

	err = json.Unmarshal(lockBytes, &lock)
	if err != nil {
		log.Panicf("parsing lock: %v", err)
	}

	var now = time.Now()

	{ // delete-env
		cmdArgs := []string{
			"delete-env",
			"bosh-deployment/bosh.yml",
			"--state", "bosh-director/state.json",
			"--vars-store", "bosh-director/vars.yml",
			"--ops-file", fmt.Sprintf("bosh-deployment/%s/cpi.yml", lock.IaaS),
			"--ops-file", "bosh-deployment/external-ip-not-recommended.yml",
			"--ops-file", "bosh-deployment/bosh-lite.yml",
			"--var", fmt.Sprintf("director_name=bosh-lite-%s", now.Format("20060102T150405")),
		}

		for varKey, varVal := range lock.Vars {
			cmdArgs = append(cmdArgs, "--var", fmt.Sprintf("%s=%s", varKey, varVal))
		}

		for _, opPath := range lock.Ops {
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
