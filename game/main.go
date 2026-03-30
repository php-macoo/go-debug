// Package main 是三消达人游戏的入口，负责初始化各层依赖并启动 HTTP 服务。
package main

import (
	"log"
	"net/http"

	"go-debug/game/config"
	"go-debug/game/dao"
	"go-debug/game/handler"
	"go-debug/game/router"
	"go-debug/game/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置（数据库、服务端口、认证密钥等）
	cfg := config.MustLoad("game/config.yaml")

	// 初始化数据库连接并自动迁移表结构
	db, err := dao.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer dao.CloseDB(db)
	log.Println("数据库初始化完成")

	// 依赖注入: DAO → Service → Handler
	userDAO := dao.NewUserDAO(db)
	scoreDAO := dao.NewScoreDAO(db)

	authSvc := service.NewAuthService(cfg.Auth)
	userSvc := service.NewUserService(userDAO, authSvc)
	scoreSvc := service.NewScoreService(scoreDAO)

	authH := handler.NewAuthHandler(userSvc, cfg.Server.StaticDir)
	scoreH := handler.NewScoreHandler(scoreSvc)

	// 创建 Gin 引擎，注册 API 路由
	engine := gin.Default()
	router.Setup(engine, authSvc, authH, scoreH)

	// 未匹配的路由交给静态文件服务（前端页面）
	fs := http.FileServer(http.Dir(cfg.Server.StaticDir))
	engine.NoRoute(gin.WrapH(fs))

	log.Printf("三消达人已启动，访问 http://localhost%s", cfg.Server.Addr)
	if err := engine.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
