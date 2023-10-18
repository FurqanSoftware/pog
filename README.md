# Pog

Pog is a simple logger for Go with a status indicator.

## Usage

``` go
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
