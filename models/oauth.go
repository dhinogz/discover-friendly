package models

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type OAuth struct {
	models.BaseModel

	// Fields
	Provider     string         `db:"provider"`
	AccessToken  string         `db:"access_token"`
	RefreshToken string         `db:"refresh_token"`
	TokenType    string         `db:"token_type"`
	Expiry       types.DateTime `db:"expiry"`

	// Relations
	UserId string `db:"user"`
	// User   *User  `db:"-"`
}

func (o *OAuth) Save(dao *daos.Dao) error {
	if err := dao.Save(o); err != nil {
		return fmt.Errorf("save oauth (%s): %w", o.Id, err)
	}

	return nil
}

func (*OAuth) TableName() string {
	return "oauth"
}

func OAuthQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&OAuth{})
}

func GetOAuthByUserId(dao *daos.Dao, userId, provider string) (*OAuth, error) {
	o := &OAuth{}
	oq := dao.ModelQuery(o)

	if err := oq.AndWhere(dbx.HashExp{"user": userId, "provider": provider}).Limit(1).One(o); err != nil {
		return nil, fmt.Errorf("db query oauth (%s): %w", userId, err)
	}

	return o, nil
}
func DeleteOAuthByUserId(dao *daos.Dao, userId, provider string) error {
	o := &OAuth{}
	oq := dao.ModelQuery(o)
	if err := oq.AndWhere(dbx.HashExp{"user": userId, "provider": provider}).Limit(1).One(o); err != nil {
		return fmt.Errorf("db query oauth (%s): %w", userId, err)
	}
	if err := dao.Delete(o); err != nil {
		return err
	}
	return nil
}
