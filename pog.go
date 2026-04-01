// Package pog provides a simple logger for Go with a terminal status indicator.
//
// Pog wraps the standard log package, adding colored log level indicators and a
// persistent status line that is redrawn at the bottom of the terminal.
package pog

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Pogger is a logger with a terminal status indicator. It wraps a [log.Logger]
// and adds leveled logging and a persistent status line.
type Pogger struct {
	out       io.Writer
	logger    *log.Logger
	status    Status
	stopCh    chan struct{}
	m         sync.Mutex
	initOnce  sync.Once
	exitHooks []func()
}

// NewPogger creates a new Pogger that writes to out. The prefix and flag
// arguments are passed to [log.New] for the underlying logger. A background
// goroutine is started to manage the status line; call [Pogger.Stop] to stop it.
func NewPogger(out io.Writer, prefix string, flag int) *Pogger {
	pogger := Pogger{
		out:    out,
		logger: log.New(out, prefix, flag),
		stopCh: make(chan struct{}),
	}
	go pogger.loop()
	return &pogger
}

// SetStatus sets the status displayed on the terminal status line.
func (p *Pogger) SetStatus(status Status) {
	p.m.Lock()
	p.status = status
	p.m.Unlock()
}

// Debug logs a message with the debug level indicator [d].
func (p *Pogger) Debug(v ...any) {
	a := []any{color.WhiteString("[d]") + " "}
	a = append(a, v...)
	p.logger.Print(a...)
}

// Debugln logs a message with the debug level indicator [d], followed by a newline.
func (p *Pogger) Debugln(v ...any) {
	a := []any{color.WhiteString("[d]")}
	a = append(a, v...)
	p.logger.Println(a...)
}

// Debugf logs a formatted message with the debug level indicator [d].
func (p *Pogger) Debugf(format string, v ...any) {
	p.logger.Printf(color.WhiteString("[d]")+" "+format, v...)
}

// Info logs a message with the info level indicator [i].
func (p *Pogger) Info(v ...any) {
	a := []any{"[i] "}
	a = append(a, v...)
	p.logger.Print(a...)
}

// Infoln logs a message with the info level indicator [i], followed by a newline.
func (p *Pogger) Infoln(v ...any) {
	a := []any{"[i]"}
	a = append(a, v...)
	p.logger.Println(a...)
}

// Infof logs a formatted message with the info level indicator [i].
func (p *Pogger) Infof(format string, v ...any) {
	p.logger.Printf("[i] "+format, v...)
}

// Warn logs a message with the warn level indicator [w].
func (p *Pogger) Warn(v ...any) {
	a := []any{color.YellowString("[w]") + " "}
	a = append(a, v...)
	p.logger.Print(a...)
}

// Warnln logs a message with the warn level indicator [w], followed by a newline.
func (p *Pogger) Warnln(v ...any) {
	a := []any{color.YellowString("[w]")}
	a = append(a, v...)
	p.logger.Println(a...)
}

// Warnf logs a formatted message with the warn level indicator [w].
func (p *Pogger) Warnf(format string, v ...any) {
	p.logger.Printf(color.YellowString("[w]")+" "+format, v...)
}

// Error logs a message with the error level indicator [E].
func (p *Pogger) Error(v ...any) {
	a := []any{color.RedString("[E]") + " "}
	a = append(a, v...)
	p.logger.Print(a...)
}

// Errorln logs a message with the error level indicator [E], followed by a newline.
func (p *Pogger) Errorln(v ...any) {
	a := []any{color.RedString("[E]")}
	a = append(a, v...)
	p.logger.Println(a...)
}

// Errorf logs a formatted message with the error level indicator [E].
func (p *Pogger) Errorf(format string, v ...any) {
	p.logger.Printf(color.RedString("[E]")+" "+format, v...)
}

// Fatal logs a message with the error level indicator [E] and then calls
// os.Exit(1), running any registered exit hooks first.
func (p *Pogger) Fatal(v ...any) {
	p.Error(v...)
	p.exit(1)
}

// Fatalln logs a message with the error level indicator [E] followed by a
// newline and then calls os.Exit(1), running any registered exit hooks first.
func (p *Pogger) Fatalln(v ...any) {
	p.Errorln(v...)
	p.exit(1)
}

// Fatalf logs a formatted message with the error level indicator [E] and then
// calls os.Exit(1), running any registered exit hooks first.
func (p *Pogger) Fatalf(format string, v ...any) {
	p.Errorf(format, v...)
	p.exit(1)
}

// AddExitHook registers a function to be called before os.Exit when a Fatal
// method is invoked.
func (p *Pogger) AddExitHook(fn func()) {
	p.m.Lock()
	p.exitHooks = append(p.exitHooks, fn)
	p.m.Unlock()
}

// Stop stops the background goroutine that manages the status line.
func (p *Pogger) Stop() {
	close(p.stopCh)
}

func (p *Pogger) loop() {
	cur := ""
	pad := ""
	{
		n := len(p.logger.Prefix())
		flags := p.logger.Flags()
		if flags&(log.Ldate|log.Ltime|log.Lmicroseconds) != 0 {
			if flags&log.Ldate != 0 {
				n += len("2006/01/02 ")
			}
			if flags&(log.Ltime|log.Lmicroseconds) != 0 {
				if flags&log.Lmicroseconds != 0 {
					n += len("15:04:05.000000 ")
				} else {
					n += len("15:04:05 ")
				}
			}
		}
		if n > 0 {
			pad = strings.Repeat(" ", n)
		}
	}
	ticker := time.NewTicker(125 * time.Millisecond)
	defer ticker.Stop()
L:
	for i := 0; ; i = (i + 1) % 10 {
		var s string
		p.m.Lock()
		b := []byte{' '}
		if p.status != nil {
			if i < 5 || !p.status.Throb() {
				b[0] = p.status.Icon()
			}
			s = "[" + string(b) + "]"
			if color := p.status.Color(); color != nil {
				s = color.Sprint(s)
			}
			s += " " + p.status.Text()
		}
		p.m.Unlock()
		if s != cur {
			fmt.Fprintf(p.out, "\033[2K\r%s%s\r", pad, s)
			cur = s
		}
		select {
		case <-p.stopCh:
			break L
		case <-ticker.C:
		}
	}
}

func (p *Pogger) exit(code int) {
	p.m.Lock()
	hooks := make([]func(), len(p.exitHooks))
	copy(hooks, p.exitHooks)
	p.m.Unlock()
	for _, fn := range hooks {
		fn()
	}
	os.Exit(code)
}

// Status represents the state displayed on the terminal status line.
type Status interface {
	// Icon returns the character shown inside the status indicator brackets.
	Icon() byte
	// Text returns the text displayed next to the status indicator.
	Text() string
	// Color returns the color used to render the status indicator, or nil for
	// no color.
	Color() *color.Color
	// Throb returns whether the status indicator should animate by toggling the
	// icon on and off.
	Throb() bool
}

var (
	defaultPogger *Pogger
)

// InitDefault initializes the default Pogger using the settings from
// [log.Default]. It must be called before using the package-level logging
// functions.
func InitDefault() {
	defaultPogger = NewPogger(log.Default().Writer(), log.Default().Prefix(), log.Default().Flags())
}

// Default returns the default Pogger. It returns nil if [InitDefault] has not
// been called.
func Default() *Pogger {
	return defaultPogger
}

// SetStatus sets the status on the default Pogger.
func SetStatus(status Status) { defaultPogger.SetStatus(status) }

// Debug calls [Pogger.Debug] on the default Pogger.
func Debug(v ...any) { defaultPogger.Debug(v...) }

// Debugln calls [Pogger.Debugln] on the default Pogger.
func Debugln(v ...any) { defaultPogger.Debugln(v...) }

// Debugf calls [Pogger.Debugf] on the default Pogger.
func Debugf(format string, v ...any) { defaultPogger.Debugf(format, v...) }

// Info calls [Pogger.Info] on the default Pogger.
func Info(v ...any) { defaultPogger.Info(v...) }

// Infoln calls [Pogger.Infoln] on the default Pogger.
func Infoln(v ...any) { defaultPogger.Infoln(v...) }

// Infof calls [Pogger.Infof] on the default Pogger.
func Infof(format string, v ...any) { defaultPogger.Infof(format, v...) }

// Warn calls [Pogger.Warn] on the default Pogger.
func Warn(v ...any) { defaultPogger.Warn(v...) }

// Warnln calls [Pogger.Warnln] on the default Pogger.
func Warnln(v ...any) { defaultPogger.Warnln(v...) }

// Warnf calls [Pogger.Warnf] on the default Pogger.
func Warnf(format string, v ...any) { defaultPogger.Warnf(format, v...) }

// Error calls [Pogger.Error] on the default Pogger.
func Error(v ...any) { defaultPogger.Error(v...) }

// Errorln calls [Pogger.Errorln] on the default Pogger.
func Errorln(v ...any) { defaultPogger.Errorln(v...) }

// Errorf calls [Pogger.Errorf] on the default Pogger.
func Errorf(format string, v ...any) { defaultPogger.Errorf(format, v...) }

// Fatal calls [Pogger.Fatal] on the default Pogger.
func Fatal(v ...any) { defaultPogger.Fatal(v...) }

// Fatalln calls [Pogger.Fatalln] on the default Pogger.
func Fatalln(v ...any) { defaultPogger.Fatalln(v...) }

// Fatalf calls [Pogger.Fatalf] on the default Pogger.
func Fatalf(format string, v ...any) { defaultPogger.Fatalf(format, v...) }

// Stop calls [Pogger.Stop] on the default Pogger.
func Stop() { defaultPogger.Stop() }

// AddExitHook calls [Pogger.AddExitHook] on the default Pogger.
func AddExitHook(fn func()) { defaultPogger.AddExitHook(fn) }
