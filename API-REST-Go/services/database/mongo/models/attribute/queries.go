package attribute

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *Model) GetAll(fromDate time.Time, toDate time.Time, filterOptions map[string]interface{}) ([]*Attribute, error) {
	filter := bson.D{}
	filterDate := bson.D{}

	if !fromDate.IsZero() {
		filterDate = append(filterDate, bson.E{"$gt", fromDate})
		filter = bson.D{{"timestamp", filterDate}}
	}
	if !toDate.IsZero() {
		filterDate = append(filterDate, bson.E{"$lt", toDate})
		filter = bson.D{{"timestamp", filterDate}}
	}

	if len(filterOptions) != 0 {
		if name, ok := filterOptions["name"]; ok {
			filter = append(filter, bson.E{"metadata.name", name})
		}
		if label, ok := filterOptions["label"]; ok {
			filter = append(filter, bson.E{"metadata.label", label})
		}
		// other options ...
	}

	var atts []*Attribute
	cur, err := m.Coll.Find(context.TODO(),
		filter,
	)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var att Attribute
		err := cur.Decode(&att)
		if err != nil {
			log.Fatal(err)
		}

		atts = append(atts, &att)
	}
	//Close the cursor once finished
	cur.Close(context.TODO())

	if len(atts) == 0 {
		return nil, errors.New("no attributes found")
	}

	return atts, nil
}
func (m *Model) Get(id string) (*Attribute, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var att Attribute
	r := m.Coll.FindOne(context.TODO(),
		bson.M{
			"_id": objID,
		},
	)
	err = r.Decode(&att)
	if err != nil {
		return nil, err
	}
	return &att, nil
}
func (m *Model) Insert(attribute *Attribute) error {
	_, err := m.Coll.InsertOne(context.TODO(),
		bson.D{
			{"metadata", bson.D{
				{"asset_name", attribute.Metadata.AssetName},
				{"name", attribute.Metadata.Name},
				{"label", attribute.Metadata.Label},
				{"unit", attribute.Metadata.Unit},
			}},
			{"timestamp", attribute.Timestamp},
			{"value", attribute.Value},
		},
	)
	return err
}
func (m *Model) InsertMany(attributes []*Attribute) error {
	var documents []interface{}
	for _, att := range attributes {
		documents = append(documents,
			bson.D{
				{"metadata", bson.D{
					{"asset_name", att.Metadata.AssetName},
					{"name", att.Metadata.Name},
					{"label", att.Metadata.Label},
					{"unit", att.Metadata.Unit},
				}},
				{"timestamp", att.Timestamp},
				{"value", att.Value},
			})
	}
	_, err := m.Coll.InsertMany(context.TODO(), documents)
	return err
}
func (m *Model) Update(attribute *Attribute) error {
	_, err := m.Coll.UpdateOne(
		context.TODO(),
		bson.D{
			{"_id", attribute.ID},
		},
		bson.D{
			{"$set", bson.D{
				{"metadata.asset_name", attribute.Metadata.AssetName},
				{"metadata.name", attribute.Metadata.Name},
				{"metadata.label", attribute.Metadata.Label},
				{"metadata.unit", attribute.Metadata.Unit},
				{"timestamp", attribute.Timestamp},
				{"value", attribute.Value},
			}},
		},
	)

	return err
}
func (m *Model) Delete(id string) error {
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
