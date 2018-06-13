# filter

This script is to update filter output docs to use the flattened event structure

## How to run service

- Run `go build`
- Run `./filter -mongo-url=<url>`

The mongodb url should look something like the following `localhost:27017` or
`127.0.0.1:27017`. If a username and password are needed follow this structure
`<username>:<password>@<host>:<port>`

