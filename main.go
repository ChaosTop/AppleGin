package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	TOKEN = ""
	PORT  = ""
)

var AppleGin *gorm.DB
var err error

type Account struct {
	gorm.Model
	Email    string
	Password string
}

type PostRecord struct {
	Email    string `form:"email"`
	Password string `form:"password"`
	Token    string `form:"token"`
}

func init() {
	viper.AddConfigPath("./")
	viper.SetConfigFile("config.yaml")
	viper.ReadInConfig()

	TOKEN = viper.GetString("TOKEN")
	PORT = viper.GetString("PORT")

	AppleGin, err = gorm.Open(sqlite.Open("appleid.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	AppleGin.AutoMigrate(&Account{})
}

func main() {

	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.LoadHTMLGlob("./*")

	r.GET("/", func(c *gin.Context) {
		var data Account
		AppleGin.Last(&data)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"AppleID":  data.Email,
			"Password": data.Password,
		})
	})
	r.GET("/all", getAllRecords)
	r.GET("/latest", getLatestRecord)
	r.POST("/account", createRecord)
	r.Run(fmt.Sprint(":" + PORT)) // listen and serve on 0.0.0.0:8080
}

func getAllRecords(c *gin.Context) {
	var data []Account
	AppleGin.Find(&data)
	c.JSON(http.StatusOK, gin.H{"message": data})
}

func getLatestRecord(c *gin.Context) {
	var data Account
	AppleGin.Last(&data)
	c.JSON(http.StatusOK, gin.H{"message": data})
}

func createRecord(c *gin.Context) {
	var data PostRecord
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if data.Token != TOKEN {
		c.JSON(http.StatusForbidden, gin.H{"error": "Authentication Failed!"})
		return
	}
	AppleGin.Create(&Account{Email: data.Email, Password: data.Password})
	c.JSON(http.StatusOK, gin.H{"message": "AppleID " + data.Email + ", Password " + data.Password})

}
