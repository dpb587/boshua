package datastore

import (
	"github.com/pkg/errors"
)

type StoreAnalysisCmd struct {
	*CmdOpts `no-flag:"true"`

	Analyzer string `long:"analyzer" description:"The analyzer used for the results"`

	Args StoreAnalysisArgs `positional-args:"true"`
}

type StoreAnalysisArgs struct {
	Artifact string `positional-arg-name:"ARTIFACT-PATH" description:"Artifact to store"`
}

func (c *StoreAnalysisCmd) Execute(_ []string) error {
	return errors.New("TODO resurrect functionality")
	// c.AppOpts.ConfigureLogger("compiledrelease/datastore/store-analysis")
	//
	// index, err := c.getDatastore()
	// if err != nil {
	// 	return errors.Wrap(err, "loading datastore")
	// }
	//
	// subject, err := index.Find(c.CompiledReleaseOpts.Reference())
	// if err != nil {
	// 	return errors.Wrap(err, "filtering")
	// }
	//
	// analysisIndex, err := index.GetAnalysisDatastore(c.CompiledReleaseOpts.Reference())
	// if err != nil {
	// 	return errors.Wrap(err, "getting analysis index")
	// }
	//
	// ref := analysis.Reference{
	// 	Subject:  subject,
	// 	Analyzer: analysis.AnalyzerName(c.Analyzer),
	// }
	//
	// path, err := filepath.Abs(c.Args.Artifact)
	// if err != nil {
	// 	return errors.Wrap(err, "expanding artifact path")
	// }
	//
	// meta4, err := metalinkutil.CreateFromFiles(fmt.Sprintf("file://%s", path))
	// if err != nil {
	// 	return errors.Wrap(err, "creating in-memory metalink")
	// }
	//
	// err = analysisIndex.Store(ref, *meta4)
	// if err != nil {
	// 	return errors.Wrap(err, "storing artifact")
	// }
	//
	// return nil
}
