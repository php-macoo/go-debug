package dao

import (
	"go-debug/game/model"

	"gorm.io/gorm"
)

type GameDAO struct {
	db *gorm.DB
}

func NewGameDAO(db *gorm.DB) *GameDAO {
	return &GameDAO{db: db}
}

// GetOnlineByKey 返回已上线且 game_key 匹配的一条游戏，否则 gorm.ErrRecordNotFound。
func (d *GameDAO) GetOnlineByKey(gameKey string) (*model.Game, error) {
	var g model.Game
	err := d.db.Where("game_key = ? AND status = ?", gameKey, "online").First(&g).Error
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (d *GameDAO) ListOnline() ([]model.Game, error) {
	var games []model.Game
	err := d.db.Where("status = ?", "online").Order("sort_order ASC, id ASC").Find(&games).Error
	return games, err
}

// SeedDefaults 按 game_key 补全缺失的默认游戏（新环境全量插入，已有库只追加新游戏）。
func (d *GameDAO) SeedDefaults() error {
	defaults := []model.Game{
		{GameKey: "match3", Name: "三消达人", Icon: "🃏", Desc: "三张相同卡片即可消除，操作轻松但很上头。", URL: "/games/match3.html", Status: "online", SortOrder: 1},
		{GameKey: "2048", Name: "数字冲冲冲", Icon: "🔢", Desc: "合成更大的数字，冲击更高分数与步数效率。", URL: "/games/2048.html", Status: "online", SortOrder: 2},
		{GameKey: "sudoku", Name: "数独计划", Icon: "🧩", Desc: "经典逻辑推理，支持多难度与计时挑战。", URL: "/games/sudoku.html", Status: "online", SortOrder: 3},
		{GameKey: "snake", Name: "贪吃蛇", Icon: "🐍", Desc: "经典吃豆成长，速度渐快，挑战反应与走位。", URL: "/games/snake.html", Status: "online", SortOrder: 4},
		{GameKey: "plane", Name: "飞机大战", Icon: "✈️", Desc: "移动射击、躲避敌机，尽可能拿到更高分。", URL: "/games/plane.html", Status: "online", SortOrder: 5},
	}
	for i := range defaults {
		g := defaults[i]
		var n int64
		d.db.Model(&model.Game{}).Where("game_key = ?", g.GameKey).Count(&n)
		if n == 0 {
			if err := d.db.Create(&g).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
