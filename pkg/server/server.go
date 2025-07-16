package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/bluesky-social/indigo/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/hbjydev/nook/internal/db"
	sloggin "github.com/samber/slog-gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Args struct {
	BindAddr string
	Did      string
	Hostname string
	DbDsn    string
	Version  string
	Logger   *slog.Logger
}

type config struct {
	Did      string
	Hostname string
	Version  string
}

type Server struct {
	logger *slog.Logger

	http  *http.Client
	httpd *http.Server

	g  *gin.Engine
	db *db.DB

	config config
}

type CustomValidator struct {
	validator *validator.Validate
}

type ValidationError struct {
	error
	Field string
	Tag   string
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		var validateErrors validator.ValidationErrors
		if errors.As(err, &validateErrors) && len(validateErrors) > 0 {
			first := validateErrors[0]
			return ValidationError{
				error: err,
				Field: first.Field(),
				Tag:   first.Tag(),
			}
		}

		return err
	}

	return nil
}

func New(args Args) (*Server, error) {
	if args.BindAddr == "" {
		return nil, errors.New("bind-addr must be set")
	}

	if _, err := syntax.ParseDID(args.Did); err != nil {
		return nil, errors.New("did must be set and valid")
	}

	if args.Hostname == "" {
		return nil, errors.New("hostname must be set")
	}

	if args.DbDsn == "" {
		return nil, errors.New("db dsn must be set")
	}

	if args.Version == "" {
		return nil, errors.New("version must be set")
	}

	if args.Logger == nil {
		args.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	}

	g := gin.New()
	g.Use(sloggin.New(args.Logger.WithGroup("http")))
	g.Use(gin.Recovery())
	g.Use(cors.Default())

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("atproto-handle", func(fl validator.FieldLevel) bool {
			if _, err := syntax.ParseHandle(fl.Field().String()); err != nil {
				return false
			}
			return true
		})
		v.RegisterValidation("atproto-did", func(fl validator.FieldLevel) bool {
			if _, err := syntax.ParseDID(fl.Field().String()); err != nil {
				return false
			}
			return true
		})
		v.RegisterValidation("atproto-rkey", func(fl validator.FieldLevel) bool {
			if _, err := syntax.ParseRecordKey(fl.Field().String()); err != nil {
				return false
			}
			return true
		})
		v.RegisterValidation("atproto-nsid", func(fl validator.FieldLevel) bool {
			if _, err := syntax.ParseNSID(fl.Field().String()); err != nil {
				return false
			}
			return true
		})
	}

	httpd := &http.Server{
		Addr:    args.BindAddr,
		Handler: g.Handler(),

		// FIXME: Should move these to top-level consts instead of magic vals
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  5 * time.Minute,
	}

	gdb, err := gorm.Open(sqlite.New(sqlite.Config{
		DriverName: "libsql",
		DSN:        args.DbDsn,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	dbw := db.New(gdb)

	srv := &Server{
		logger: args.Logger,
		http:   util.RobustHTTPClient(),
		httpd:  httpd,
		g:      g,
		db:     dbw,
		config: config{
			Did:      args.Did,
			Hostname: args.Hostname,
			Version:  args.Version,
		},
	}

	return srv, nil
}

func (s *Server) setupRoutes() {
	// General non-interactive stuff
	s.g.GET("/", s.handleRoot)
	s.g.GET("/robots.txt", s.handleRobotsTxt)
	s.g.GET("/xrpc/_health", s.handleHealth)
	s.g.GET("/.well-known/did.json", s.handleWellKnown)
	s.g.GET("/.well-known/oauth-protected-resource", s.handleOauthProtectedResource)

	// Public routes
	s.g.GET("/xrpc/com.atproto.identity.resolveHandle", s.handleIdentityResolveHandle)
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
