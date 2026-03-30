// Package router 负责注册所有 HTTP 路由，按公开/需认证两组划分。
package router

import (
	"go-debug/game/handler"
	"go-debug/game/middleware"
	"go-debug/game/service"

	"github.com/gin-gonic/gin"
)

// Setup 将所有 API 路由挂载到 Gin 引擎上。
//
// 路由分组:
//   - 公开路由：/api/register, /api/login
//   - 需认证路由：/api/user, /api/user/profile, /api/user/avatar, /api/score, /api/leaderboard
func Setup(engine *gin.Engine, authSvc *service.AuthService, authH *handler.AuthHandler, scoreH *handler.ScoreHandler) {
	api := engine.Group("/api")
	{
		api.POST("/register", authH.Register)
		api.POST("/login", authH.Login)
	}

	authed := api.Group("", middleware.Auth(authSvc))
	{
		authed.GET("/user", authH.GetUser)
		authed.PUT("/user/profile", authH.UpdateProfile)
		authed.POST("/user/avatar", authH.UploadAvatar)
		authed.PUT("/user/avatar", authH.SetAvatar)
		authed.POST("/score", scoreH.Submit)
		authed.GET("/leaderboard", scoreH.Leaderboard)
	}
}
