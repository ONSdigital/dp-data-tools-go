fix-edition-docs
==================

This small application is to update edition resources on cmd due to structural
changes.

### How to run service
* Run `go build`
* Run `./fix-edition-docs -mongo-url=<url>`

The url should look something like the following `localhost:27017` or
`127.0.0.1:27017`. If a username and password are needed follow this structure
`<username>:<password>@<host>:<port>`
