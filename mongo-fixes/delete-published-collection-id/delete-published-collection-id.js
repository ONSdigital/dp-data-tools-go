// delete-published-collection-id.js
//
// if a dataset doc has:
//     next.state:"published" and non-blank    next.collection_id ---> delete latter
//  current.state:"published" and non-blank current.collection_id ---> delete latter

// collection
ds_collection = 'datasets'

// sub-document parts
ds_subdocs = ['next', 'current']

// do a find first, show what would be changed
show_find         = true
// if show_find is true, show_old_ids_only limits output to collection_id
show_old_ids_only = true

// o determines what find() outputs (null shows all)
o = null

// set to false to only use show_find
do_update = true

//////////////////////////

// utility to printjson(r) even if r is a cursor (iterate over it)
function show_(r){
        if (r==null) { print(''); return; }
        if (r.hasNext != undefined) {
                while (r.hasNext()) {
                        printjson(r.next())
                }
        } else {
                printjson(r)
        }
}

//////////////////////////

for (i = 0; i < ds_subdocs.length; i++) {
        // build paths for this sub-doc
        state_path  = ds_subdocs[i] + '.state'
        collid_path = ds_subdocs[i] + '.collection_id'

        // build the modifier: {$unset:{sub_doc.collection_id:""}
        setter_cid              = {}
        setter_cid[collid_path] = ""
        setter                  = {"$unset":setter_cid}

        // build the filter query
        q = {}
        q[state_path]  = "published"
        q[collid_path] = {"$exists":true,"$ne":''}

        show_(
                'collection:' + ds_collection + ' sub-doc:' + ds_subdocs[i] +
                ' state@' + state_path + ' c_id@' + collid_path
        )

        if (show_find) {
                if (show_old_ids_only) {
                        o = {}
                        o[collid_path] = 1
                }
                r = db.getCollection(ds_collection).find(q, o)
                show_(r)
        }

        if (do_update) {
                r = db.getCollection(ds_collection).updateMany(q, setter)
                show_(r)
        }
}
