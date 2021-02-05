package routes

import (
	"gitee-image-hosting/controller"
	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {
	router := gin.Default()
	//router := gin.New()
	router.LoadHTMLGlob("views/*")
	router.Static("/static", "./static")


	router.GET("/", controller.Index)
	router.POST("/upload", controller.ImgUpload)
	router.GET("/images", controller.Images)
	return router
}
