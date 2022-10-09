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
func (a *AssetModel) GetAll(daterange ...time.Time) ([]*Asset, error) {
	return nil, nil
}
func (a *AssetModel) Get(id int) (*Asset, error) {
	return nil, nil
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
	// OLA
	return nil
}

// ...

// PAYLOADS (json input) ---------------------------------------------------------------

// ...
