#!/usr/bin/env bash

if [[ -z "$1" ]]
  then
    echo "Please supply the mongo connection string as the first parameter, e.g mongodb://localhost:27017"
    exit 1
fi

mongo $1 <<EOF

 use datasets
 db.dimension.options.renameCollection("dimension.options_neo4j", false)
 db.createCollection("dimension.options")
 db.instances.renameCollection("instances_neo4j", false)
 db.createCollection("instances")
 db.editions.renameCollection("editions_neo4j", false)
 db.createCollection("editions")
 db.datasets.renameCollection("datasets_neo4j", false)
 db.createCollection("datasets")

 use imports
 db.imports.renameCollection("imports_neo4j", false)
 db.createCollection("imports")

 use filters
 db.filters.renameCollection("filters_neo4j", false)
 db.createCollection("filters")
 db.filterOutputs.renameCollection("filterOutputs_neo4j", false)
 db.createCollection("filterOutputs")

 use recipes
 db.recipes.renameCollection("recipes_neo4j", false)
 db.createCollection("recipes")

 var file = cat('./recipes.json');
 use recipes
 var recipes = JSON.parse(file);
 db.recipes.insert(recipes)

EOF

