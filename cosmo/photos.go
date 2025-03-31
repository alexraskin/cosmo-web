package cosmo

type ImageConfig struct {
	BaseURL         string
	ThumbnailParams string
	FullsizeParams  string
	Folder          string
}

type CatImage struct {
	Name string
}

type GalleryData struct {
	Title       string
	Description string
	Images      []CatImage
	ImageConfig ImageConfig
}

func DefaultGalleryConfig() GalleryData {
	return GalleryData{
		Title:       "Cosmo the Cat",
		Description: "a collection of photos of my adorable cat Cosmo",
		Images:      defaultCatImages,
		ImageConfig: defaultImageConfig,
	}
}

// Default configurations
var (
	defaultImageConfig = ImageConfig{
		BaseURL:         "https://cdn.cosmothecat.net/cdn-cgi/image",
		ThumbnailParams: "width=300,height=200,fit=cover,quality=80",
		FullsizeParams:  "quality=80",
		Folder:          "cosmo",
	}

	defaultCatImages = []CatImage{
		{Name: "0f1z8b8s6eqakqa9zdxk7wrh13zzae36.jpg"},
		{Name: "0ghgth0axzr1wcx355e6qvg6a07y1fz3.jpg"},
		{Name: "1h2qf5sae9mhgd877k1dhvmzqf7yc7ff.jpg"},
		{Name: "5yp8ehxxxxaemkqg64cwpz68tffycdpt.jpg"},
		{Name: "IMG_0025-optimized.jpg"},
		{Name: "IMG_0122-optimized.jpg"},
		{Name: "IMG_0196-optimized.jpg"},
		{Name: "IMG_0491-optimized.jpg"},
		{Name: "IMG_1263-optimized.jpg"},
		{Name: "IMG_1345-optimized.jpg"},
		{Name: "IMG_1385-optimized.jpg"},
		{Name: "IMG_1585-optimized.jpg"},
		{Name: "IMG_1599-optimized.jpg"},
		{Name: "IMG_1704-optimized.jpg"},
		{Name: "IMG_1707-optimized.jpg"},
		{Name: "IMG_1828-optimized.jpg"},
		{Name: "IMG_1919-optimized.jpg"},
		{Name: "bcdf5krb6yb9jn30thaqtv7n35w87t26.jpg"},
		{Name: "bs4r87sxhh57a1t6m5jvqr4kzf7jwx0a.jpg"},
		{Name: "dcp3q6ptv45wgb057c0759d5xztqnfeq.jpg"},
		{Name: "drmx5zftnz6f72dd1j8v7mvxzra49kbj.jpg"},
		{Name: "g6ydxahn69pf9mvxqkrc7xw9xyh4yk1r.jpg"},
		{Name: "yvpqzd77q3xs1qr8vsn0077ehk2gg2mz.jpg"},
		{Name: "IMG_20250221_210034.JPEG"},
		{Name: "IMG_1385.JPEG"},
	}
)
