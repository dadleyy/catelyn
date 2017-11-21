package catelyn

type ConfluenceSpace struct {
	ID   int    `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
	Kind string `json:"type"`
}

type ConfluencePaging struct {
	Limit uint8 `json:"limit"`
	Start uint8 `json:"start"`
	Size  uint8 `json:"size"`
	Links struct {
		Next     string `json:"next"`
		Previous string `json:"previous"`
	} `json:"_links"`
}
