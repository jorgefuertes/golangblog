package main

import (
	"embed"
	"io/fs"
	"net/http"
	"path/filepath"

	"golangblog/controller"
	"golangblog/internal/banner"
	"golangblog/internal/build"
	"golangblog/internal/cfg"
	"golangblog/internal/dev"
	"golangblog/internal/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/jet"
	"gopkg.in/alecthomas/kingpin.v2"
)

//go:embed asset
var asset embed.FS

//go:embed view
var view embed.FS

func main() {
	// command line flags and params
	cfg.Main.Env = kingpin.Flag("environment", "dev or prod mode").
		Short('e').Default("prod").String()
	cfg.Main.Server.IP = kingpin.Flag("ip", "IP to listen").
		Short('i').Default("127.0.0.1").String()
	cfg.Main.Server.Port = kingpin.Flag("port", "Port to listen").
		Short('p').Default("8080").String()
	cfg.Main.Server.BodyLimitMb = kingpin.Flag("body-limit", "Body limit in MiB").
		Default("4").Int()
	cfg.Main.Server.RTimeout = kingpin.Flag("read-timeout", "Read timeout").
		Short('r').Default("10s").Duration()
	cfg.Main.Server.WTimeout = kingpin.Flag("write-timeout", "Write timeout").
		Short('w').Default("10s").Duration()
	cfg.Main.Server.Concurrency = kingpin.Flag("concurrency", "Maximum number of concurrent connections in MiB").
		Default("256").Int()

	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(cfg.Version).Author(cfg.Author)
	kingpin.CommandLine.Help = "Web Application Server"
	kingpin.Parse()

	// root path
	cfg.Main.Root = filepath.Dir(".")

	// template engine
	var engine *jet.Engine
	if cfg.IsDev() {
		engine = jet.New("view", ".jet.html")
	} else {
		tFS, _ := fs.Sub(view, "view")
		engine = jet.NewFileSystem(http.FS(tFS), ".jet.html")
	}

	if err := engine.Load(); err != nil {
		log.Fatal("TMPL:ENGINE", err)
	}
	if cfg.IsDev() {
		engine.Reload(true)
	}

	engine.AddFunc("isDev", cfg.IsDev)
	engine.AddFunc("isProd", cfg.IsProd)

	// app and configuration
	app := fiber.New(
		fiber.Config{
			ReadTimeout:           *cfg.Main.Server.RTimeout,
			WriteTimeout:          *cfg.Main.Server.WTimeout,
			BodyLimit:             *cfg.Main.Server.BodyLimitMb * 1024 * 1024,
			Concurrency:           *cfg.Main.Server.Concurrency * 1024,
			ServerHeader:          "QRRadionics_Server_" + cfg.Version,
			DisableStartupMessage: true,
			Views:                 engine,
			ViewsLayout:           "./layout/main",
			ErrorHandler: func(ctx *fiber.Ctx, err error) error {
				code := fiber.StatusInternalServerError
				txt := banner.Error500
				if e, ok := err.(*fiber.Error); ok {
					code = e.Code
				}
				if code == 404 {
					txt = banner.Error404
				}
				err = ctx.Status(code).SendString(
					banner.Title + "\n" + txt + banner.Separator + err.Error() + banner.Separator)
				if err != nil {
					return ctx.Status(500).SendString("Internal Server Error")
				}

				return nil
			},
		},
	)

	// compression
	if cfg.IsProd() {
		app.Use(compress.New())
	}

	// logger
	app.Use(logger.New())

	// cors
	app.Use(cors.New())

	// routes
	controller.Setup(app)

	// static assets
	if cfg.IsDev() {
		app.Static("/asset", "asset")
	} else {
		subFS, _ := fs.Sub(asset, "asset")
		app.Use("/asset", filesystem.New(filesystem.Config{
			Root:         http.FS(subFS),
			NotFoundFile: "Static file not found",
		}))
	}

	// recover from panic
	app.Use(recover.New())

	// startup banner and info
	log.Info("SERVER:LOADING", "\n", banner.Console)
	if cfg.IsDev() {
		log.Info("SERVER:ENV", "Development mode ON")
	} else {
		log.Info("SERVER:ENV", "Production mode ON")
	}
	log.Info("SERVER:VERSION", build.Version())
	log.Info("SERVER:START", "Listening in", *cfg.Main.Server.IP+":"+*cfg.Main.Server.Port)

	// livereload
	if cfg.IsDev() {
		dev.StartLiveReload()
	}

	// server UP
	app.Listen(*cfg.Main.Server.IP + ":" + *cfg.Main.Server.Port)
}
