package color

import "fmt"

/// Formatter ...
type ZshColoredOutput struct{}

/// Foreground ...
func (f ZshColoredOutput) foreground(color string, format string, a ...interface{}) string {
	return "%F{" + color + "}" + fmt.Sprintf(format, a...) + "%f"
}
