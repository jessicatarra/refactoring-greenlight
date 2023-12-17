package config

import (
	"flag"
	"fmt"
	"strings"
)

var (
	Version string
)

type Config struct {
	Port int
	Env  string
	DB   struct {
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Smtp struct {
		Host     string
		Port     int
		Username string
		Password string
		From     string
	}
	Cors struct {
		TrustedOrigins []string
	}
	Jwt struct {
		Secret string
	}
	Auth struct {
		HttpBaseURL    string
		GrpcBaseURL    string
		GrpcServerPort int
		HttpPort       int
	}
}

func Init() (cfg Config, err error) {
	flag.IntVar(&cfg.Port, "port", 8080, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.DB.Dsn, "db-dsn", "postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable", "PostgreSQL DSN")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.StringVar(&cfg.Smtp.Host, "smtp-host", "example.smtp.host", "SMTP host")
	flag.IntVar(&cfg.Smtp.Port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.Smtp.Username, "smtp-username", "example_username", "SMTP username")
	flag.StringVar(&cfg.Smtp.Password, "smtp-password", "example_password", "SMTP password")
	flag.StringVar(&cfg.Smtp.From, "smtp-sender", "Example Name <no-reply@example.org>", "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.Cors.TrustedOrigins = strings.Fields(val)
		return nil
	})

	flag.StringVar(&cfg.Jwt.Secret, "jwt-secret", "56vphh6sheco5sbtfkxwesy3wx7fpiip", "JWT secret")

	flag.StringVar(&cfg.Auth.HttpBaseURL, "base-url", "http://localhost:8082", "base URL for the application")
	flag.StringVar(&cfg.Auth.GrpcBaseURL, "auth-grpc-client-base-url", "localhost:50051", "GRPC client")

	flag.IntVar(&cfg.Auth.HttpPort, "auth-http-port", 8082, "port to listen on for HTTP requests for auth module")

	flag.IntVar(&cfg.Auth.GrpcServerPort, "auth-grpc-port", 50051, "port to listen on for GRPC methods for auth module")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", Version)
	}

	return cfg, nil
}
