// copy ALL datasets.options collection from datasets database

db = db.getSiblingDB('datasets')

db.dimension.options.find().forEach(function(doc) {
    printjson(doc);
})
