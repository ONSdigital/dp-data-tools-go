# dataset

This script is to link all old timeseries CDIDs to New timeseries CDIDs for PST

## How to run service

- Run `go build`
- Run 
```
./update-cdids -zebedee-url=<url> \
 -environment=<env-url> -username=<> -password=<> -filepath=<xlsx-file-path> -sheetname=<sheet-name-to-use> \
-limit=<count-of-rows-to-process>
```

The zebedee url should be correct for the environment you are running on.
E.g `http://localhost:8082`


The environment url should be correct for the environment you are running on.
E.g `https://publishing.develop.onsdigital.co.uk`
