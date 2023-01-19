package asset

import (
	"API-REST/services/database/mongo/models/attribute"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Asset struct {
	ID         primitive.ObjectID     `json:"id" bson:"_id, omitempty"`
	Name       string                 `json:"name,omitempty" bson:"name"`
	Date       time.Time              `json:"date,omitempty" bson:"date"`
	CreatedAt  time.Time              `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at,omitempty" bson:"updated_at"`
	Attributes []*attribute.Attribute `json:"attributes,omitempty" bson:"attributes"`
	// ...
}
