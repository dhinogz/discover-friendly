package models

import (
	"fmt"
	"log/slog"

	"github.com/pocketbase/pocketbase/daos"
	pbmodels "github.com/pocketbase/pocketbase/models"
)

func CreateUser(dao *daos.Dao, email, username string) (*pbmodels.Record, error) {
	collection, err := dao.FindCollectionByNameOrId("users")
	if err != nil {
		return nil, err
	}

	record := pbmodels.NewRecord(collection)

	// set individual fields
	// or bulk load with record.Load(map[string]any{...})
	if err := record.SetEmail(email); err != nil {
		slog.Error("setting email failed", "error", err)
		return nil, fmt.Errorf("Internal error")
	}
	if err := record.SetUsername(username); err != nil {
		slog.Error("setting username failed", "error", err)
		return nil, fmt.Errorf("Internal error")
	}

	if err := dao.SaveRecord(record); err != nil {
		return nil, err
	}
	return record, nil
}
