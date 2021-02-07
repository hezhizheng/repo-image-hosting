package routes

import (
	"gitee-image-hosting/bindata"
	"gitee-image-hosting/controller"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"html/template"
	"strings"
)

func InitRoute() *gin.Engine {
	//router := gin.Default()
	router := gin.New()
	//router.LoadHTMLGlob("views/*")

	//加载模板文件
	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}
	router.SetHTMLTemplate(t)

	//router.Static("/static", "./static")

	//加载静态资源，例如网页的css、js
	fs := assetfs.AssetFS{
		Asset:     bindata.Asset,
		AssetDir:  bindata.AssetDir,
		AssetInfo: nil,
		Prefix:    "static", //一定要加前缀
	}
	router.StaticFS("/static", &fs)


	router.GET("/", controller.Index)
	router.POST("/upload", controller.ImgUpload)
	router.GET("/images", controller.Images)
	router.POST("/delete", controller.ImageDel)
	return router
}

//加载模板文件
func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for _, name := range bindata.AssetNames() {
		if !strings.HasSuffix(name, ".html") {
			continue
		}
		asset, err := bindata.Asset(name)
		if err != nil {
			continue
		}
		name := strings.Replace(name, "views/", "", 1) //这里将templates替换下，在控制器中调用就不用加templates/了
		t, err = t.New(name).Parse(string(asset))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}