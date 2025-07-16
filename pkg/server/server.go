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
	Did      string
	Hostname string
	Logger   *slog.Logger
}

type config struct {
	Did      string
	Hostname string
}

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

	if args.Did == "" {
		return nil, errors.New("did must be set")
	}

	if args.Hostname == "" {
		return nil, errors.New("hostname must be set")
	}

	if args.Logger == nil {
		args.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	}

	g := gin.New()
	g.Use(sloggin.New(args.Logger.WithGroup("http")))
	g.Use(gin.Recovery())
	g.Use(cors.Default())

	httpd := &http.Server{
		Addr:    args.BindAddr,
		Handler: g.Handler(),

		// FIXME: Should move these to top-level consts instead of magic vals
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  5 * time.Minute,
	}

	srv := &Server{
		logger: args.Logger,
		http:   http.DefaultClient,
		httpd:  httpd,
		g:      g,
		config: config{
			Did:      args.Did,
			Hostname: args.Hostname,
		},
	}

	return srv, nil
}

func (s *Server) setupRoutes() {
	// Static Routes
	s.g.GET("/", s.handleRoot)
	s.g.GET("/robots.txt", s.handleRobotsTxt)

	// Metadata Routes
	s.g.GET("/.well-known/did.json", s.handleWellKnown)
	s.g.GET("/.well-known/oauth-protected-resource", s.handleOauthProtectedResource)
}

func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("mounting routes...")
	s.setupRoutes()

	s.logger.Info("starting nook")

	go func() {
		if err := s.httpd.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()

	return nil
}
