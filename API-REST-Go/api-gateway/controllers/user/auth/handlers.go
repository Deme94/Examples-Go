package auth

import (
	"API-REST/api-gateway/controllers/user/auth/payloads"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/conf"
	"API-REST/services/database/postgres/models/user"
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

	"github.com/gin-gonic/gin"
	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

func (c *Controller) Login(ctx *gin.Context) {
	var req payloads.LoginRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	var u *user.User
	if req.Email != "" {
		u, err = c.Model.GetByEmailWithPassword(req.Email)
		if err != nil {
			util.ErrorJSON(ctx, err, http.StatusUnauthorized)
			return
		}
	} else {
		u, err = c.Model.GetByUsernameWithPassword(req.Username)
		if err != nil {
			util.ErrorJSON(ctx, err, http.StatusUnauthorized)
			return
		}
	}

	if !u.DeletedAt.IsZero() {
		util.ErrorJSON(ctx, errors.New("user deleted"), http.StatusUnauthorized)
		return
	}

	hashedPassword := u.Password

	err = c.compareHashAndPassword(hashedPassword, req.Password)
	if err != nil {
		util.ErrorJSON(ctx, err, http.StatusUnauthorized)
		return
	}

	// Generate jwt token after successful login
	token, err := c.generateJwtToken(fmt.Sprint(u.ID), conf.Env.GetString("JWT_SECRET"))
	if err != nil {
		util.ErrorJSON(ctx, err, http.StatusNotImplemented)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, payloads.LoginResponse{ID: u.ID, Token: string(token)}, "response")
}

func (c *Controller) Get(ctx *gin.Context) {
	claimerID := ctx.GetInt("Claimer-ID")

	u, err := c.Model.Get(claimerID)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	userResponse := payloads.GetResponse{
		ID:                 u.ID,
		CreatedAt:          u.CreatedAt,
		Username:           u.Username,
		Email:              u.Email,
		Nick:               u.Nick,
		FirstName:          u.FirstName,
		LastName:           u.LastName,
		Phone:              u.Phone,
		Address:            u.Address,
		LastPasswordChange: u.LastPasswordChange,
		VerifiedMail:       u.VerifiedMail,
		VerifiedPhone:      u.VerifiedPhone,
	}

	util.WriteJSON(ctx, http.StatusOK, userResponse, "user")
}
func (c *Controller) GetPhoto(ctx *gin.Context) {
	claimerID := ctx.GetInt("Claimer-ID")

	imageName, err := c.Model.GetPhoto(claimerID)
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
	claimerID := ctx.GetInt("Claimer-ID")

	cvName, err := c.Model.GetCV(claimerID)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	filePath := "./storage/users/" + cvName + ".pdf"
	ctx.File(filePath)
}
func (c *Controller) Update(ctx *gin.Context) {
	claimerID := ctx.GetInt("Claimer-ID")

	var req payloads.UpdateRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	var u user.User
	u.ID = claimerID
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
func (c *Controller) ChangePassword(ctx *gin.Context) {
	claimerID := ctx.GetInt("Claimer-ID")

	var req payloads.ChangePasswordRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	hashedPassword, err := c.Model.GetPassword(claimerID)
	if err != nil {
		util.ErrorJSON(ctx, err, http.StatusInternalServerError)
		return
	}

	err = c.compareHashAndPassword(hashedPassword, req.OldPassword)
	if err != nil {
		util.ErrorJSON(ctx, err, http.StatusUnauthorized)
		return
	}

	hashedNewPassword, err := c.hashPassword(req.NewPassword)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	err = c.Model.UpdatePassword(claimerID, hashedNewPassword)
	if err != nil {
		util.ErrorJSON(ctx, err, http.StatusInternalServerError)
		return
	}

	ok := payloads.OkResponse{
		OK: true,
	}
	util.WriteJSON(ctx, http.StatusOK, ok, "OK")
}
func (c *Controller) ResetPassword(ctx *gin.Context) {
	var req payloads.ResetPasswordRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	u, err := c.Model.GetByEmailWithPassword(req.Email)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	password := c.generateRandomPassword()

	err = c.Model.UpdatePassword(u.ID, password)
	if err != nil {
		util.ErrorJSON(ctx, err, http.StatusInternalServerError)
		return
	}

	// Send password to email
	fmt.Println(password)
	// ...

	ok := payloads.OkResponse{
		OK: true,
	}
	util.WriteJSON(ctx, http.StatusOK, ok, "OK")
}
func (c *Controller) UpdateRoles(ctx *gin.Context) {
	claimerID := ctx.GetInt("Claimer-ID")

	var req payloads.UpdateRolesRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	err = c.Model.UpdateRoles(claimerID, req.RoleIDs...)
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
	claimerID := ctx.GetInt("Claimer-ID")

	var req payloads.UpdatePhotoRequest

	err := ctx.BindJSON(&req)
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
	imageName := "user" + fmt.Sprint(claimerID)
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

	err = c.Model.UpdatePhoto(claimerID, imageName)
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
	claimerID := ctx.GetInt("Claimer-ID")

	// Retrieve the file
	var req payloads.UpdateCVRequest
	err := ctx.ShouldBind(&req)
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
	fileName := "usercv" + fmt.Sprint(claimerID)
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
	c.Model.UpdateCV(claimerID, fileName)

	ok := payloads.OkResponse{
		OK: true,
	}
	util.WriteJSON(ctx, http.StatusOK, ok, "OK")
}
func (c *Controller) Delete(ctx *gin.Context) {
	claimerID := ctx.GetInt("Claimer-ID")

	err := c.Model.Delete(claimerID)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	ok := payloads.OkResponse{
		OK: true,
	}
	util.WriteJSON(ctx, http.StatusOK, ok, "OK")
}
