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

		collection, err := dao.FindCollectionByNameOrId("2u2cz76lsw1v3rg")
		if err != nil {
			return err
		}

		// update
		edit_deadline := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "gaga20u4",
			"name": "deadline",
			"type": "date",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), edit_deadline); err != nil {
			return err
		}
		collection.Schema.AddField(edit_deadline)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("2u2cz76lsw1v3rg")
		if err != nil {
			return err
		}

		// update
		edit_deadline := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "gaga20u4",
			"name": "available_date",
			"type": "date",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), edit_deadline); err != nil {
			return err
		}
		collection.Schema.AddField(edit_deadline)

		return dao.SaveCollection(collection)
	})
}
