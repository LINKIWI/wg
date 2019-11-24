package webgrep

const (
	// EndpointSearch is the path to the code search query endpoint.
	EndpointSearch = "/api/search"
	// EndpointMetadata is the path to the metadata endpoint.
	EndpointMetadata = "/api/meta/info"
)

// CodeSearchResult formalizes fields for a single code search result.
type CodeSearchResult struct {
	Repo    string `json:"repo"`
	Version string `json:"version"`
	Path    string `json:"path"`
	Lines   []struct {
		Line   string `json:"line"`
		Number int    `json:"number"`
		Bounds []int  `json:"bounds"`
	} `json:"lines"`
}

// FileSearchResult formalizes fields for a single file path result.
type FileSearchResult struct {
	Repo    string `json:"repo"`
	Version string `json:"version"`
	Path    string `json:"path"`
	Bounds  []int  `json:"bounds"`
}

// SearchStats formalizes fields in server-side search statistics.
type SearchStats struct {
	ExitReason int `json:"exitReason"`
	Time       int `json:"time"`
}

// SearchQueryRequest describes the top-level request for a search query.
type SearchQueryRequest struct {
	Query         string   `json:"query"`
	File          string   `json:"file"`
	Repos         []string `json:"repos"`
	Regex         bool     `json:"regex"`
	CaseSensitive bool     `json:"caseSensitive"`
	MaxMatches    int      `json:"maxMatches"`
	Context       int      `json:"context"`
}

// SearchQueryResponse describes the top-level response for a search query.
type SearchQueryResponse struct {
	Stats SearchStats        `json:"stats"`
	Code  []CodeSearchResult `json:"code"`
	Files []FileSearchResult `json:"files"`
}

// MetadataResponse describes the top-level response for a metadata request.
type MetadataResponse struct {
	Name         string `json:"name"`
	Timestamp    int    `json:"timestamp"`
	Version      string `json:"version"`
	Repositories []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		URL     string `json:"url"`
		Remote  string `json:"remote"`
	} `json:"repositories"`
}
