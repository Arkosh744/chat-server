package chat

import "go.uber.org/zap"

var _ Service = (*service)(nil)

type Service interface {
}

type service struct {
	log *zap.SugaredLogger
}

func NewService(log *zap.SugaredLogger) *service {
	return &service{
		log: log,
	}
}
