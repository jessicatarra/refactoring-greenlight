package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	Version string
	port    string
	env     string
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
		Sender   string
	}
	Cors struct {
		TrustedOrigins []string
	}
	Jwt struct {
		Secret string
	}
}

func Init() (cfg Config, err error) {
	intPort, _ := strconv.Atoi(port)
	flag.IntVar(&cfg.Port, "port", intPort, "API server port")
	flag.StringVar(&cfg.Env, "env", env, "Environment (development|staging|production)")
	flag.StringVar(&cfg.DB.Dsn, "db-dsn", os.Getenv(
		"DATABASE_URL"), "PostgreSQL DSN")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	flag.StringVar(&cfg.Smtp.Host, "smtp-host", os.Getenv("SMTP_HOST"), "SMTP host")
	flag.IntVar(&cfg.Smtp.Port, "smtp-port", smtpPort, "SMTP port")
	flag.StringVar(&cfg.Smtp.Username, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.Smtp.Password, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.Smtp.Sender, "smtp-sender", os.Getenv("SMTP_SENDER"), "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.Cors.TrustedOrigins = strings.Fields(val)
		return nil
	})

	flag.StringVar(&cfg.Jwt.Secret, "jwt-secret", os.Getenv("JWT_SECRET"), "JWT secret")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", Version)
	}

	return cfg, nil
}
