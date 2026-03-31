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
	cfg := config.MustLoad("game/config.yaml")

	db, err := dao.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer dao.CloseDB(db)
	log.Println("数据库初始化完成")

	// DAO
	userDAO := dao.NewUserDAO(db)
	scoreDAO := dao.NewScoreDAO(db)
	gameDAO := dao.NewGameDAO(db)
	apiLogDAO := dao.NewApiLogDAO(db)
	runDAO := dao.NewGameRunDAO(db)

	// 种子数据：按 game_key 补全缺失的默认游戏（含后续新增条目）
	if err := gameDAO.SeedDefaults(); err != nil {
		log.Printf("种子数据写入失败（可忽略）: %v", err)
	}

	// Service
	authSvc := service.NewAuthService(cfg.Auth)
	userSvc := service.NewUserService(userDAO, authSvc)
	scoreSvc := service.NewScoreService(db, scoreDAO, gameDAO, runDAO, cfg)
	gameSvc := service.NewGameService(gameDAO)

	// Handler
	authH := handler.NewAuthHandler(userSvc, cfg.Server.StaticDir)
	scoreH := handler.NewScoreHandler(scoreSvc, userSvc)
	gameH := handler.NewGameHandler(gameSvc)

	// Gin 引擎 & 路由
	engine := gin.Default()
	router.Setup(engine, &router.Deps{
		AuthSvc:   authSvc,
		ApiLogDAO: apiLogDAO,
		AuthH:     authH,
		ScoreH:    scoreH,
		GameH:     gameH,
	})

	fs := http.FileServer(http.Dir(cfg.Server.StaticDir))
	engine.NoRoute(gin.WrapH(fs))

	log.Printf("小游戏空间已启动，访问 http://localhost%s", cfg.Server.Addr)
	if err := engine.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
