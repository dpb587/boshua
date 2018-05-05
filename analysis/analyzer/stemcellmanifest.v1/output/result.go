package output

type Result struct {
	Raw    string      `json:"raw"`
	Parsed interface{} `json:"parsed"`
}
