package cosmo

type ErrorResponse struct {
	Message   string `json:"message"`
	Status    int    `json:"status"`
	Path      string `json:"path"`
	RequestID string `json:"request_id"`
}

type ImageConfig struct {
	BaseURL         string
	ThumbnailParams string
	FullsizeParams  string
	Folder          string
}

type CatImage struct {
	URL  string
	Name string
}

type GalleryData struct {
	Title       string
	Description string
	Images      []CatImage
	ImageConfig ImageConfig
}
