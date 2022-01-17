package log

import (
	"fmt"
	"log"

	"golangblog/internal/cfg"

	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
)

var red = color.New(color.FgRed).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()
var magenta = color.New(color.FgMagenta).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()

// Info - Log info
func Info(prefix string, data ...interface{}) {
	log.Println(cyan("ℹ ", "[", prefix, "]"), fmt.Sprint(data...))
}

// Warn - Log a warning
func Warn(prefix string, data ...interface{}) {
	log.Println(yellow("★ ", "[", prefix, "]"), fmt.Sprint(data...))
}

// Debug - Log debug info
func Debug(prefix string, data ...interface{}) {
	if cfg.IsDev() {
		log.Println(magenta("● ", "[", prefix, "]"), fmt.Sprint(data...))
	}
}

// Error - Log an error
func Error(prefix string, err error) {
	code := 0
	// retreive the custom statuscode if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	log.Println(red("⚠ ", "[", prefix, "/", code, "]"), err.Error())
}

// Fatal - Log fatal error and panics
func Fatal(prefix string, err error) {
	Error(prefix, err)
	panic(err)
}
