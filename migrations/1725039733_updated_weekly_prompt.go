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

		collection.Name = "playlist_prompt"

		// remove
		collection.Schema.RemoveField("gaga20u4")

		// remove
		collection.Schema.RemoveField("ymg9qldz")

		// add
		new_is_spotify_playlist_generated := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "k8auofig",
			"name": "is_spotify_playlist_generated",
			"type": "bool",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {}
		}`), new_is_spotify_playlist_generated); err != nil {
			return err
		}
		collection.Schema.AddField(new_is_spotify_playlist_generated)

		// update
		edit_title := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "ofuvoka6",
			"name": "title",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_title); err != nil {
			return err
		}
		collection.Schema.AddField(edit_title)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("2u2cz76lsw1v3rg")
		if err != nil {
			return err
		}

		collection.Name = "weekly_prompt"

		// add
		del_deadline := &schema.SchemaField{}
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
		}`), del_deadline); err != nil {
			return err
		}
		collection.Schema.AddField(del_deadline)

		// add
		del_week_number := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "ymg9qldz",
			"name": "week_number",
			"type": "number",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"min": 1,
				"max": 52,
				"noDecimal": true
			}
		}`), del_week_number); err != nil {
			return err
		}
		collection.Schema.AddField(del_week_number)

		// remove
		collection.Schema.RemoveField("k8auofig")

		// update
		edit_title := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "ofuvoka6",
			"name": "name",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_title); err != nil {
			return err
		}
		collection.Schema.AddField(edit_title)

		return dao.SaveCollection(collection)
	})
}
