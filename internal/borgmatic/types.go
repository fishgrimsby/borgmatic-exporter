package borgmatic

type ListResult struct {
	Archives   []Archive  `json:"archives"`
	Encryption Encryption `json:"encryption"`
	Repository Repository `json:"repository"`
}

type Repository struct {
	Id           string `json:"id"`
	LastModified string `json:"last_modified"`
	Location     string `json:"location"`
}

type Encryption struct {
	Mode string `json:"mode"`
}

type Archive struct {
	Archive  string `json:"archive"`
	Barchive string `json:"barchive"`
	Id       string `json:"id"`
	Name     string `json:"name"`
	Start    string `json:"start"`
	Time     string `json:"time"`
}

type InfoResult struct {
	Cache       Cache      `json:"cache"`
	Encryption  Encryption `json:"encryption"`
	Repository  Repository `json:"repository"`
	SecurityDir string     `json:"security_dir"`
}

type Cache struct {
	Path  string `json:"path"`
	Stats Stats  `json:"stats"`
}

type Stats struct {
	TotalChunks       int `json:"total_chunks"`
	TotalCsize        int `json:"total_csize"`
	TotalSize         int `json:"total_size"`
	TotalUniqueChunks int `json:"total_unique_chunks"`
	UniqueCsize       int `json:"unique_csize"`
	UniqueSize        int `json:"unique_size"`
}
