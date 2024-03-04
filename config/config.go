package config

import (
	"os"
	"path"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type CommitInfo struct {
	FullSHA1 string
	Time     string
}

type DataBase struct {
	Host     string `envconfig:"PG_HOST" required:"true"`
	Port     int    `envconfig:"PG_PORT" required:"true"`
	User     string `envconfig:"PG_USER" required:"true"`
	Password string `envconfig:"PG_PASSWORD" required:"true"`
	Name     string `envconfig:"PG_DATABASE" required:"true"`
}

type Config struct {
	Port       string `envconfig:"SERVER_PORT" default:"1300"`
	Domain     string `envconfig:"DOMAIN"`
	AppCode    string `envconfig:"APP_Code"`
	AppName    string `envconfig:"APP_NAME"`
	AppVersion string `envconfig:"API_VERSION"`
}

type Settings struct {
	Config
	DataBase
	RootDIR     string
	Env         string
	DSN         string
	UploadsPath string
	JwtExpiry   time.Duration
	JwtSecret   []byte
}

var (
	JwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
	DOMAIN    = os.Getenv("DOMAIN")
	RootDIR   = os.Getenv("ROOT_DIR")
	ENV       = os.Getenv("ENV")
	DSN       = os.Getenv("CONNECTION_STRING")
	Port      = os.Getenv("PORT")
	AppCode   = os.Getenv("APP_CODE")
	AppName   = os.Getenv("APP_NAME")
	JwtExpiry = getJWTExpiryHours()
)

func GetSettings() *Settings {
	return &Settings{
		Env:         ENV,
		DSN:         DSN,
		UploadsPath: GetUploadsPath(""),
		JwtExpiry:   getJWTExpiryHours(),
	}
}

func GetUploadsPath(dir string) string {
	return path.Join(RootDIR, "public", "uploads", dir)
}

func GetRootPath(dir string) string {
	return path.Join(RootDIR, dir)
}

func getJWTExpiryHours() time.Duration {
	expEnv := os.Getenv("JWT_EXPIRY_IN_HOURS")
	if expInt, err := strconv.ParseFloat(expEnv, 64); err != nil {
		return 1 * time.Hour
	} else {
		return time.Duration(expInt) * time.Hour
	}
}

func TimeNow() time.Time {
	return time.Now().UTC()
}

func NewParsedConfig() (*Settings, error) {
	_ = godotenv.Load(".env")
	cfg := Settings{}
	err := envconfig.Process("", &cfg)

	cfg.JwtExpiry = getJWTExpiryHours()
	cfg.JwtSecret = JwtSecret
	cfg.Env = ENV

	return &cfg, err
}
