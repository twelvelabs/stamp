package iostreams

import "fmt"

type Logger interface {
	Info(string, ...any)
	Success(string, ...any)
	Warning(string, ...any)
	Failure(string, ...any)
}

var (
	// ensure interface
	_ Logger = &IconLogger{}
)

func NewIconLogger(ios *IOStreams, formatter *Formatter) *IconLogger {
	return &IconLogger{
		ios:       ios,
		formatter: formatter,
	}
}

type IconLogger struct {
	ios       *IOStreams
	formatter *Formatter
}

func (l *IconLogger) Info(line string, args ...any) {
	icon := l.formatter.InfoIcon()
	fmt.Fprintf(l.ios.Err, (icon + " " + line), args...)
}

func (l *IconLogger) Success(line string, args ...any) {
	icon := l.formatter.SuccessIcon()
	fmt.Fprintf(l.ios.Err, (icon + " " + line), args...)
}

func (l *IconLogger) Warning(line string, args ...any) {
	icon := l.formatter.WarningIcon()
	fmt.Fprintf(l.ios.Err, (icon + " " + line), args...)
}

func (l *IconLogger) Failure(line string, args ...any) {
	icon := l.formatter.FailureIcon()
	fmt.Fprintf(l.ios.Err, (icon + " " + line), args...)
}
