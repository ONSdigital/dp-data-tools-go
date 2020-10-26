update-dimension-links
==================

This program updates dimension URLs in the `instances` collection -  removing any versioning (i.e. `/v1/` substring)

### How to run service
* Run `go build`
* Run `./update-dimension-links -mongo-url=<url>`

The mongodb url should look something like the following `localhost:27017` or
`127.0.0.1:27017`. If a username and password are needed follow this structure
`<username>:<password>@<host>:<port>`
