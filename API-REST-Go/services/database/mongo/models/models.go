package models

import (
	"API-REST/services/database/mongo/models/asset"
	"API-REST/services/database/mongo/models/attribute"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Asset     *asset.Model
	Attribute *attribute.Model
	// ...
)

func Build(mongoDB *mongo.Database) error {
	var err error

	Asset, err = asset.New(mongoDB.Collection("assets"))
	if err != nil {
		return err
	}
	Attribute, err = attribute.New(mongoDB.Collection("attributes"))
	if err != nil {
		return err
	}
	// ...

	return nil
}
