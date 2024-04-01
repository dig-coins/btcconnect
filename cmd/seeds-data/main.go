package main

import (
	"log"
	"os"

	"github.com/dig-coins/btcconnect/internal/edfile"
	"gopkg.in/yaml.v3"
)

type Item struct {
	Dest   string   `yaml:"dest"`
	SecKey string   `yaml:"secKey"`
	Seeds  []string `yaml:"seeds"`
}

type Config struct {
	Items []Item `yaml:"items"`
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

	for _, item := range config.Items {
		if len(item.Seeds) == 0 {
			log.Panicln("empty seeds")
		}

		d, err = yaml.Marshal(item.Seeds)
		if err != nil {
			log.Panicln(err)
		}

		err = edfile.WriteSecFile(item.Dest, item.SecKey, d)
		if err != nil {
			log.Panicln(err)
		}
	}
}
