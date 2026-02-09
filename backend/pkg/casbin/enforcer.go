// Package casbin - Casbin enforcer
package casbin

import (
	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func InitEnforcer(db *gorm.DB, modelPath string) (*casbin.Enforcer, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, err
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}

	enforcer.EnableAutoSave(true)

	return enforcer, nil
}
