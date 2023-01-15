package controllers

import (
	"API-REST/api-gateway/controllers/asset"
	"API-REST/api-gateway/controllers/attribute"
	"API-REST/api-gateway/controllers/permission"
	"API-REST/api-gateway/controllers/role"
	"API-REST/api-gateway/controllers/user"
	"API-REST/api-gateway/controllers/user/auth"
	mongo "API-REST/services/database/mongo/models"
	psql "API-REST/services/database/postgres/models"
)

var (
	Auth       *auth.Controller
	User       *user.Controller
	Role       *role.Controller
	Permission *permission.Controller
	Asset      *asset.Controller
	Attribute  *attribute.Controller
	// ...
)

func Build() {
	Auth = &auth.Controller{Model: psql.User}
	User = &user.Controller{Model: psql.User}
	Role = &role.Controller{Model: psql.Role}
	Permission = &permission.Controller{Model: psql.Permission}
	Asset = &asset.Controller{Model: mongo.Asset}
	Attribute = &attribute.Controller{Model: mongo.Attribute}
	// ...
}
