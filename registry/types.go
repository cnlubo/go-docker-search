package registry

type Environment struct {
	StorePath  string
	DockerID   string
	DockerPass string
}

type RepositoryRequest struct {
	RepositoryUrl string
	User          string
	Password      string
	Repo          string
	Endpoint      string
	PageSize      uint8
	Major         int64
	Minor       int64
}

type TagResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// type tagsResponse struct {
// 	Tags []string `json:"tags"`
// }

type SearchResults struct {
	Results  dockerContainers `json:"results"`
	Query    string           `json:"query"`
	NumPages int              `json:"num_pages"`
}
type dockerContainers []*dockerImage
type dockerImage struct {
	Description string
	IsOfficial  bool `json:"is_official"`
	IsTrusted   bool `json:"is_trusted"`
	Name        string
	StarCount   int `json:"star_count"`
	// Dockerfile string `json:"dockerfile"`
}

func (I dockerContainers) Len() int {
	return len(I)
}
func (I dockerContainers) Less(i, j int) bool {
	return I[i].StarCount > I[j].StarCount
}
func (I dockerContainers) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}
