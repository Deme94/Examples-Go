package controllers

import (
	"API-REST/api-gateway/controllers/asset"
	"API-REST/api-gateway/controllers/attribute"
	"API-REST/api-gateway/controllers/user"
	"API-REST/services/database/models"
)

var (
	User      *user.Controller
	Asset     *asset.Controller
	Attribute *attribute.Controller
	// ...
)

func Build() {
	User = &user.Controller{Model: models.User}
	Asset = &asset.Controller{Model: models.Asset}
	Attribute = &attribute.Controller{Model: models.Attribute}
	// ...
}
