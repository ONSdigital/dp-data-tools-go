filter-doc-version-identifier
==================

This script is to update filter blueprint and output resources on cmd due to 
identifying the version in which filter is being applied to has changed from
using the `instance_id` to a dataset object containing `id`, `edition` and 
`version`.

### How to run service
* Run `go build`
* Run `./filter-doc-version-identifier -mongo-url=<url>`

The url should look something like the following `localhost:27017` or
`127.0.0.1:27017`. If a username and password are needed follow this structure
`<username>:<password>@<host>:<port>`
