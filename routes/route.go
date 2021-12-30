package routes

import (
	"github.com/gin-gonic/gin"
	"hezhizheng/repo-image-hosting/controller"
	"hezhizheng/repo-image-hosting/static"
	"hezhizheng/repo-image-hosting/views"
	"net/http"
)

func InitRoute() *gin.Engine {
	//router := gin.Default()
	router := gin.New()

	//加载模板文件
	router.SetHTMLTemplate(views.GoTpl)
	// live 模式 打包用
	router.StaticFS("/static", http.FS(static.EmbedStatic))
	// dev 开发用 避免修改静态资源需要重启服务
	//router.StaticFS("/static", http.Dir("static"))

	router.GET("/", controller.Index)
	router.POST("/upload", controller.ImgUpload)
	router.GET("/images", controller.Images)
	router.POST("/delete", controller.ImageDel)
	return router
}
