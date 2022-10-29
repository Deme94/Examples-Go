package main

import (
	c "API-REST/cmd/api/controllers"
)

type controllers struct {
	user      *c.UserController
	asset     *c.AssetController
	attribute *c.AttributeController
	// ...
}

func (s *server) createControllers(dbs *databases) *controllers {
	return &controllers{
		user:      c.NewUserController(dbs.postgresql, s.logger, secret, domain),
		asset:     c.NewAssetController(dbs.mongodb.Collection("assets"), s.logger),
		attribute: c.NewAttributeController(dbs.mongodb.Collection("attributes"), s.logger),
		// ...
	}
}
