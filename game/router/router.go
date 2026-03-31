// Package router 按模块划分注册所有 HTTP 路由。
package router

import (
	"go-debug/game/dao"
	"go-debug/game/handler"
	"go-debug/game/middleware"
	"go-debug/game/service"

	"github.com/gin-gonic/gin"
)

// Deps 汇总路由层所需的全部依赖。
type Deps struct {
	AuthSvc   *service.AuthService
	ApiLogDAO *dao.ApiLogDAO
	AuthH     *handler.AuthHandler
	ScoreH    *handler.ScoreHandler
	GameH     *handler.GameHandler
}

// Setup 将所有 API 路由按模块挂载到 Gin 引擎上。
//
// 路由分组:
//   - /api/games              公开 - 游戏大厅列表
//   - /api/auth/*             公开 - 注册/登录
//   - /api/user/*             需认证 - 用户资料
//   - /api/game/:gameKey/*    可选认证 - 登录或 X-Guest-Device-Id 匿名：run/start、交成绩；排行榜公开
func Setup(engine *gin.Engine, deps *Deps) {
	api := engine.Group("/api")
	api.Use(middleware.AccessLog(deps.ApiLogDAO))

	// ─── 公开: 游戏大厅 ───
	api.GET("/games", deps.GameH.List)

	// ─── 公开: 认证模块 ───
	auth := api.Group("/auth")
	{
		auth.POST("/register", deps.AuthH.Register)
		auth.POST("/login", deps.AuthH.Login)
	}

	// ─── 需认证: 用户模块 ───
	user := api.Group("/user", middleware.Auth(deps.AuthSvc))
	{
		user.GET("", deps.AuthH.GetUser)
		user.PUT("/profile", deps.AuthH.UpdateProfile)
		user.POST("/avatar", deps.AuthH.UploadAvatar)
		user.PUT("/avatar", deps.AuthH.SetAvatar)
	}

	// ─── 游戏模块：可选登录；未登录时须由前端带 X-Guest-Device-Id（见 auth.js）───
	game := api.Group("/game", middleware.OptionalAuth(deps.AuthSvc))
	{
		game.POST("/:gameKey/run/start", deps.ScoreH.StartRun)
		game.POST("/:gameKey/score", deps.ScoreH.Submit)
		game.GET("/:gameKey/leaderboard", deps.ScoreH.Leaderboard)
	}
}
