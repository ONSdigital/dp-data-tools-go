## Neptune migration

** Please note - once this script runs, the MongoDB indexes will need recreating by running the MongoDB Ansible in dp-setup **

This script contains all the Mongo database updates for the migration from Neo4j to Neptune. 

All existing collection in Mongo DB will be renamed for backup with a `_neo4j` suffix
 - instances
 - editions
 - dimensions
 - datasets
 - recipes
 - filters
 
The original collections will be recreated, and updated recipes get imported from `recipes.json`

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

### restore the Neo4j collections

If you need to restore the Neo4j backup collections to their original names you can use the `restore-neo4j/restore-neo4j.sh` script. 
The script takes the MongoDB connection string as a parameter, in the same way as the Neptune migration script detailed above.

```
./restore-neo4j/restore-neo4j.sh mongodb://localhost:27017
```

### remove the backup collections

Once the migration is complete, and it's safe to remove the backups they can be removed using the `remove-backups/remove-neo4j-backups.sh` script.

```
./remove-backups/remove-neo4j-backups.sh mongodb://localhost:27017
```