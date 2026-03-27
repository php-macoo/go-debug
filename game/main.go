package main

import (
	"log"
	"net/http"

	"go-debug/game/config"
	"go-debug/game/dao"
	"go-debug/game/handler"
	"go-debug/game/service"
)

func main() {
	cfg := config.MustLoad("game/config.yaml")

	db, err := dao.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer dao.CloseDB(db)
	log.Println("数据库初始化完成")

	userDAO := dao.NewUserDAO(db)
	scoreDAO := dao.NewScoreDAO(db)

	authSvc := service.NewAuthService(cfg.Auth)
	userSvc := service.NewUserService(userDAO, authSvc)
	scoreSvc := service.NewScoreService(scoreDAO)

	h := handler.New(userSvc, scoreSvc, authSvc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	mux.Handle("/", http.FileServer(http.Dir(cfg.Server.StaticDir)))

	log.Printf("三消达人已启动，访问 http://localhost%s", cfg.Server.Addr)
	if err := http.ListenAndServe(cfg.Server.Addr, mux); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
