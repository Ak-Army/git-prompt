package color

import (
	"fmt"
	"strconv"
)
// Foreground text colors
const (
	FgBlack int = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

/// Formatter ...
type DefaultColoredOutput struct{}

var colorMap = map[string]int{
	"black":   FgBlack,
	"red":     FgRed,
	"green":   FgGreen,
	"yellow":  FgYellow,
	"blue":    FgBlue,
	"magenta": FgMagenta,
	"cyan":    FgCyan,
	"white":   FgWhite,
}

/// Foreground ...
func (f DefaultColoredOutput) foreground(color string, format string, a ...interface{}) string {
	return "\001\033[" + f.getColor(color) + "m\002" + fmt.Sprintf(format, a...) + "\001\033[0m\002"
}
func (f DefaultColoredOutput) getColor(color string) string {
	return strconv.Itoa(colorMap[color])
}
