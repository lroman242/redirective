package tracer

// ScreenSize type describe common screen size to capture correct size screenshot.
type ScreenSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// NewScreenSize convert width and height value to ScreenSize type.
func NewScreenSize(width, height int) *ScreenSize {
	return &ScreenSize{
		Width:  width,
		Height: height,
	}
}
