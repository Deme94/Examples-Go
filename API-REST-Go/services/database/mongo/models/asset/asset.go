package asset

import (
	"API-REST/services/database/mongo/models/attribute"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Asset struct {
	ID         primitive.ObjectID     `bson:"_id, omitempty"`
	Name       string                 `bson:"name"`
	Date       time.Time              `bson:"date"`
	CreatedAt  time.Time              `bson:"created_at"`
	UpdatedAt  time.Time              `bson:"updated_at"`
	Attributes []*attribute.Attribute `bson:"attributes"`
	// ...
}
