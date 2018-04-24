package boshua

type Subject interface {
	SubjectReference() Reference
	SubjectMetalinkStorage() map[string]interface{}
}

type Reference struct {
	Context string
	ID      string
}
