package models

import (
	"API-REST/services/database/models/asset"
	"API-REST/services/database/models/attribute"
	"API-REST/services/database/models/user"

	"github.com/arthurkushman/buildsqlx"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	User      *user.Model
	Asset     *asset.Model
	Attribute *attribute.Model
	// ...
)

func Build(postgresDB *buildsqlx.DB, mongoDB *mongo.Database) {
	User = &user.Model{Db: postgresDB}
	Asset = &asset.Model{Coll: mongoDB.Collection("assets")}
	Attribute = &attribute.Model{Coll: mongoDB.Collection("attributes")}
	// ...
}
