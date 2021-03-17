// copy editions collection from datasets database

db = db.getSiblingDB('datasets')

db.editions.find().forEach(function(doc) {
    printjson(doc);
})
