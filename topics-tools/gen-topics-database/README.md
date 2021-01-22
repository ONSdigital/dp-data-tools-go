# gen-topics-database
1. Crawl enough of ONS site to build collections for dp-topic-api mongo topic database

2. Crawl through whole ONS site and Audit all URI links

## read the instructions in the first 100 lines of comments in main.go

The notes will help you decide on how you want to run the app.

With the default flags the app will output two scripts that can be used to init the
topics and content collections in the topics database, thus:

mongo topics-init.js

mongo content-init.js

