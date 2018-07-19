check-audit
==================

This utility updates checks audit process by consuming kafka messages from
`audit-events` topic.

### How to run the utility locally

Run:
- `go build`
- `./check-audit -kafka-brokers='<brokers>'`

Where `<brokers>` should be a comma-separated list like `localhost:9092` (i.e. `<ip-address>:<port>`).

### How to run the utility on an environment

- get the list of IP addresses for the kafka brokers on the target environment (e.g. `<ip1:port1>,<ip2:port2>`)
- ssh onto a target box that has access to these kafka brokers (e.g. an instance of dimension-extractor service)
  - on that target box, check the operating system and architecture with `uname -a`
- locally, cross-compile (using environment variables `GOOS` and `GOARCH`) a binary for running on the target box:
  - `GOOS=linux GOARCH=amd64 go build` (theses are typical values for AWS hosts)
- copy that cross-compiled binary onto the target box using one of:
  - `scp check-audit <username>@<target-ip-address>:`
  - `cd ../../dp-setup/ansible && scp -F ./ssh.cfg ../../dp-data-tools/check-audit/check-audit <username>@<target-ip-address>:`
- on the target box (via ssh), run `./check-audit -kafka-brokers='<ip1:port1>,<ip2:port2>'`
