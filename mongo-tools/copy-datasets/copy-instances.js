// copy instances collection from datasets database

db = db.getSiblingDB('datasets')

db.instances.find().forEach(function(doc) {
    printjson(doc);
})
