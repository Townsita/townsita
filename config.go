package townsita

import (
	"encoding/json"
	"log"
	"os"
)

const defaultConfigFile = "./config/townsita.json"

type Config struct{}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(args []string) error {
	fileName := defaultConfigFile
	if len(args) > 1 {
		fileName = args[1]
	}
	log.Printf("Loading config from %s", fileName)
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(c); err != nil {
		return err
	}
	return nil
}

func (c *Config) templatePath(templateName string) string {
	return "./templates/" + templateName
}
