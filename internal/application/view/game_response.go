package view

type ImageType string

const (
	ImageNone   ImageType = "NONE"
	ImageAsset  ImageType = "ASSET"
	ImageBytes  ImageType = "BYTES"
	ImageBase64 ImageType = "BASE64"
)

type Image struct {
	Type   ImageType
	Path   string
	Data   []byte
	Base64 string
	Name   string
}

type GameResponse struct {
	Text       string
	Image      *Image
	MenuKey    string
	MenuParams map[string]any
}
