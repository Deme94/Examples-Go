package user

import (
	"API-REST/api-gateway/controllers/user/payloads"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/database/postgres/models/user"
	psql "API-REST/services/database/postgres/predicates"
	"API-REST/services/logger"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

func (c *Controller) GetAll(ctx *gin.Context) {

	// Query parameters
	predicates := psql.Predicates{}
	pageParam := ctx.Query("page")
	pageSizeParam := ctx.Query("pageSize")
	yearParam := ctx.Query("year")
	deletedParam := ctx.Query("deleted")
	if len(pageParam) != 0 && len(pageSizeParam) != 0 {
		page, err := strconv.Atoi(pageParam)
		if err != nil {
			util.ErrorJSON(ctx, err)
			return
		}
		pageSize, err := strconv.Atoi(pageSizeParam)
		if err != nil {
			util.ErrorJSON(ctx, err)
			return
		}
		if page == 0 || pageSize == 0 {
			util.ErrorJSON(ctx, errors.New("page and pageSize params must be greater than 0"))
			return
		}
		limit := pageSize
		offset := (page - 1) * pageSize
		predicates.Offset(offset).Limit(limit)
	}
	if len(yearParam) != 0 {
		startDate := fmt.Sprint(yearParam, "-01-01")
		endDate := fmt.Sprint(yearParam, "-12-31")
		predicates.Where("created_at", ">=", startDate).AndWhere("created_at", "<=", endDate)
	}
	if len(deletedParam) != 0 {
		deleted, err := strconv.ParseBool(deletedParam)
		if err == nil {
			if predicates.HasWhere() {
				if deleted {
					predicates.AndWhereNotNull("deleted_at")
				} else {
					predicates.AndWhereNull("deleted_at")
				}
			} else {
				if deleted {
					predicates.WhereNotNull("deleted_at")
				} else {
					predicates.WhereNull("deleted_at")
				}
			}
		}
	}

	usrs, err := c.Model.GetAll(&predicates)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	var all []*payloads.GetResponse
	for _, user := range usrs {
		all = append(all, &payloads.GetResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		})
	}

	allResponse := payloads.GetAllResponse{Users: all}
	util.WriteJSON(ctx, http.StatusOK, allResponse, "response")
}
func (c *Controller) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(ctx, err)
		return
	}

	u, err := c.Model.Get(id)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, u, "user")
}
func (c *Controller) GetPhoto(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(ctx, err)
		return
	}

	imageName, err := c.Model.GetPhoto(id)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	// Get webp file
	f, err := os.OpenFile("./storage/users/"+imageName+".webp", os.O_RDWR, 0644)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}
	defer f.Close()
	// Decode webp file to image
	image, err := webp.Decode(f, &decoder.Options{})
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}
	// Encode image into buffer
	var buf bytes.Buffer
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	if err != nil {
		log.Fatalln(err)
	}
	webp.Encode(&buf, image, options)
	// Get bytes and encode to base64
	imageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	util.WriteJSON(ctx, http.StatusOK, imageBase64, "photo")
}
func (c *Controller) GetCV(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(ctx, err)
		return
	}

	cvName, err := c.Model.GetCV(id)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	filePath := "./storage/users/" + cvName + ".pdf"
	ctx.File(filePath)
}
func (c *Controller) Insert(ctx *gin.Context) {
	var req payloads.InsertRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	hashedPassword, err := c.HashPassword(req.Password)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	err = c.Model.Insert(&user.User{Username: req.Username, Email: req.Email, Password: hashedPassword})
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, "user created successfully", "response")
}
func (c *Controller) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(ctx, err)
		return
	}

	var req payloads.UpdateRequest
	err = ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	var u user.User
	u.ID = id
	u.Nick = req.Nick
	u.FirstName = req.FirstName
	u.LastName = req.LastName
	u.Phone = req.Phone
	u.Address = req.Address

	err = c.Model.Update(&u)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	ok := payloads.OkResponse{
		OK: true,
	}
	util.WriteJSON(ctx, http.StatusOK, ok, "OK")
}
func (c *Controller) UpdateRoles(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(ctx, err)
		return
	}

	var req payloads.UpdateRolesRequest

	err = ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	err = c.Model.UpdateRoles(id, req.RoleIDs...)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	ok := payloads.OkResponse{
		OK: true,
	}
	util.WriteJSON(ctx, http.StatusOK, ok, "OK")
}
func (c *Controller) UpdatePhoto(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(ctx, err)
		return
	}

	var req payloads.UpdatePhotoRequest

	err = ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	// Decode base64 webp to bytes
	unbased, err := base64.StdEncoding.DecodeString(req.PhotoBase64)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}
	// Decode bytes to image
	reader := bytes.NewReader(unbased)
	image, err := webp.Decode(reader, &decoder.Options{})
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}
	// Create our own file
	imageName := "user" + fmt.Sprint(id)
	f, err := os.OpenFile("./storage/users/"+imageName+".webp", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}
	defer f.Close()
	// Encode image into file
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	if err != nil {
		log.Fatalln(err)
	}
	webp.Encode(f, image, options)
	logger.Logger.Println("user's photo saved")

	err = c.Model.UpdatePhoto(id, imageName)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	ok := payloads.OkResponse{
		OK: true,
	}
	util.WriteJSON(ctx, http.StatusOK, ok, "OK")
}
func (c *Controller) UpdateCV(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(ctx, err)
		return
	}

	// Retrieve the file
	var req payloads.UpdateCVRequest
	err = ctx.ShouldBind(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	// Check if file is .pdf
	if filepath.Ext(req.File.Filename) != ".pdf" {
		util.ErrorJSON(ctx, errors.New("file extension must be pdf"))
		return
	}

	// Open file
	file, err := req.File.Open()
	if err != nil {
		util.ErrorJSON(ctx, err)
	}
	// Create our own file
	fileName := "usercv" + fmt.Sprint(id)
	f, err := os.OpenFile("./storage/users/"+fileName+".pdf", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}
	defer f.Close()
	// Save file in storage
	io.Copy(f, file)
	logger.Logger.Printf("user's cv saved. Name: %s | Size: %d", req.File.Filename, req.File.Size)

	// Save fileName in DB
	c.Model.UpdateCV(id, fileName)

	ok := payloads.OkResponse{
		OK: true,
	}
	util.WriteJSON(ctx, http.StatusOK, ok, "OK")
}
func (c *Controller) Restore(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(ctx, err)
		return
	}

	err = c.Model.Restore(id)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	ok := payloads.OkResponse{
		OK: true,
	}
	util.WriteJSON(ctx, http.StatusOK, ok, "OK")
}
func (c *Controller) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(ctx, err)
		return
	}

	err = c.Model.Delete(id)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	ok := payloads.OkResponse{
		OK: true,
	}
	util.WriteJSON(ctx, http.StatusOK, ok, "OK")
}
