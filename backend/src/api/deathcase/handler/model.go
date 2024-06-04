package handler

import (
	deathCase "covid/src/pkg/deathcase"
	"covid/src/pkg/utils/logger"
)

type HandlerParam struct {
	Service deathCase.UseCase
	Logs    logger.Logger
}