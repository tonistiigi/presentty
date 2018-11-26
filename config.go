package main

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Config provides containerd configuration data for the server
type Config struct {
	Demos map[string]DemoConfig `toml:"demo"`
}

type DemoConfig struct {
	// ID string `toml`
	Build       string `toml:"build"`
	MainCommand string `toml:"main"`
	Command     string `toml:"cmd"`
}

func Load(r io.Reader) (Config, *toml.MetaData, error) {
	var c Config
	md, err := toml.DecodeReader(r, &c)
	if err != nil {
		return c, nil, errors.Wrap(err, "failed to parse config")
	}
	return c, &md, nil
}

func LoadFile(fp string) (Config, *toml.MetaData, error) {
	f, err := os.Open(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil, nil
		}
		return Config{}, nil, errors.Wrapf(err, "failed to load config from %s", fp)
	}
	defer f.Close()
	return Load(f)
}
