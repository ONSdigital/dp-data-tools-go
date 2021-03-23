package main

import (
	"context"
	"flag"
	"github.com/ONSdigital/log.go/log"
)

var (
	zebedeeURL         string
	mapperPath string
)

func main() {
	setupFlags()

	ctx := context.Background()

	if validateMandatoryParams(ctx) {
		return
	}

	log.Event(ctx, "successfully updated all documents.", log.INFO)
}

func validateMandatoryParams(ctx context.Context) bool {
	if zebedeeURL == "" {
		log.Event(ctx, "missing zebedeeURL flag", log.ERROR)
		return true
	}

	if mapperPath == "" {
		log.Event(ctx, "missing mapper-path flag", log.ERROR)
		return true
	}
	return false
}

func setupFlags() {
	flag.StringVar(&zebedeeURL, "zebedee-url", zebedeeURL, "Zebedee API URL")
	flag.StringVar(&mapperPath, "mapper-path", mapperPath, "Path to the mapper")
	flag.Parse()
}
