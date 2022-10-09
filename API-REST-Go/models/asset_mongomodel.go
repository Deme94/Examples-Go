package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Asset struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
func (a *AssetModel) GetAll(daterange time.Time) ([]*Asset, error) {
	var assets []*Asset
	r, err := a.Coll.Find(context.TODO(),
		bson.D{
			{"date", bson.D{{"$gt", "date"}}},
		},
	)
	// err := r.Decode(asset)
	// if err != nil {
	// 	return nil, err
	// }
	return assets, nil
}
func (a *AssetModel) Get(id int) (*Asset, error) {
	var asset *Asset
	r := a.Coll.FindOne(context.TODO(),
		bson.D{
			{"id", id},
		},
	)
	err := r.Decode(asset)
	if err != nil {
		return nil, err
	}
	return asset, nil
}
func (a *AssetModel) Insert(asset *Asset) error {
	_, err := a.Coll.InsertOne(context.TODO(),
		bson.D{
			{"name", asset.Name},
			{"date", asset.Date},
			{"created_at", time.Now()},
			{"updated_at", time.Now()},
		},
	)
	return err
}
func (a *AssetModel) Update(asset *Asset) error {
	return nil
}
func (a *AssetModel) Delete(id int) error {
	return nil
}

// ...

// PAYLOADS (json input) ---------------------------------------------------------------

// ...
