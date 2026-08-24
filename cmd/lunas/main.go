// Command lunas runs the API server and, at deploy time, serves the built SPA
// from the same binary — one artifact, one port, no CORS.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "modernc.org/sqlite"

	assets "github.com/AlfinRy/lunas"
	"github.com/AlfinRy/lunas/internal/api"
	"github.com/AlfinRy/lunas/internal/config"
	db "github.com/AlfinRy/lunas/internal/db"
	"github.com/AlfinRy/lunas/internal/seed"
)


func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	conn, err := sql.Open("sqlite", cfg.DBPath+"?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, assets.SchemaSQL()); err != nil {
		log.Fatalf("apply schema: %v", err)
	}
	q := db.New(conn)
	if err := q.InitSettings(ctx); err != nil {
		log.Fatalf("init settings: %v", err)
	}
	if n, err := q.CountClients(ctx); err == nil && n == 0 {
		if err := seed.Run(ctx, q); err != nil {
			log.Fatalf("seed: %v", err)
		}
		log.Printf("seeded demo data")
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORS()) // dev: Vite proxy on 5173

	h := api.New(q, cfg)
	strict := api.NewStrictHandler(h, nil)
	api.RegisterHandlersWithBaseURL(e, strict, "/api")

	// Serve the SPA (from disk in dev, embedded at deploy).
	var webFS fs.FS = assets.Web()
	if _, err := os.Stat("web/dist/index.html"); err == nil {
		webFS = os.DirFS("web/dist")
	}
	e.GET("/*", staticHandler(webFS))
	e.GET("/", staticHandler(webFS))

	go func() {
		addr := ":" + cfg.Port
		log.Printf("Lunas listening on http://localhost%s (template_mode=%v)", addr, cfg.TemplateMode())
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("start: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = e.Shutdown(shutdownCtx)
}

func staticHandler(webFS fs.FS) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Param("*")
		if path == "" {
			path = "index.html"
		}
		if _, err := fs.Stat(webFS, path); err != nil {
			// SPA fallback: unknown paths get index.html.
			path = "index.html"
		}
		c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
		return fsFile(c, webFS, path)
	}
}

func fsFile(c echo.Context, fsys fs.FS, name string) error {
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Not found.")
	}
	c.Response().Header().Set(echo.HeaderContentType, contentTypeFor(name))
	if _, err := c.Response().Write(data); err != nil {
		return err
	}
	return nil
}

func contentTypeFor(name string) string {
	switch {
	case strings.HasSuffix(name, ".html"):
		return "text/html; charset=utf-8"
	case strings.HasSuffix(name, ".js"):
		return "text/javascript; charset=utf-8"
	case strings.HasSuffix(name, ".css"):
		return "text/css; charset=utf-8"
	case strings.HasSuffix(name, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(name, ".woff2"):
		return "font/woff2"
	case strings.HasSuffix(name, ".json"):
		return "application/json"
	default:
		return fmt.Sprintf("application/octet-stream")
	}
}
