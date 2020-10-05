# New check-audit

This utility consumes messages from the `audit` kafka topic and builds a
list of accessed paths and their respective results (as `successful` and `unsuccessful`
counters) stored in-memory. The results (list of actions) are logged once the
utility has received a signal to terminate.

## How to run the utility locally

Run:
- `go build`
- `HUMAN_LOG=1 ./check-audit-new -kafka-brokers='<brokers>'`

Where `<brokers>` should be a comma-separated list like `localhost:9092` (i.e. `<ip-address>:<port>`).

## How to run the utility on an environment

- get the list of IP addresses for the kafka brokers on the target environment (e.g. `<ip1:port1>,<ip2:port2>`)
- ssh onto a target box that has access to these kafka brokers (e.g. an instance of dimension-extractor service)
  - on that target box, check the operating system and architecture with `$ uname -a`
- locally, cross-compile (using environment variables `GOOS` and `GOARCH`) a binary for running on the target box:
  - `$ GOOS=linux GOARCH=amd64 go build` (values typical for AWS hosts)
- copy that cross-compiled binary onto the target box using (replace `publishing 2` with your target box):
  - `$ dp scp develop publishing 2 check-audit-new .`
- login to that box:
  - `$ dp ssh develop publishing 2`
  - on the target box (via ssh), run
  - `$ ./check-audit-new -kafka-brokers='<ip1:port1>,<ip2:port2>'`

### Note about HUMAN_LOG and special characters encoding

If you run check-audit-new without HUMAN_LOG, then any special character will be presented as a unicode escaped ASCII hexadecimal code.
For example, if an audit event contains `query_param` with value `k1=v1&k2=v1`, the value will be displayed as `k1=v1\u0026k2=v2` unless the tool is run with HUMAN_LOG=1
