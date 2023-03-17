package user

import (
	"API-REST/services/conf"
	"API-REST/services/database"
	"API-REST/services/logger"
	"API-REST/services/mail"
	"API-REST/services/sms"
	"API-REST/services/storage"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// Setup
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Set default log flags (print file and line)

	// Conf
	log.Println("Loading configuration service...")
	err := conf.Setup("../../../test.env", "../../../conf.test")
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

	// SMS
	log.Println("Loading sms service...")
	err = sms.Setup()
	if err != nil {
		log.Fatal("\033[31m"+"SMS SERVICE FAILED"+"\033[0m"+" -> ", err)
	}
	log.Println("\033[32m" + "SMS SERVICE IS RUNNING" + "\033[0m")

	// Storage
	log.Println("Loading storage service...")
	err = storage.SetupLocal()
	if err != nil {
		log.Fatal("\033[31m"+"STORAGE SERVICE FAILED"+"\033[0m"+" -> ", err)
	}
	// err = storage.SetupGCS()
	// if err != nil {
	// 	log.Fatal("\033[31m"+"STORAGE SERVICE FAILED"+"\033[0m"+" -> ", err)
	// }
	log.Println("\033[32m" + "STORAGE SERVICE IS RUNNING" + "\033[0m")

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

	// RUN TESTS
	code := m.Run()
	os.Exit(code)
}

func TestAll(t *testing.T) {
	testHelloWorld(t)
}

func testHelloWorld(t *testing.T) {
	t.Fatal("Hello world!")
}
