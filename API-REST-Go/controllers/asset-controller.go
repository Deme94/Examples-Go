package controllers

import (
	"log"
	"net/http"
)

// CONTROLLER ***************************************************************
type AssetController struct {
	/*db *mongodbType*/
	logger *log.Logger
}

func NewAssetController( /*mongdb*/ logger *log.Logger) *AssetController {
	c := AssetController{}
	//c.model = &m.UserModel{Db: db}
	c.logger = logger

	return &c
}

// METHODS CONTROLLER ---------------------------------------------------------------
// ...

// API HANDLERS ---------------------------------------------------------------
func (c *AssetController) GetAll(w http.ResponseWriter, r *http.Request) {

}
func (c *AssetController) Get(w http.ResponseWriter, r *http.Request) {

}
func (c *AssetController) Insert(w http.ResponseWriter, r *http.Request) {

}
func (c *AssetController) Update(w http.ResponseWriter, r *http.Request) {

}
func (c *AssetController) Delete(w http.ResponseWriter, r *http.Request) {

}

// ...
