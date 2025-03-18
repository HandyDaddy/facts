package config

import "github.com/BurntSushi/toml"

type (
	Config struct {
		App        App        `toml:"Application"`
		HttpClient HttpClient `toml:"HttpClient"`
	}

	App struct {
		Name          string `toml:"Name"`
		Version       string `toml:"Version"`
		MaxBufferSize uint   `toml:"MaxBufferSize"`
	}

	HttpClient struct {
		Addr  string `toml:"Addr"`
		Token string `env:"BEARERTOKEN"`
	}
)

func Parse(path string) (*Config, error) {
	var conf Config
	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
