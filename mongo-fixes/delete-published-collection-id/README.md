delete-published-collection-id
==================

This utility updates `dataset` resources on CMD
due to it no longer storing `collection_id` when `{state:"published"}`.

### How to run the utility

Run
```
mongo <mongo_url> <options> delete-published-collection-id.js
```

The `<mongo_url>` part should look like:
- `localhost:27017/datasets` or `127.0.0.1:27017/datasets`
  - if authentication is needed, use:
    `mongodb://<username>:<password>@<host>:<port>/datasets?authSource=admin`
    (use single-quotes for protection from your shell)
- in the above, `/datasets` indicates the database to be modified

Example of the (optional) `<options>` part:

- `--eval 'cfg={verbose:true}'` (e.g. use for debugging)
- `cfg` defaults to: `{verbose:true, ids:true, update: true}`
- if you specify `cfg`, all missing options default to `false`

Full example (e.g. for capturing IDs and not wiping them):

```
mongo localhost:27017/datasets --eval 'cfg={verbose:true, ids:true}' delete-published-collection-id.js
```
