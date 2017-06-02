package color

type Colors interface {
	foreground(color string, format string, a ...interface{}) string
}
type Color struct {
	Colors
}

func (f Color) Blue(format string, a ...interface{}) string {
	return f.foreground("blue", format, a...)
}

func (f Color) Cyan(format string, a ...interface{}) string {
	return f.foreground("cyan", format, a...)
}

func (f Color) Yellow(format string, a ...interface{}) string {
	return f.foreground("yellow", format, a...)
}

func (f Color) Green(format string, a ...interface{}) string {
	return f.foreground("green", format, a...)
}

func (f Color) Red(format string, a ...interface{}) string {
	return f.foreground("red", format, a...)
}
