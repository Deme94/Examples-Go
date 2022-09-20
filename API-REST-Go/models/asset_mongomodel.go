package models

import (
	"time"
)

type Asset struct {
	ID   int       `json:"id"`
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

// DB MODEL ***************************************************************
type AssetModel struct {
	/*db *mongodbType*/
}

func NewAssetModel( /*db mongotype*/ ) *AssetModel {
	return &AssetModel{ /*db*/ }
}

// DB QUERIES ---------------------------------------------------------------
func (a *AssetModel) GetAll(daterange ...time.Time) ([]*Asset, error) {
	return nil, nil
}
func (a *AssetModel) Get(id int) (*Asset, error) {
	return nil, nil
}
func (a *AssetModel) Insert(Asset *Asset) error {
	return nil
}
func (a *AssetModel) Update(Asset *Asset) error {
	return nil
}
func (a *AssetModel) Delete(id int) error {
	return nil
}

// ...

// PAYLOADS (json input) ---------------------------------------------------------------

// ...
