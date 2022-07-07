package main

import (
	"log"
)

var cfg Config

func main() {

	err := loadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	entryLog()
	Start()
}

func Start() {
	for _, v := range cfg.Jobs {
		v.start()
	}
}
