package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

type Args struct {
	BindAddr string
	Logger   *slog.Logger
}

type config struct{}

type Server struct {
	logger *slog.Logger

	http  *http.Client
	httpd *http.Server

	g *gin.Engine

	config config
}

func New(args Args) (*Server, error) {
	if args.BindAddr == "" {
		return nil, errors.New("bind-addr must be set")
	}

	if args.Logger == nil {
		args.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	}

	g := gin.New()
	g.Use(sloggin.New(args.Logger.WithGroup("http")))
	g.Use(gin.Recovery())
	g.Use(cors.Default())

	httpd := &http.Server{
		Addr: args.BindAddr,
		Handler: g.Handler(),

		// FIXME: Should move these to top-level consts instead of magic vals
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 5 * time.Minute,
	}

	srv := &Server{
		logger: args.Logger,
		http:   http.DefaultClient,
		httpd:  httpd,
		g:      g,
		config: config{},
	}

	return srv, nil
}

func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("starting nook")

	go func() {
		if err := s.httpd.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()

	return nil
}
