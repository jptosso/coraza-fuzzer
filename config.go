package fuzzer

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Fuzzer struct {
		Rules           []string `yaml:"rules"`
		Transformations []string `yaml:"transformations"`
		Operators       []struct {
			Name string   `yaml:"name"`
			Args []string `yaml:"args"`
		} `yaml:"operators"`
		MinLength  int64  `yaml:"min_length"`
		MaxLength  int64  `yaml:"max_length"`
		Steps      string `yaml:"steps"`
		Iterations int64  `yaml:"iterations"`
	} `yaml:"fuzzer"`
}

func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := new(Config)
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
