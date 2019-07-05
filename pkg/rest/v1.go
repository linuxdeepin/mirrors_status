package rest

import (
	"github.com/gin-gonic/gin"
	"mirrors_status/internal/middleware"
	"mirrors_status/pkg/rest/controller"
)

func InitGuestController(engine *gin.Engine) {
	r := engine.Group("/api/v1")

	r.GET("/mirrors", controller.GetAllMirrors)
	r.GET("/upstreams", controller.GetAllUpstreams)
	r.GET("/publish", controller.GetPublishUpstreams)
	r.POST("/session", controller.Login)
}

func InitAdminController(engine *gin.Engine) {
	r := engine.Group("/api/v1/admin")
	r.Use(middleware.Auth())

	r.POST("/mirrors", controller.CreateMirror)
	r.DELETE("/mirrors/:id", controller.DeleteMirror)
	r.PUT("/mirrors", controller.UpdateMirror)
	r.POST("/tasks", controller.CreateTask)
	r.GET("/tasks/:id", controller.GetTaskById)
	r.GET("/tasks", controller.GetOpenTasks)
	r.GET("/check", controller.CheckAllMirrors)
	r.GET("/check/:upstream", controller.CheckMirrorsByUpstream)
	r.POST("/check", controller.CheckMirrors)
	r.DELETE("/tasks/:id", controller.AbortTask)
	r.PATCH("/tasks/:id/:status", controller.UpdateTaskStatus)
	r.GET("/archives/:id", controller.GetArchiveByTaskId)
	r.GET("/archives", controller.GetAllArchives)
	r.DELETE("/session", controller.Logout)
	r.GET("/session", controller.GetLoginStatus)
	r.POST("/mail", controller.SendMail)
}
