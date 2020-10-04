#!/usr/bin/env bash

if [[ -z "$1" ]]
  then
    echo "Please supply the mongo connection string as the first parameter, e.g mongodb://localhost:27017"
    exit 1
fi

mongo $1 <<EOF

 use datasets
 db.dimension.options_neo4j.renameCollection("dimension.options", true)
 db.instances_neo4j.renameCollection("instances", true)
 db.editions_neo4j.renameCollection("editions", true)
 db.datasets_neo4j.renameCollection("datasets", true)

 use imports
 db.imports_neo4j.renameCollection("imports", true)

 use filters
 db.filters_neo4j.renameCollection("filters", true)
 db.filterOutputs_neo4j.renameCollection("filterOutputs", true)

 use recipes
 db.recipes_neo4j.renameCollection("recipes", true)

EOF

