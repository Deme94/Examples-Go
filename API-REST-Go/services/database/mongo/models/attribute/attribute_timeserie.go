package attribute

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attribute struct {
	ID        primitive.ObjectID `json:"id" bson:"_id, omitempty"`
	Metadata  AttributeMetadata  `json:"metadata,omitempty" bson:"metadata"`
	Timestamp *time.Time         `json:"timestamp,omitempty" bson:"timestamp"`
	Value     float64            `json:"value,omitempty" bson:"value"`
}
type AttributeMetadata struct {
	AssetName string `json:"asset_name,omitempty" bson:"asset_name"`
	Name      string `json:"name,omitempty" bson:"name"`
	Label     string `json:"label,omitempty" bson:"label"`
	Unit      string `json:"unit,omitempty" bson:"unit"`
}
