# dp-data-tools

 See below for various tools to update data.

## Current tools/scripts

### mongodb related

* [Edition document restructure](./mongo-fixes/edition-doc-structure)
* [Filter blueprint and output documents include new dataset object](./mongo-fixes/filter-doc-version-identifier)
* [Instance/version documents include new downloads structure](./mongo-fixes/download-structure/dataset)
* [Filter output documents include new downloads structure](./mongo-fixes/download-structure/filter)
* [Remove collection_id from published datasets](./mongo-fixes/delete-published-collection-id)

### kafka related

* [Check audit messages have been added to kafka](./kafka-tools/check-audit)
* [Queue a kafka message to rebuild full downloads for a dataset](./kafka-tools/generate-downloads)
