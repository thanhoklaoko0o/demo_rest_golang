package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"trainning/datasource"
	"trainning/models"
	"trainning/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	router := gin.Default()

	router.Use(aCustomMiddleware())

	router.GET("/ping", getPing)
	router.POST("/ping", postPing)
	router.GET("/detail/:id", getDetail)

	// Access data from db

	datasource.ConnectDatabase()
	router.GET("/user", GetUsers)
	router.GET("/user/:id", GetUser)
	router.POST("/user", CreateUser)
	router.PUT("/user/:id", UpdateUser)
	router.DELETE("/user/:id", DeleteUser)

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			//Khởi tạo middleware chỉ dùng cho group có prefix là v1
			v1.Use(groupV1Middleware())
			v1.GET("/ping", func(context *gin.Context) {
				context.JSON(http.StatusOK, gin.H{
					"message": "Ping",
				})
			})
			v1.GET("/pong", func(context *gin.Context) {
				context.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})
		}
		v2 := api.Group("/v2")
		{
			v2.GET("/a", func(context *gin.Context) {
				context.JSON(http.StatusOK, gin.H{
					"message": "a",
				})
			})
			v2.GET("/b", func(context *gin.Context) {
				context.JSON(http.StatusOK, gin.H{
					"message": "b",
				})
			})
		}
	}
	//upload single file
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(context *gin.Context) {
		// single file
		file, _ := context.FormFile("file")
		log.Println(file.Filename)

		// Upload the file to specific dst.
		context.SaveUploadedFile(file, "./assets/upload/"+file.Filename)
		// context.SaveUploadedFile(file, "./assets/upload/"+uuid.New().String()+filepath.Ext(file.Filename))

		context.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	//upload multiple file
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload_multiple_file", func(context *gin.Context) {
		// Multipart form
		form, _ := context.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			// Upload the file to specific dst.
			context.SaveUploadedFile(file, "./assets/upload/"+uuid.New().String()+filepath.Ext(file.Filename))
		}
		context.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})
	router.Run(":333")
}
func groupV1Middleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		log.Println("I'm in a group v1 middleware")
		context.Next()
	}
}

func aCustomMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		log.Println("I'm in a global middleware")
		// jwt
		// auth
		if true {
			context.Next()
		}
	}
}

// func aCustomMiddleware(context *gin.Context) {
// 	// jwt
// 	// auth
// 	log.Println("This is global middleware")
// 	context.Next()
// }

func getDetail(context *gin.Context) {
	id := context.Param("id")
	context.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func postPing(context *gin.Context) {
	address := context.DefaultPostForm("addr", "VietNam")
	context.JSON(http.StatusOK, gin.H{
		"message": "Hello from " + address + " POST ping",
	})
}

func getPing(context *gin.Context) {
	name := context.DefaultQuery("name", "guest")
	var data = map[string]interface{}{
		"message": "Hello " + name + "from GET ping",
	}
	// context.String(http.StatusOK, "Ping")
	//context.JSON(http.StatusOK, gin.H{"message": "ok..........."})
	context.JSON(http.StatusOK, data)
}

// GetUsers from the DB
func GetUsers(context *gin.Context) {

	var users []models.User
	datasource.DB.Find(&users)
	response.JSON(context, http.StatusOK, users)

}

// CreateUser from the DB
func CreateUser(context *gin.Context) {
	var input models.User
	if err := context.ShouldBindJSON(&input); err != nil {
		response.ERROR(context, http.StatusBadRequest)
		return
	}
	// Create user
	user := models.User{Nickname: input.Nickname, Email: input.Email, Password: input.Password}
	datasource.DB.Create(&user)
	response.JSON(context, http.StatusOK, user)
}

// GetUser from the DB
func GetUser(context *gin.Context) {
	// Get model if exist
	var user models.User
	if err := datasource.DB.Where("id = ?", context.Param("id")).First(&user).Error; err != nil {
		response.ERROR(context, http.StatusBadRequest)
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	response.JSON(context, http.StatusOK, user)
}

//  UpdateUser from the DB
func UpdateUser(context *gin.Context) {
	// Get model if exist
	var user models.User
	if err := datasource.DB.Where("id = ?", context.Param("id")).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input models.User
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	datasource.DB.Model(&user).Updates(input)
	response.JSON(context, http.StatusOK, user)
}

//  DeleteUser from the DB
func DeleteUser(context *gin.Context) {
	// Get model if exist
	var user models.User
	if err := datasource.DB.Where("id = ?", context.Param("id")).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	datasource.DB.Delete(&user)

	response.JSON(context, http.StatusOK, true)
}
