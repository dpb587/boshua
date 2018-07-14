package clicommon

import (
	"fmt"

	"github.com/pkg/errors"
)

type AnalyzersCmd struct{}

func (c *AnalyzersCmd) Execute(loader SubjectLoader) error {
	subject, err := loader()
	if err != nil {
		return errors.Wrap(err, "finding subject")
	}

	for _, analyzer := range subject.SupportedAnalyzers() {
		fmt.Println(analyzer)
	}

	return nil
}
