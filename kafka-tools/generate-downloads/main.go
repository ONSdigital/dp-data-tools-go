package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ONSdigital/dp-dataset-api/download"
	adapter "github.com/ONSdigital/dp-dataset-api/kafka"
	"github.com/ONSdigital/dp-dataset-api/schema"
	kafka "github.com/ONSdigital/dp-kafka"

	"github.com/kelseyhightower/envconfig"
)

type configuration struct {
	DatasetID  string `envconfig:"DATASET_ID"`
	InstanceID string `envconfig:"INSTANCE_ID"`
	Edition    string `envconfig:"EDITION"`
	Version    string `envconfig:"VERSION"`

	Timeout                time.Duration `envconfig:"TIMEOUT"`
	KafkaAddr              []string      `envconfig:"KAFKA_ADDR"`
	GenerateDownloadsTopic string        `envconfig:"GENERATE_DOWNLOADS_TOPIC"`
}

var defaultCfg = configuration{
	KafkaAddr:              []string{"kafka-1:9092", "kafka-2:9092", "kafka-3:9092"},
	GenerateDownloadsTopic: "filter-job-submitted",
	Timeout:                30 * time.Second,
	Edition:                "time-series",
}

func main() {

	cfg := getConfig()
	fmt.Printf("Config: %+v\n", cfg)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	pChannels := kafka.CreateProducerChannels()

	// kafka may block, so do work in goroutine (cancel context when done)
	go func() {
		defer cancel()

		prod, err := kafka.NewProducer(ctx, cfg.KafkaAddr, cfg.GenerateDownloadsTopic, 0, pChannels)
		if err != nil {
			fmt.Fprintf(os.Stderr, "NewProducer failed: %v\n", err)
			return
		}
		defer func() {
			fmt.Println("closing producer...")
			prod.Close(ctx)
			fmt.Println("producer closed")
		}()
		downloadGenerator := &download.Generator{
			Producer:   adapter.NewProducerAdapter(prod),
			Marshaller: schema.GenerateDownloadsEvent,
		}

		fmt.Println("waiting for kafka initialisation...")
		select {
		case <-ctx.Done():
			return
		case <-pChannels.Init:
		}
		fmt.Println("kafka initialised")
		tmr := time.NewTimer(1 * time.Second) // settle time seems necessary
		select {
		case <-ctx.Done():
			tmr.Stop()
			return
		case <-tmr.C:
		}

		fmt.Println("message send...")
		if err = downloadGenerator.Generate(ctx, cfg.DatasetID, cfg.InstanceID, cfg.Edition, cfg.Version); err != nil {
			fmt.Fprintf(os.Stderr, "Generate failed: %v\n", err)
			return
		}

		fmt.Println("message sent, pausing...")
		tmr2 := time.NewTimer(15 * time.Second) // pause before closing, necessary to allow message to depart
		select {
		case <-ctx.Done():
			tmr2.Stop()
			return
		case <-tmr2.C:
		}
	}()

	// wait for kafka work to complete (or timeout, or error)
	select {
	case err := <-pChannels.Errors:
		if err != nil {
			fmt.Fprintf(os.Stderr, "producer error: %s\n", err)
		}
		cancel()
	case <-ctx.Done():
		if err := ctx.Err(); err != nil && err != context.Canceled {
			panic(err)
		}
	}
	fmt.Println("done")
}

func getConfig() (cfg *configuration) {
	cfg = &configuration{}
	*cfg = defaultCfg

	if err := envconfig.Process("", cfg); err != nil {
		panic(err)
	}
	if cfg.DatasetID == "" {
		panic("no dataset id")
	}
	if cfg.InstanceID == "" {
		panic("no instance id")
	}
	if cfg.Edition == "" {
		panic("no edition")
	}
	if cfg.Version == "" {
		panic("no version")
	}
	return
}
