package middleware_test

import (
	"API-REST/api-gateway"
	authPayloads "API-REST/api-gateway/controllers/user/auth/payloads"
	userPayloads "API-REST/api-gateway/controllers/user/payloads"
	"API-REST/api-gateway/utilities/test"
	"API-REST/services/conf"
	"API-REST/services/database"
	"API-REST/services/logger"
	"API-REST/services/mail"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var app *fiber.App

var basePath string
var headers map[string]string
var cookies []*http.Cookie

var testUserAdmin = userPayloads.InsertRequest{
	Username: "admin",
	Email:    "admin@gmail.com",
	Password: "test1234",
}
var testUser = userPayloads.InsertRequest{
	Username: "test",
	Email:    "test@gmail.com",
	Password: "test1234",
}

func TestMain(m *testing.M) {
	// Setup
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Set default log flags (print file and line)

	// Conf
	log.Println("Loading configuration service...")
	err := conf.Setup("../../test.env", "../../conf.test")
	if err != nil {
		log.Fatal("\033[31m"+"CONFIGURATION SERVICE FAILED"+"\033[0m"+" -> ", err)
	}
	log.Println("\033[32m" + "CONFIGURATION SERVICE IS RUNNING" + "\033[0m")

	// Logger
	log.Println("Loading logging service...")
	err = logger.Setup()
	if err != nil {
		log.Fatal("\033[31m"+"LOGGING SERVICE FAILED"+"\033[0m"+" -> ", err)
	}
	log.Println("\033[32m" + "LOGGING SERVICE IS RUNNING" + "\033[0m")

	// Mail
	log.Println("Loading mail service...")
	err = mail.Setup()
	if err != nil {
		log.Fatal("\033[31m"+"MAIL SERVICE FAILED"+"\033[0m"+" -> ", err)
	}
	log.Println("\033[32m" + "MAIL SERVICE IS RUNNING" + "\033[0m")

	// DB
	log.Println("Loading database service...")
	err = database.SetupPostgresDockertest()
	if err != nil {
		log.Fatal("\033[31m"+"DATABASE SERVICE FAILED"+"\033[0m"+" -> ", err)
	}
	// err = database.SetupMongo()
	// if err != nil {
	// 	log.Fatal("\033[31m"+"DATABASE SERVICE FAILED"+"\033[0m"+" -> ", err)
	// }
	log.Println("\033[32m" + "DATABASE SERVICE IS RUNNING" + "\033[0m")

	// Build Test App (api)
	app = api.NewRouter()

	basePath = "http://" + conf.Env.GetString("HOST") + ":" + conf.Env.GetString("PORT") + conf.Conf.GetString("apiBasePath")

	// RUN TESTS
	code := m.Run()
	os.Exit(code)
}

func TestAll(t *testing.T) {
	testCORS(t)
	testAuth(t)
}

func testCORS(t *testing.T) {
	msTimeout := 2000
	req := test.NewRequest(&test.RequestParams{
		Method: "GET",
		Path:   basePath + "/",
	})

	// Test allowed Origin
	req.Header.Add("Origin", "test.com")
	res, err := app.Test(req, msTimeout)
	if err != nil {
		log.Fatal(err)
	}
	assert.Contains(t, fmt.Sprint(res.Header), "Origin:[test.com]")

	// Test unallowed Origin
	req.Header.Del("Origin")
	req.Header.Add("Origin", "hacker.com")
	res, err = app.Test(req, msTimeout)
	if err != nil {
		log.Fatal(err)
	}
	assert.NotContains(t, fmt.Sprint(res.Header), "Origin:[test.com]")
	assert.Contains(t, fmt.Sprint(res.Header), "Origin:[]")
}

func testAuth(t *testing.T) {
	msTimeout := 2000
	res, err := app.Test(test.NewRequest(&test.RequestParams{
		Method: "GET",
		Path:   basePath + "/private/",
	}), msTimeout)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	signUpTestUsers()
	loginTestUserAdmin()

	res, err = app.Test(test.NewRequest(&test.RequestParams{
		Method:  "GET",
		Path:    basePath + "/private/",
		Headers: headers,
		Cookies: cookies,
	}), msTimeout)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

// utils
func signUpTestUsers() {
	// signup admin
	msTimeout := 15000
	_, err := app.Test(test.NewRequest(&test.RequestParams{
		Method: "POST",
		Path:   basePath + "/public/users",
		Body:   testUserAdmin,
	}), msTimeout)
	if err != nil {
		log.Fatal(err)
	}

	//signup user
	_, err = app.Test(test.NewRequest(&test.RequestParams{
		Method: "POST",
		Path:   basePath + "/public/users",
		Body:   testUser,
	}), msTimeout)
	if err != nil {
		log.Fatal(err)
	}
}

func loginTestUser() {
	msTimeout := 2000
	res, err := app.Test(test.NewRequest(&test.RequestParams{
		Method: "POST",
		Path:   basePath + "/public/auth/login",
		Body:   testUser,
	}), msTimeout)
	if err != nil {
		log.Fatal(err)
	}
	// Save cookies
	cookies = res.Cookies()
	// Save Jwt Token Header
	loginResponse := authPayloads.LoginResponse{}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	json.Unmarshal(bodyBytes, &loginResponse)

	headers = make(map[string]string)
	headers["Authorization"] = "Bearer " + loginResponse.Token
}

func loginTestUserAdmin() {
	msTimeout := 2000
	res, err := app.Test(test.NewRequest(&test.RequestParams{
		Method: "POST",
		Path:   basePath + "/public/auth/login",
		Body:   testUserAdmin,
	}), msTimeout)
	if err != nil {
		log.Fatal(err)
	}
	// Save cookies
	cookies = res.Cookies()
	// Save Jwt Token Header
	loginResponse := authPayloads.LoginResponse{}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	json.Unmarshal(bodyBytes, &loginResponse)

	headers = make(map[string]string)
	headers["Authorization"] = "Bearer " + loginResponse.Token
}
