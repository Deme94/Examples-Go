package models

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MAIN STRUCT
type Asset struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	Name      string             `bson:"name"`
	Date      time.Time          `bson:"date"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	// ... description, etc
}

// other structs

// ...

// DB COLLECTION ***************************************************************
type AssetModel struct {
	Coll *mongo.Collection
}

func NewAssetModel(coll *mongo.Collection) *AssetModel {
	return &AssetModel{coll}
}

// DB QUERIES ---------------------------------------------------------------
func (m *AssetModel) GetAll(fromDate time.Time, toDate time.Time, filterOptions map[string]interface{}) ([]*Asset, error) {
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

	if len(filterOptions) != 0 {
		if name, ok := filterOptions["name"]; ok {
			filter = append(filter, bson.E{"name", name})
		}
		// other options ...
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
func (m *AssetModel) GetNames(fromDate time.Time, toDate time.Time) ([]*Asset, error) {
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

	projection := bson.D{
		{"name", 1}, // select name field
		{"_id", 0},  // remove _id field from selection (it is always returned by default)
	}
	cur, err := m.Coll.Find(context.TODO(),
		filter,
		options.Find().SetProjection(projection),
	)
	if err != nil {
		return nil, err
	}

	var assets []*Asset
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
func (m *AssetModel) InsertMany(assets []*Asset) error {
	var documents []interface{}
	for _, a := range assets {
		documents = append(documents,
			bson.D{
				{"name", a.Name},
				{"date", a.Date},
				{"created_at", time.Now()},
			})
	}
	_, err := m.Coll.InsertMany(context.TODO(), documents)
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
