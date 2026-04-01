# Pog

[![CI](https://github.com/FurqanSoftware/pog/actions/workflows/ci.yml/badge.svg)](https://github.com/FurqanSoftware/pog/actions/workflows/ci.yml)

Pog is a simple logger for Go with a status indicator.

[![Go Reference](https://pkg.go.dev/badge/github.com/FurqanSoftware/pog.svg)](https://pkg.go.dev/github.com/FurqanSoftware/pog)

We use Pog in [Toph Printd](https://github.com/FurqanSoftware/toph-printd).

## Usage

``` go
import (
	"github.com/FurqanSoftware/pog"
	"github.com/fatih/color"
)

type pogStatus struct {
	icon  byte
	text  string
	color *color.Color
	throb bool
}

func (s pogStatus) Icon() byte          { return s.icon }
func (s pogStatus) Text() string        { return s.text }
func (s pogStatus) Color() *color.Color { return s.color }
func (s pogStatus) Throb() bool         { return s.throb }

var (
	statusReady    = pogStatus{'~', "Ready", color.New(color.FgGreen), true}
	statusPrinting = pogStatus{'~', "Printing", color.New(color.FgBlue), false}
	statusOffline  = pogStatus{'!', "Offline", color.New(color.FgRed), false}
)
```

``` go
log.SetPrefix("\033[2K\r")
log.SetFlags(log.Ldate | log.Ltime)

pog.InitDefault()

pog.Info("I am ready.")
pog.SetStatus(statusReady)
```

``` txt
2023/10/18 11:59:11 [i] I am ready.
                    [~] Ready
```

``` go
pog.Info("Got a print.")
pog.SetStatus(statusPrinting)
pog.Warn("That's a lot of words.")
```

``` txt
2023/10/18 11:59:11 [i] I am ready.
2023/10/18 11:59:19 [i] Got a print.
2023/10/18 11:59:19 [w] That's a lot of words.
                    [~] Printing
```

``` go
pog.Error("Lost connection.")
pog.SetStatus(statusOffline)
```

``` txt
2023/10/18 11:59:11 [i] I am ready.
2023/10/18 11:59:19 [i] Got a print.
2023/10/18 11:59:19 [w] That's a lot of words.
2023/10/18 11:59:57 [E] Lost connection.
                    [!] Offline
```

## Log Levels

| Level | Indicator | Variants |
|-------|-----------|----------|
| Debug | `[d]` | `Debug`, `Debugln`, `Debugf` |
| Info  | `[i]` | `Info`, `Infoln`, `Infof` |
| Warn  | `[w]` | `Warn`, `Warnln`, `Warnf` |
| Error | `[E]` | `Error`, `Errorln`, `Errorf` |
| Fatal | `[E]` | `Fatal`, `Fatalln`, `Fatalf` |

Fatal functions log at the Error level and then exit, running any registered exit hooks.

## Exit Hooks

Register functions to run before a `Fatal` exit:

``` go
pog.AddExitHook(func() {
	// Clean up resources.
})
```
