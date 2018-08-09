package result

type Record struct {
	Raw    string      `json:"raw"`
	Parsed interface{} `json:"parsed"`
}
