package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/sureshk/birthday-server/src/api/controllers"
	"github.com/sureshk/birthday-server/src/api/models"
)

var server = controllers.Server{}
var userInstance = models.User{}

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("./../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()
	os.Exit(m.Run())
}

func Database() {
	var err error
	TestDbDriver := os.Getenv("TEST_DB_DRIVER")
	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_NAME"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)

		if err != nil {
			fmt.Printf("Can not connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error: ", err)
		} else {
			fmt.Printf("We are connected to the %s database \n", TestDbDriver)
		}
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user := models.User{
		Name:     "sureshk",
		Email:    "er.sureshprajapati@gmail.com",
		Birthday: "1990-06-04",
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func seedUsers() ([]models.User, error) {
	users := []models.User{
		models.User{
			Name:     "sureshk",
			Email:    "er.sureshprajapati@gmail.com",
			Birthday: "1990-06-04",
		},
		models.User{
			Name:     "monu",
			Email:    "er.monuprajapati@gmail.com",
			Birthday: "1991-04-05",
		},
	}
	for i := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}
	return users, nil
}
