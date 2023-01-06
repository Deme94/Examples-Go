package attribute

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MAIN STRUCT
type Attribute struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	Metadata  AttributeMetadata  `bson:"metadata"`
	Timestamp time.Time          `bson:"timestamp"`
	Value     float64            `bson:"value"`
}
type AttributeMetadata struct {
	AssetName string `bson:"asset_name"`
	Name      string `bson:"name"`
	Label     string `bson:"label"`
	Unit      string `bson:"unit"`
}

// ...

// DB COLLECTION ***************************************************************
type Model struct {
	Coll *mongo.Collection
}

func New(coll *mongo.Collection) (*Model, error) {
	return &Model{Coll: coll}, nil
}
