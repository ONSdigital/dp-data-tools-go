# generate-downloads

`generate-downloads` is a utility to queue a message in Kafka which should trigger the full downloads to be rebuilt for a given published dataset, when they are missing.

## Configuration

You will need different values (and possibly more env vars) than shown in the examples and/or `Makefile`, below, see [config](./main.go) for the defaults.

## Queueing the message

Queue the message 'locally' (using an ssh tunnel) or remotely (run the binary inside the environment).

### Run locally

Either run locally:

- ssh tunnel to kafka (if targeting a cloud environment)
- edit the [Makefile](./Makefile) with the appropriate config (env vars), see the make target `all`
- to queue the message
  - `$ make`

### Run in an environment

Alternatively, to run the binary inside an environment:

- build the binary (ensure the correct `GOOS` setting in the `Makefile` for the remote host)
  - `$ make build`
- push that binary to the env (this example is for `develop`)
  - `$ dp scp develop publishing 2 generate-downloads .`
- login to that host
  - `$ dp ssh develop publishing 2`
- on the remote host :warning:, queue the message
  - `$ INSTANCE_ID="xb1ae3d1-913e-43e0-b4c9-2c741744f12" DATASET_ID="weekly-deaths-local-authoritay" VERSION="2" ./generate-downloads`
