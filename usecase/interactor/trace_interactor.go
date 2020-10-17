package interactor

import (
	"github.com/lroman242/redirective/domain"
	"github.com/lroman242/redirective/infrastructure/logger"
	"github.com/lroman242/redirective/infrastructure/tracer"
	"github.com/lroman242/redirective/usecase/presenter"
	"github.com/lroman242/redirective/usecase/repository"
	"math/rand"
	"net/url"
	"time"
)

// TraceInteractor represent interactor for actions related to trace results
type TraceInteractor interface {
	Trace(*url.URL, string) (*domain.TraceResults, error)
	Screenshot(*url.URL, int, int, string) (string, error)
	FindTrace(interface{}) (*domain.TraceResults, error)
}

// traceInteractor implement TraceInteractor interface
type traceInteractor struct {
	Tracer     tracer.Tracer
	Log        logger.Logger
	Repository repository.TraceRepository
	Presenter  presenter.TracePresenter
}

// NewTraceInteractor will construct TraceInteractor
func NewTraceInteractor(tracer tracer.Tracer, presenter presenter.TracePresenter, repo repository.TraceRepository, log logger.Logger) TraceInteractor {
	return &traceInteractor{
		Tracer:     tracer,
		Log:        log,
		Repository: repo,
		Presenter:  presenter,
	}
}

// Trace func will trace provided url
// and will return trace results (including screenshot)
func (ti *traceInteractor) Trace(url *url.URL, fileName string) (*domain.TraceResults, error) {
	results, err := ti.Tracer.Trace(url, fileName)
	if err != nil {
		return nil, err
	}

	id, err := ti.Repository.SaveTraceResults(results)
	if err != nil {
		return nil, err
	}

	//TODO: return id

	return ti.Presenter.ResponseTraceResults(results), err
}

// Screenshot function will make screenshot of landing url
func (ti *traceInteractor) Screenshot(url *url.URL, width int, height int, assetsFolderPath string) (string, error) {
	screenSize := &tracer.ScreenSize{
		Width:  width,
		Height: height,
	}

	path, err := ti.Tracer.Screenshot(url, screenSize, assetsFolderPath+"/"+ti.randomScreenshotFileName())
	if err != nil {
		return "", err
	}

	path = ti.Presenter.ResponseScreenshot(path)

	return path, nil
}

// FindTrace function will search and return trace results using provided ID
func (ti *traceInteractor) FindTrace(id interface{}) (*domain.TraceResults, error) {
	results, err := ti.Repository.FindTraceResults(id)
	if err != nil {
		return nil, err
	}

	return ti.Presenter.ResponseTraceResults(results), err
}

func (ti *traceInteractor) randomScreenshotFileName() string {
	var charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 16)

	for i := range b {
		b[i] = charset[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(charset))]
	}

	return string(b) + `.png`
}