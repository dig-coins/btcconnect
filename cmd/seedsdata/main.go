package main

import (
	"log"
	"os"

	"github.com/dig-coins/btcconnect/internal/edfile"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Dest   string   `yaml:"dest"`
	SecKey string   `yaml:"secKey"`
	Seeds  []string `yaml:"seeds"`
}

func main() {
	d, err := os.ReadFile("seeds-data.yaml")
	if err != nil {
		log.Panicln(err)
	}

	var config Config

	err = yaml.Unmarshal(d, &config)
	if err != nil {
		log.Panicln(err)
	}

	if len(config.Seeds) == 0 {
		log.Panicln("empty seeds")
	}

	d, err = yaml.Marshal(config.Seeds)
	if err != nil {
		log.Panicln(err)
	}

	err = edfile.WriteSecFile(config.Dest, config.SecKey, d)
	if err != nil {
		log.Panicln(err)
	}
}
