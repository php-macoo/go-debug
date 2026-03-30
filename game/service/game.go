package service

import (
	"go-debug/game/dao"
	"go-debug/game/model"
)

type GameService struct {
	gameDAO *dao.GameDAO
}

func NewGameService(gameDAO *dao.GameDAO) *GameService {
	return &GameService{gameDAO: gameDAO}
}

func (s *GameService) ListOnline() ([]model.Game, error) {
	return s.gameDAO.ListOnline()
}
