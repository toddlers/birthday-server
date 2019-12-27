package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/sureshk/birthday-server/src/api/models"
)

var users = []models.User{
	models.User{
		Name:     "sureshk",
		Birthday: "1990-06-04",
	},
	models.User{
		Name:     "monu",
		Birthday: "1994-06-04",
	},
	models.User{
		Name:     "alex",
		Birthday: "1970-01-04",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		log.Fatalf("can not drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		log.Fatalf("can not drop table: %v", err)
	}
	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("can not drop table: %v", err)
		}
	}
}
