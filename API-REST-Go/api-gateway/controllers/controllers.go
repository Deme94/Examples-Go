package controllers

import (
	"API-REST/api-gateway/controllers/asset"
	"API-REST/api-gateway/controllers/attribute"
	"API-REST/api-gateway/controllers/feature"
	"API-REST/api-gateway/controllers/permission"
	"API-REST/api-gateway/controllers/role"
	"API-REST/api-gateway/controllers/user"
	"API-REST/api-gateway/controllers/user/auth"
	mongo "API-REST/services/database/mongo/models"
	psql "API-REST/services/database/postgres/models"

	"github.com/go-playground/validator/v10"
)

var (
	User       *user.Controller
	Role       *role.Controller
	Permission *permission.Controller
	Feature    *feature.Controller
	Asset      *asset.Controller
	Attribute  *attribute.Controller
	// ...
)

func Build() {
	validate := validator.New()

	User = &user.Controller{Validate: validate, Model: psql.User, Auth: &auth.Controller{Validate: validate, Model: psql.User}}
	Role = &role.Controller{Validate: validate, Model: psql.Role}
	Permission = &permission.Controller{Validate: validate, Model: psql.Permission}
	Feature = &feature.Controller{Validate: validate, Model: psql.Feature}
	Asset = &asset.Controller{Validate: validate, Model: mongo.Asset}
	Attribute = &attribute.Controller{Validate: validate, Model: mongo.Attribute}
	// ...
}
