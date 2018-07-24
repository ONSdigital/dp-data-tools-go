package main

import (
	"errors"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/kafka"
	"github.com/ONSdigital/go-ns/log"
)

// AuditEvent represents the structure of the audit message received
type AuditEvent struct {
	Created         string            `avro:"created"`
	Service         string            `avro:"service"`
	RequestID       string            `avro:"request_id"`
	User            string            `avro:"user"`
	AttemptedAction string            `avro:"attempted_action"`
	ActionResult    string            `avro:"action_result"`
	Params          map[string]string `avro:"params"`
}

// Action represents stats of the consumed actions
type Action struct {
	Attempted    int
	Successful   int
	Unsuccessful int
	Total        int
}

var kafkaBrokers string

func main() {
	flag.StringVar(&kafkaBrokers, "kafka-brokers", kafkaBrokers, "kafka broker addresses, comma separated")
	flag.Parse()

	if kafkaBrokers == "" {
		log.Error(errors.New("missing kafka brokers, must be comma separated"), log.Data{"kafka_brokers": kafkaBrokers})
		return
	}

	kafkaBrokerList := strings.Split(kafkaBrokers, ",")

	syncConsumerGroup, err := kafka.NewSyncConsumer(kafkaBrokerList, "audit-events", "check-audit", kafka.OffsetOldest)
	if err != nil {
		log.ErrorC("could not obtain consumer", err, log.Data{"list_of_kafka_brokers": kafkaBrokerList})
		os.Exit(1)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	actions := make(map[string]Action)
	for {
		select {
		case message := <-syncConsumerGroup.Incoming():
			event, err := readMessage(message.GetData())
			if err != nil {
				log.Error(err, log.Data{"schema": "failed to unmarshal event"})
				break
			}

			log.Info("Received message", log.Data{"audit_event": event})

			addResult(actions, event)

			syncConsumerGroup.CommitAndRelease(message)
		case <-signals:
			log.Info("Audit stats", log.Data{"audit": actions})
			os.Exit(0)
		}
	}
}

func readMessage(eventValue []byte) (*AuditEvent, error) {
	var e AuditEvent

	if err := audit.EventSchema.Unmarshal(eventValue, &e); err != nil {
		return nil, err
	}

	return &e, nil
}

func addResult(actions map[string]Action, event *AuditEvent) {
	var successful, attempted, unsuccessful int

	switch event.ActionResult {
	case "successful":
		successful = 1
	case "unsuccessful":
		unsuccessful = 1
	case "attempted":
		attempted = 1
	}

	action, ok := actions[event.AttemptedAction]
	if !ok {
		actions[event.AttemptedAction] = Action{
			Attempted:    attempted,
			Successful:   successful,
			Unsuccessful: unsuccessful,
			Total:        1,
		}
	} else {
		if attempted == 1 {
			action.Attempted++
		} else {
			if successful == 1 {
				action.Successful++
			} else {
				action.Unsuccessful++
			}
		}
		action.Total++
		actions[event.AttemptedAction] = action
	}
}
