package pog

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func init() {
	color.NoColor = true
}

type testStatus struct {
	icon  byte
	text  string
	color *color.Color
	throb bool
}

func (s testStatus) Icon() byte          { return s.icon }
func (s testStatus) Text() string        { return s.text }
func (s testStatus) Color() *color.Color { return s.color }
func (s testStatus) Throb() bool         { return s.throb }

func newTestPogger() (*Pogger, *bytes.Buffer) {
	var buf bytes.Buffer
	p := NewPogger(&buf, "", 0)
	return p, &buf
}

func TestDebug(t *testing.T) {
	p, buf := newTestPogger()
	defer p.Stop()

	p.Debug("hello")
	if !strings.Contains(buf.String(), "[d] hello") {
		t.Errorf("got %q, want it to contain %q", buf.String(), "[d] hello")
	}
}

func TestDebugln(t *testing.T) {
	p, buf := newTestPogger()
	defer p.Stop()

	p.Debugln("hello")
	if !strings.Contains(buf.String(), "[d] hello") {
		t.Errorf("got %q, want it to contain %q", buf.String(), "[d] hello")
	}
}

func TestDebugf(t *testing.T) {
	p, buf := newTestPogger()
	defer p.Stop()

	p.Debugf("hello %s", "world")
	if !strings.Contains(buf.String(), "[d] hello world") {
		t.Errorf("got %q, want it to contain %q", buf.String(), "[d] hello world")
	}
}

func TestInfo(t *testing.T) {
	p, buf := newTestPogger()
	defer p.Stop()

	p.Info("hello")
	if !strings.Contains(buf.String(), "[i] hello") {
		t.Errorf("got %q, want it to contain %q", buf.String(), "[i] hello")
	}
}

func TestInfoln(t *testing.T) {
	p, buf := newTestPogger()
	defer p.Stop()

	p.Infoln("hello")
	if !strings.Contains(buf.String(), "[i] hello") {
		t.Errorf("got %q, want it to contain %q", buf.String(), "[i] hello")
	}
}

func TestInfof(t *testing.T) {
	p, buf := newTestPogger()
	defer p.Stop()

	p.Infof("hello %s", "world")
	if !strings.Contains(buf.String(), "[i] hello world") {
		t.Errorf("got %q, want it to contain %q", buf.String(), "[i] hello world")
	}
}

func TestWarn(t *testing.T) {
	p, buf := newTestPogger()
	defer p.Stop()

	p.Warn("hello")
	if !strings.Contains(buf.String(), "[w] hello") {
		t.Errorf("got %q, want it to contain %q", buf.String(), "[w] hello")
	}
}

func TestWarnln(t *testing.T) {
	p, buf := newTestPogger()
	defer p.Stop()

	p.Warnln("hello")
	if !strings.Contains(buf.String(), "[w] hello") {
		t.Errorf("got %q, want it to contain %q", buf.String(), "[w] hello")
	}
}

func TestWarnf(t *testing.T) {
	p, buf := newTestPogger()
	defer p.Stop()

	p.Warnf("hello %s", "world")
	if !strings.Contains(buf.String(), "[w] hello world") {
		t.Errorf("got %q, want it to contain %q", buf.String(), "[w] hello world")
	}
}

func TestError(t *testing.T) {
	p, buf := newTestPogger()
	defer p.Stop()

	p.Error("hello")
	if !strings.Contains(buf.String(), "[E] hello") {
		t.Errorf("got %q, want it to contain %q", buf.String(), "[E] hello")
	}
}

func TestErrorln(t *testing.T) {
	p, buf := newTestPogger()
	defer p.Stop()

	p.Errorln("hello")
	if !strings.Contains(buf.String(), "[E] hello") {
		t.Errorf("got %q, want it to contain %q", buf.String(), "[E] hello")
	}
}

func TestErrorf(t *testing.T) {
	p, buf := newTestPogger()
	defer p.Stop()

	p.Errorf("hello %s", "world")
	if !strings.Contains(buf.String(), "[E] hello world") {
		t.Errorf("got %q, want it to contain %q", buf.String(), "[E] hello world")
	}
}

func TestSetStatus(t *testing.T) {
	p, _ := newTestPogger()
	defer p.Stop()

	s := testStatus{icon: '~', text: "Ready", throb: false}
	p.SetStatus(s)

	p.m.Lock()
	got := p.status
	p.m.Unlock()

	if got == nil {
		t.Fatal("status is nil after SetStatus")
	}
	if got.Text() != "Ready" {
		t.Errorf("got text %q, want %q", got.Text(), "Ready")
	}
	if got.Icon() != '~' {
		t.Errorf("got icon %q, want %q", got.Icon(), '~')
	}
}

func TestAddExitHook(t *testing.T) {
	p, _ := newTestPogger()
	defer p.Stop()

	called := false
	p.AddExitHook(func() { called = true })

	if len(p.exitHooks) != 1 {
		t.Fatalf("got %d exit hooks, want 1", len(p.exitHooks))
	}

	p.exitHooks[0]()
	if !called {
		t.Error("exit hook was not called")
	}
}

func TestStop(t *testing.T) {
	p, _ := newTestPogger()
	p.Stop()

	// Verify stopCh is closed by reading from it (should not block).
	select {
	case <-p.stopCh:
	default:
		t.Error("stopCh not closed after Stop()")
	}
}

func TestDefaultPogger(t *testing.T) {
	InitDefault()
	defer Stop()

	if Default() == nil {
		t.Fatal("Default() returned nil after InitDefault()")
	}
}

func TestDefaultPackageFunctions(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	log.SetOutput(logger.Writer())
	log.SetPrefix("")
	log.SetFlags(0)

	InitDefault()
	defer Stop()

	Info("pkg-info")
	if !strings.Contains(buf.String(), "[i] pkg-info") {
		t.Errorf("got %q, want it to contain %q", buf.String(), "[i] pkg-info")
	}
}

func TestLogWithFlags(t *testing.T) {
	var buf bytes.Buffer
	p := NewPogger(&buf, "", log.Ldate|log.Ltime)
	defer p.Stop()

	p.Info("flagged")
	if !strings.Contains(buf.String(), "[i] flagged") {
		t.Errorf("got %q, want it to contain %q", buf.String(), "[i] flagged")
	}
}
