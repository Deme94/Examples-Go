package user

import (
	"API-REST/api-gateway/controllers/user/payloads"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/conf"
	"API-REST/services/database/models/user"
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
	"golang.org/x/crypto/bcrypt"
)

func (c *Controller) GetAll(ctx *gin.Context) {

	var usrs []*user.User
	var err error

	// if query parameters
	y := ctx.Query("year")
	if len(y) != 0 {
		usrs, err = c.Model.GetAllByYear(y)
	} else {
		usrs, err = c.Model.GetAll()

	}
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	all := payloads.GetAllResponse{Users: usrs}
	util.WriteJSON(ctx, http.StatusOK, all, "users")
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

// Example secured route with different behaviour for role user and role admin
func (c *Controller) GetSecuredUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(ctx, err)
		return
	}

	if ctx.Param("Claimer-Role") == "admin" || fmt.Sprint(id) == ctx.GetString("Claimer-ID") {
		u, err := c.Model.Get(id)
		if err != nil {
			util.ErrorJSON(ctx, err)
			return
		}
		util.WriteJSON(ctx, http.StatusOK, u, "user")
	} else {
		util.ErrorJSON(ctx, errors.New("unauthorized - cannot get other user's data"), http.StatusForbidden)
	}
}

// Example secured route only for admins - already checked by middleware
func (c *Controller) GetSecuredAdmin(ctx *gin.Context) {
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
func (c *Controller) Insert(ctx *gin.Context) {
	var req payloads.LoginRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	err = c.Model.Insert(&user.User{Name: req.Name, Email: req.Email, Password: string(hashedPassword)})
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	var u user.User
	u.ID = id
	u.Name = req.Name // cambiar por ruta del archivo creado
	u.Email = req.Email
	u.Password = string(hashedPassword)

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
		fmt.Println("HOLA3")
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

func (c *Controller) Login(ctx *gin.Context) {
	var req payloads.LoginRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	var u *user.User
	if len(req.Email) != 0 {
		u, err = c.Model.GetByEmailWithPassword(req.Email)
		if err != nil {
			util.ErrorJSON(ctx, err, http.StatusUnauthorized)
			return
		}
	} else {
		u, err = c.Model.GetByNameWithPassword(req.Name)
		if err != nil {
			util.ErrorJSON(ctx, err, http.StatusUnauthorized)
			return
		}
	}

	hashedPassword := u.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		util.ErrorJSON(ctx, err, http.StatusUnauthorized)
		return
	}

	// Generate jwt token after successful login
	token, err := c.generateJwtToken(fmt.Sprint(u.ID), conf.JwtSecret)
	if err != nil {
		util.ErrorJSON(ctx, err, http.StatusNotImplemented)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, payloads.LoginResponse{Id: u.ID, Token: string(token)}, "response")
}
