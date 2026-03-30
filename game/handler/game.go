package handler

import (
	"go-debug/game/pkg/resp"
	"go-debug/game/service"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameSvc *service.GameService
}

// gameListItem 大厅列表对外字段（不含库内主键、gameKey、排序与时间戳）。
type gameListItem struct {
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Desc   string `json:"desc"`
	URL    string `json:"url"`
	Status string `json:"status"`
}

func NewGameHandler(gameSvc *service.GameService) *GameHandler {
	return &GameHandler{gameSvc: gameSvc}
}

// List 返回所有上线状态的游戏列表，供大厅首页展示。
// GET /api/games
func (h *GameHandler) List(c *gin.Context) {
	games, err := h.gameSvc.ListOnline()
	if err != nil {
		resp.Fail500(c, "查询失败")
		return
	}
	list := make([]gameListItem, len(games))
	for i := range games {
		g := games[i]
		list[i] = gameListItem{
			Name:   g.Name,
			Icon:   g.Icon,
			Desc:   g.Desc,
			URL:    g.URL,
			Status: g.Status,
		}
	}
	resp.OK(c, gin.H{"list": list})
}
