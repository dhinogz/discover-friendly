package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("wvqmiwnd7kq7sjd")
		if err != nil {
			return err
		}

		// add
		new_is_lotw := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "ywk3b7uo",
			"name": "is_lotw",
			"type": "bool",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {}
		}`), new_is_lotw); err != nil {
			return err
		}
		collection.Schema.AddField(new_is_lotw)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("wvqmiwnd7kq7sjd")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("ywk3b7uo")

		return dao.SaveCollection(collection)
	})
}
