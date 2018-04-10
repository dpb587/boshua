package analysis

import (
	"encoding/json"
	"io"
)

type Analyzer interface {
	Analyze(Writer) error
}

type Writer interface {
	Write(result interface{}) error
}

type JSONWriter struct {
	encoder *json.Encoder
}

func NewJSONWriter(writer io.Writer) *JSONWriter {
	return &JSONWriter{
		encoder: json.NewEncoder(writer),
	}
}

func (w *JSONWriter) Write(result interface{}) error {
	return w.encoder.Encode(result)
}
