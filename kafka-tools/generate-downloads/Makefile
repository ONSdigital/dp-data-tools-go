# see main.go for all env var options/defaults

.PHONY: all
all:
	KAFKA_ADDR=localhost:9092 INSTANCE_ID="xb1ae3d1-913e-43e0-b4c9-2c741744f12" DATASET_ID="weekly-deaths-local-authority" VERSION="2" go run -race main.go

# you may want to build a binary to run on AWS:
.PHONY: build
build:
	GOOS=linux go build -o generate-downloads main.go
