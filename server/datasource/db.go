package datasource

import (
	"trainning/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

// Connect to the DATABASE
func ConnectDatabase() {
	database, err := gorm.Open("mysql", "root:sa123456@/golangdbnew?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("Failed to connect to database!")
	}
	database.AutoMigrate(&models.User{})
	DB = database
}
