package catelyn

// ConfluenceSpace represents a space object in confluence.
type ConfluenceSpace struct {
	ID   int    `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
	Kind string `json:"type"`
}

// ConfluenceContent represents individual json items returned from the content api.
type ConfluenceContent struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

// ConfluencePaging represents paging information returned from the confluence rest api.
type ConfluencePaging struct {
	Limit uint8 `json:"limit"`
	Start uint8 `json:"start"`
	Size  uint8 `json:"size"`
	Links struct {
		Next     string `json:"next"`
		Previous string `json:"previous"`
	} `json:"_links"`
}

// ConfluencePageSearchInput represents the space search query parameters.
type ConfluencePageSearchInput struct {
	ConfluencePaging
	Title    string
	SpaceKey string
}

// ConfluenceSpaceSearchInput represents the space search query parameters.
type ConfluenceSpaceSearchInput struct {
	ConfluencePaging
	Type string
}
