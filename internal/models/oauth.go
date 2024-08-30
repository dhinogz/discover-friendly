package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
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

var ErrNoOAuthRows = errors.New("no rows in oauth query")

func GetOAuthByUserId(dao *daos.Dao, userId, provider string) (*OAuth, error) {
	o := &OAuth{}
	oq := dao.ModelQuery(o)

	if err := oq.AndWhere(dbx.HashExp{"user": userId, "provider": provider}).Limit(1).One(o); err != nil {
		if err == sql.ErrNoRows {
			o.UserId = userId
			o.Provider = provider
			return o, ErrNoOAuthRows
		}
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

func (o *OAuth) UpdateOAuth(dao *daos.Dao, token *oauth2.Token) error {
	o.AccessToken = token.AccessToken
	o.RefreshToken = token.RefreshToken
	o.TokenType = token.TokenType

	expiry := types.DateTime{}
	if err := expiry.Scan(token.Expiry); err != nil {
		return fmt.Errorf("scanning expire token time: %w", err)
	}
	o.Expiry = expiry

	if err := dao.Save(o); err != nil {
		return fmt.Errorf("update oauth (%s): %w", o.Id, err)
	}

	return nil
}
