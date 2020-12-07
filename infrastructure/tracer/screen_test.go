package tracer_test

import (
	"testing"

	"github.com/lroman242/redirective/infrastructure/tracer"
)

func TestNewScreenSize(t *testing.T) {
	ss := tracer.NewScreenSize(15, 25)

	if ss.Width != 15 {
		t.Errorf("Invalid ScreenSize Width on creating. Expect %d but get %d", 15, ss.Width)
	}

	if ss.Height != 25 {
		t.Errorf("Invalid ScreenSize Height on creating. Expect %d but get %d", 25, ss.Height)
	}
}
