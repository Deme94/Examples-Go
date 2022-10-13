package models

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Asset struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	Name      string             `json:"name"`
	Date      time.Time          `json:"date"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	// ...
}

// DB COLLECTION ***************************************************************
type AssetModel struct {
	Coll *mongo.Collection
}

func NewAssetModel(coll *mongo.Collection) *AssetModel {
	return &AssetModel{coll}
}

// DB QUERIES ---------------------------------------------------------------
func (m *AssetModel) GetAll(fromDate time.Time, toDate time.Time) ([]*Asset, error) {
	filter := bson.D{}
	filterDate := bson.D{}

	if !fromDate.IsZero() {
		filterDate = append(filterDate, bson.E{"$gt", fromDate})
		filter = bson.D{{"date", filterDate}}
	}
	if !toDate.IsZero() {
		filterDate = append(filterDate, bson.E{"$lt", toDate})
		filter = bson.D{{"date", filterDate}}
	}

	var assets []*Asset
	cur, err := m.Coll.Find(context.TODO(),
		filter,
	)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var asset Asset
		err := cur.Decode(&asset)
		if err != nil {
			log.Fatal(err)
		}

		assets = append(assets, &asset)

	}
	//Close the cursor once finished
	cur.Close(context.TODO())

	if len(assets) == 0 {
		return nil, errors.New("no assets found")
	}

	return assets, nil
}
func (m *AssetModel) Get(id string) (*Asset, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var asset Asset
	r := m.Coll.FindOne(context.TODO(),
		bson.M{
			"_id": objID,
		},
	)
	err = r.Decode(&asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}
func (m *AssetModel) Insert(asset *Asset) error {
	_, err := m.Coll.InsertOne(context.TODO(),
		bson.D{
			{"name", asset.Name},
			{"date", asset.Date},
			{"created_at", time.Now()},
		},
	)
	return err
}
func (m *AssetModel) Update(asset *Asset) error {
	_, err := m.Coll.UpdateOne(
		context.TODO(),
		bson.D{
			{"_id", asset.ID},
		},
		bson.D{
			{"$set", bson.D{
				{"name", asset.Name},
				{"date", asset.Date},
				{"updated_at", time.Now()},
			}},
		},
	)

	return err
}
func (m *AssetModel) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = m.Coll.DeleteOne(
		context.TODO(),
		bson.D{
			{"_id", objID},
		},
	)

	return err
}

// ...

// PAYLOADS (json input) ---------------------------------------------------------------

// ...
