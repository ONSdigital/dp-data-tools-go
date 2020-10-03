
if [ -z "$1" ]
  then
    echo "Please supply the mongo connection string as the first parameter, e.g mongodb://localhost:27017"
    exit 1
fi

mongo $1 <<EOF

 use datasets
 db.dimension.options.remove({})
 db.instances.remove({})
 db.editions.remove({})
 db.datasets.remove({})

 use imports
 db.imports.remove({})

 use filters
 db.filters.remove({})
 db.filterOutputs.remove({})

 use recipes
 db.recipes.remove({})

 var file = cat('./recipes.json');
 use recipes
 var recipes = JSON.parse(file);
 db.recipes.insert(recipes)

EOF

