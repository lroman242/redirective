package screenshoter

import (
	"github.com/lroman242/redirective/phantomjs"
	//"github.com/benbjohnson/phantomjs"
	"net/url"
	"time"
)

type ScreenSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type ScreenShoter struct {
	Url     *url.URL     `json:"url"`
	Size    *ScreenSize `json:"size"`
	Timeout int         `json:"timeout"`
	Image   string      `json:"image"`
}

func NewScreenShot(url *url.URL, width int, height int, timeout int) *ScreenShoter {
	return &ScreenShoter{
		Size:    &ScreenSize{Width: width, Height: height},
		Timeout: timeout,
		Url:     url,
	}
}

func (ss *ScreenShoter) CaptureScreenShot(imagePath string) error {
	page, err := phantomjs.DefaultProcess.CreateWebPage()
	if err != nil {
		return err
	}
	defer page.Close()

	err = page.SetViewportSize(ss.Size.Width, ss.Size.Height)
	if err != nil {
		return err
	}

	//err = page.SetClipRect(phantomjs.Rect{0,0,1366,768})
	//if err != nil {
	//	return err
	//}

	if err := page.Open(ss.Url.String()); err != nil {
		return err
	}
	time.Sleep(time.Duration(ss.Timeout) * time.Second)

	err = page.Render(imagePath, "PNG", 100)
	if err != nil {
		return err
	}

	ss.Image = imagePath

	return nil
}
