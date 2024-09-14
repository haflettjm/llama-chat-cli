package ingestyaml

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
  Api API `yaml:"api"`
  ChatConfig CHATCONF `yaml:"chat_config"`
}

type API struct{
  Key string `yaml:"key"`
  ChatEndpoint string `yaml:"chat-endpoint"`
}

type CHATCONF struct {
  Org string `yaml:"organization"`
  Project string `yaml:"project"`
  Model string `yaml:"model"`
  Temp float64 `yaml:"temperature"`
  URole string `yaml:"user"`
}

type Entry interface{}

func UnmarshalYaml(data []byte) (Config, error) {
	unmarshalled := Config{}
	if err := yaml.Unmarshal(data, &unmarshalled); err != nil {
    return Config{}, err
	}

	return unmarshalled, nil
}

func Ingest(filepath string) (Config, error) {

  dat, err := os.ReadFile(filepath)

  if err != nil {
    return Config{}, err
  }
  unmarshal, err := UnmarshalYaml(dat)
  if err != nil {
    return Config{}, err
  }
  return unmarshal, nil
}
