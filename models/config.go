package models

import (
	"crypto/sha256"
	"io"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

var (
	configPath = "config/config.toml"
)

type duration time.Duration

// Config struct
type Config struct {
	ServerOpt ServerOpt `toml:"ServerOpt"`

	HashSum []byte
}

func (d *duration) UnmarshalText(text []byte) error {
	temp, err := time.ParseDuration(string(text))
	*d = duration(temp)
	return err
}

// ServerOpt struct
type ServerOpt struct {
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	IdleTimeout          time.Duration
	CacheCleanupInterval time.Duration
}

// LoadConfig from path
func LoadConfig(c *Config) {
	_, err := toml.DecodeFile(configPath, c)
	if err != nil {
		return
	}

	c.HashSum = GetHashSum()
}

// GetHashSum of config file
func GetHashSum() []byte {
	h := sha256.New()

	f, err := os.Open(configPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	if _, err = io.Copy(h, f); err != nil {
		return nil
	}

	return h.Sum(nil)
}
