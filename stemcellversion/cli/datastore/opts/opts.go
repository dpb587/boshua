package opts

type Opts struct {
	Datastore string `long:"datastore" description:"The datastore name to use" default:"default"`
}
