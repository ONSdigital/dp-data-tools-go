## Neptune migration

This script contains all the Mongo database updates for the migration from Neo4j to Neptune. All existing data in Mongo DB is to be removed:
 - instances
 - editions
 - dimensions
 - datasets
 - recipes
 - filters
 
Once the existing data has been removed, the updated recipes get imported from `recipes.json`

### How to run the utility

Run
```
./neptune-migration.sh <mongo_url> 
```

The `<mongo_url>` part should look like:
- `mongodb://localhost:27017`
  - if authentication is needed, use:
    `mongodb://<username>:<password>@<host>:<port>?authSource=admin`
    (use single-quotes for protection from your shell)

Full example 

```
./neptune-migration.sh mongodb://localhost:27017
```
