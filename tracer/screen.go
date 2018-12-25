package tracer

type screenSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func NewScreenSize(width, height int) *screenSize {
	return &screenSize{
		Width:  width,
		Height: height,
	}
}
