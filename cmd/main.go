package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	cdn_checker "mirrors_status/pkg/business/cdn-checker"
	"mirrors_status/pkg/config"
	"mirrors_status/pkg/log"
	"mirrors_status/pkg/modules/model"
	"mirrors_status/pkg/modules/service"
	"net/http"
	"strconv"

)

type App struct {
	serverConfig *configs.ServerConf
	cdnChecker *cdn_checker.CDNChecker
}

func Init() (app App) {
	log.Info("Initializing APP")
	var sc configs.ServerConf
	serverConfig := sc.GetConfig()
	app = App{
		serverConfig: serverConfig,
	}

	configs.InitDB(*serverConfig)
	app.cdnChecker = cdn_checker.NewCDNChecker(app.serverConfig.CdnChecker)
	configs.InitScheme()
	return
}

func(app App) GetAllMirrors(c *gin.Context) {
	data := service.GetAllMirrors()
	c.JSON(200, gin.H{
		"res": data,
	})
}

func(app App) GetAllMirrorsCdn(c *gin.Context) {
	data := service.GetAllMirrorsCdn()
	c.JSON(200, gin.H{
		"res": data,
	})
}

func (app App) AddMirror(c *gin.Context) {
	var reqMirror model.MirrorsPoint
	err := c.ShouldBindJSON(&reqMirror)
	if err != nil {
		log.Errorf("Bind json found error:%v", err)
	}
	err = service.AddMirror(reqMirror)
	if err != nil {
		log.Errorf("Insert data found error:%v", err)
	}
	c.JSON(200, gin.H{
		"res": err,
	})
}

func (app App) AddMirrorCdn(c *gin.Context) {
	var reqMirrorCdn model.MirrorsCdnPoint
	err := c.ShouldBindJSON(&reqMirrorCdn)
	if err != nil {
		log.Errorf("Bind json found error:%v", err)
	}
	err = service.AddMirrorCdn(reqMirrorCdn)
	if err != nil {
		log.Errorf("Insert data found error:%v", err)
	}
	c.JSON(200, gin.H{
		"res": err,
	})
}

func (app App) TestApi(c *gin.Context) {
	query := c.PostForm("query")
	data := service.TestApi(query)
	c.JSON(200, gin.H{
		"res": data,
	})
}

func (app App) SyncAllMirrors(c *gin.Context) {
	username := c.Param("username")
	log.Infof("User:%s trying sync all mirrors")
	index := app.cdnChecker.CheckAllMirrors(app.serverConfig.CdnChecker, username)
	c.JSON(http.StatusAccepted, gin.H{
		"index": index,
	})
}

func (app App) SyncMirror(c *gin.Context) {
	mirrorName := c.Param("mirror")
	username := c.Param("username")

	log.Infof("Username:%s, Mirror ID:%s", username, mirrorName)
	index := app.cdnChecker.CheckMirror(mirrorName, app.serverConfig.CdnChecker, username)
	c.JSON(http.StatusAccepted, gin.H{
		"index": index,
	})
}

func (app App) OperationHistory(c *gin.Context) {
	data := service.GetOperationsByDateDesc()
	c.JSON(http.StatusOK, gin.H{
		"history": data,
	})
}

func (app App) OperationHistoryByMirror(c *gin.Context) {
	mirror := c.Param("mirror")
	data := service.GetOperationsByMirror(mirror)
	c.JSON(http.StatusOK, gin.H{
		"history": data,
	})
}

func (app App) GetOperationByIndex(c *gin.Context) {
	index := c.Param("index")
	log.Info(index)
	data, err := service.GetOperationByIndex(index)
	if err != nil {
		log.Infof("%#v", err)
		c.JSON(http.StatusNoContent, gin.H{
			"msg": "get operation data found error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"operation": data,
	})
}

func main() {
	app := Init()
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins=[]string{app.serverConfig.Http.AllowOrigin}
	r.Use(cors.New(corsConfig))

	//r.GET("/mirrors", app.GetAllMirrors)
	//r.GET("/mirrors_cdn", app.GetAllMirrorsCdn)
	//
	//r.POST("/mirrors", app.AddMirror)
	//r.POST("/mirrors_cdn", app.AddMirrorCdn)
	//
	//r.POST("/test", app.TestApi)

	r.GET("/check/:username", app.SyncAllMirrors)

	r.GET("/check/:username/:mirror", app.SyncMirror)

	r.GET("/history", app.OperationHistory)

	r.GET("/history/:mirror", app.OperationHistoryByMirror)

	r.GET("/operation/:index", app.GetOperationByIndex)

	r.Run(":" + strconv.Itoa(app.serverConfig.Http.Port))
}
