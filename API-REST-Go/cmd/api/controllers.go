package main

import (
	c "API-REST/controllers"
)

type controllers struct {
	user  *c.UserController
	asset *c.AssetController
	// ...
}

func (s *server) createControllers(dbs *databases) *controllers {
	return &controllers{
		user: c.NewUserController(dbs.postgresql, s.logger, secret, domain),
		//asset: &c.AssetController{ /*dbs.mongodb*/ },
		// ...
	}
}
