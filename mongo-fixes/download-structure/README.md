# edition-doc-structure

This script is to update edition resources on cmd due to structural changes.

## How to run service

- Run `go build`
- Run `./edition-doc-structure -mongo-url=<url> -download-service-url=<url>`

The mongodb url should look something like the following `localhost:27017` or
`127.0.0.1:27017`. If a username and password are needed follow this structure
`<username>:<password>@<host>:<port>`

The download service url should be correct for the environment you are running on.
E.g `http://localhost:23600` in dev
