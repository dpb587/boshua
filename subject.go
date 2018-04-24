package boshua

type Subject interface {
	SubjectReference() Reference
}

type Reference struct {
	Context string
	ID      string
}
