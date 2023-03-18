package middleware_test

import (
	"API-REST/api-gateway"
	"API-REST/api-gateway/utilities/test"
	"API-REST/services/conf"
	"API-REST/services/database"
	"API-REST/services/logger"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var app *fiber.App

var basePath string

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
