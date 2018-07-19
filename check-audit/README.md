check-audit
==================

This utility updates checks audit process by consuming kafka messages from
`audit-events` topic.

### How to run the utility

- Run `go build`
- Run `./check-audit -kafka-brokers='<brokers>localhost:9092,localhost:9093<brokers>'`

The kafka brokers should look something like the following `localhost:9092` or
if running utility on environment `<ip-address>:<port>`.

To be able to run script against environment, one must find a list of ip addresses
for the kafka brokers for which the environment is using. You should then ssh onto
a box that has access to these kafka brokers (e.g. ip address for a running instance
of dimension extractor service), check the operating system with `uname -a` and
then set environment variables `GOOS` and `GOARCH` where you plan to build `check-audit`
binary (e.g. `export GOOS=linux GOARCH=amd64; go build`). Now copy binary onto box
with access to kafka brokers using scp command, something like the following:

`cd dp-setup/ansible; scp <relative-file-path>/dp-data-tools/check-audit/check-audit <username>@<ip-address>:. `

Then ssh onto box and run `./check-audit`
