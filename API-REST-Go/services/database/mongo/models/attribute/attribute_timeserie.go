package attribute

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
