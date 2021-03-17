Copying **`datasets`** database from mongodb on develop to mongdb on your MacBook to enable dp-datasets-api to use local mongodb (thus saving you the hassle of running the import process locally to get the local stack running).

Follow each step **`VERY`** carefully ...

1. Shut down any local running instance of mongodb on MackBook

2. Connect/create tunnel to develop mongo, thus:
    ```shell
    dp ssh develop mongodb 1 -v -- -L 27017:mongodb-1:27017
    ```
3. Determine the `< secret key >` in the usual way for mongodb on develop (get MONGODB_BIND_ADDR from dp-configs/secrets/develop/dp-dataset-api-web.json) and use in the following commands on your MacBook that copies the collections into `.json` files:
    ```shell
    mongo mongodb://root:< secret key >@localhost:27017 copy-1000-dimension-options.js >dimension-options.json
    mongo mongodb://root:< secret key >@localhost:27017 copy-datasets.js >datasets.json
    mongo mongodb://root:< secret key >@localhost:27017 copy-editions.js >editions.json
    mongo mongodb://root:< secret key >@localhost:27017 copy-instances.js >instances.json
    ```
4. Close the connection to develop mongo

5. Run the `go` code that will create new `.js` scripts for populating local mongdb, thus:
    ```shell
    go run main.go
    ```
6. Start Mongodb on your MackBook

7. Run new `.js` scripts to populate local mongodb as follows:
    ```shell
    mongo insert-editions.js
    mongo insert-instances.js
    mongo insert-datasets.js
    mongo insert-dimension-options.js
    ```

8. Inspect the datasets collection your local mongodb with Robo 3T to see that its OK

9. NOTE: The dimension.options collection on develop (as of 15th March 2021) is over 1.5 G Bytes and when it is downloaded takes over an hour to populate mongodb on MacBook, so in step 3. above we call a script that only copies the first 1,000 documents from dimension.options that results in a much more reasonable file of ~ half a megabyte.

    If you want to copy ALL of dimension options, replace the `copy-1000-dimension-options.js` in step 3 with: `copy-dimension-options.js` and repeat steps 1 to 8

10. When you are happy with the results, delete temporary files with:
    ```shell
    rm insert-*.js
    rm *.json
    ```

11. Feel free to adjust this code, etc for your other databases/collections ...