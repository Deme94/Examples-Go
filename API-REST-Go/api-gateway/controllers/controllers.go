package controllers

import (
	"API-REST/api-gateway/controllers/asset"
	"API-REST/api-gateway/controllers/attribute"
	"API-REST/api-gateway/controllers/user"
	mongo "API-REST/services/database/mongo/models"
	psql "API-REST/services/database/postgres/models"
)

var (
	User      *user.Controller
	Asset     *asset.Controller
	Attribute *attribute.Controller
	// ...
)

func Build() {
	User = &user.Controller{Model: psql.User}
	Asset = &asset.Controller{Model: mongo.Asset}
	Attribute = &attribute.Controller{Model: mongo.Attribute}
	// ...
}
