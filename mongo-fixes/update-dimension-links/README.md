update-dimension-links
==================

This script is to update the dimensions URLs in each instances to be unversioned (to not contain /v1)
Maybe also rename the URL domain to localhost??? 

### How to run service
* Run `go build`
* Run `./update-dimension-links -mongo-url=<url>`

The mongodb url should look something like the following `localhost:27017` or
`127.0.0.1:27017`. If a username and password are needed follow this structure
`<username>:<password>@<host>:<port>`