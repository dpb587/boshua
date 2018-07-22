package main

import (
	"bytes"
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
	var varsPath = "bosh-director/vars.yml"

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

	{ // create-env
		cmdArgs := []string{
			"create-env",
			"bosh-deployment/bosh.yml",
			"--state", "bosh-director/state.json",
			"--vars-store", varsPath,
			"--ops-file", fmt.Sprintf("bosh-deployment/%s/cpi.yml", env.IaaS),
			"--ops-file", "bosh-deployment/bosh-lite.yml",
			"--ops-file", "bosh-deployment/external-ip-not-recommended.yml",
			"--ops-file", "bosh-deployment/jumpbox-user.yml",
			"--ops-file", "bosh-compiled-releases/ci/tasks/create-bosh-director/ops/remove-cpi-jobs.yml",
			// "--ops-file", "bosh-compiled-releases/ci/tasks/create-bosh-director/without-persistent-disk.yml",
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

	{ // envrc
		fh, err := os.OpenFile("bosh-director/.envrc", os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Panicf("opening .envrc: %v", err)
		}

		cmd := exec.Command("bosh", "interpolate", "--vars-file", varsPath, "--vars-file", os.Args[1], "--path", "/envrc", "-")
		cmd.Stdin = bytes.NewBufferString(`envrc: |
  export BOSH_ENVIRONMENT="((vars.external_ip))"
  export BOSH_CA_CERT="((director_ssl.ca))"
  export BOSH_CLIENT="admin"
  export BOSH_CLIENT_SECRET="((admin_password))"
`)
		cmd.Stdout = fh
		cmd.Stderr = os.Stderr

		err = cmd.Start()
		if err != nil {
			log.Panicf("executing envrc interpolate: %v", err)
		}

		err = cmd.Wait()
		if err != nil {
			log.Panicf("waiting for envrc interpolate: %v", err)
		}

		err = fh.Close()
		if err != nil {
			log.Panicf("writing envrc: %v", err)
		}
	}

	{ // bosh
		err = ioutil.WriteFile("bosh-director/bosh", []byte(`#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

source "$DIR/.envrc"
exec bosh "$@"`), 0750)
		if err != nil {
			log.Panicf("writing bosh: %v", err)
		}
	}

	{ // cloud-config
		cmdArgs := []string{
			"update-cloud-config", "-n",
			"bosh-deployment/warden/cloud-config.yml",
			"--vars-store", varsPath,
			"--var", fmt.Sprintf("director_name=bosh-lite-%s", now.Format("20060102T150405")),
		}

		for varKey, varVal := range env.Vars {
			cmdArgs = append(cmdArgs, "--var", fmt.Sprintf("%s=%s", varKey, varVal))
		}

		cmd := exec.Command("bosh-director/bosh", cmdArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Start()
		if err != nil {
			log.Panicf("executing update-cloud-config: %v", err)
		}

		err = cmd.Wait()
		if err != nil {
			log.Panicf("waiting for update-cloud-config: %v", err)
		}
	}
}
