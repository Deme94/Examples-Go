package asset

import "go.mongodb.org/mongo-driver/mongo"

type Model struct {
	Coll *mongo.Collection
}

func New(coll *mongo.Collection) (*Model, error) {
	return &Model{Coll: coll}, nil
}
