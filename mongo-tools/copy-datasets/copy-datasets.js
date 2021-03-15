// copy datasets collection from datasets database

// run on macbook, thus:
//   have a shell that is logged into mongo on develop
//   in another shell run command:
//     mongo mongodb://root:< get secrets value for mongo and place here >@localhost:27017 copy-datasets.js >datasets.json

db = db.getSiblingDB('datasets')

db.datasets.find().forEach(function(doc) {
    printjson(doc);
})
