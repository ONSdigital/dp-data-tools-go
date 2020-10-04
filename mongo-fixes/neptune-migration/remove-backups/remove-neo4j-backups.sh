#!/usr/bin/env bash

if [[ -z "$1" ]]
  then
    echo "Please supply the mongo connection string as the first parameter, e.g mongodb://localhost:27017"
    exit 1
fi

mongo $1 <<EOF

 use datasets
 db.dimension.options_neo4j.drop()
 db.instances_neo4j.drop()
 db.editions_neo4j.drop()
 db.datasets_neo4j.drop()

 use imports
 db.imports_neo4j.drop()

 use filters
 db.filters_neo4j.drop()
 db.filterOutputs_neo4j.drop()

 use recipes
 db.recipes_neo4j.drop()

EOF

