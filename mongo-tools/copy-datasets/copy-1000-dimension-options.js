// copy out first 1000 documents from datasets.options collection from datasets database

db = db.getSiblingDB('datasets')

var counter = 0
db.dimension.options.find().forEach(function(doc) {
    if (counter < 1000) {
        printjson(doc);
    }
    counter++;
})
