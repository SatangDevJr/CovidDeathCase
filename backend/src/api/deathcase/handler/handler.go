package handler

import (
	requestHeader "covid/src/api/requestheader"
	deathCase "covid/src/pkg/deathcase"
	covidError "covid/src/pkg/utils/error"
	"covid/src/pkg/utils/logger"
	"encoding/json"
	"net/http"
)

type DeathCaseHandler struct {
	Service deathCase.UseCase
	Logs    logger.Logger
}

func MakeDeathCaseHandler(handlerParam HandlerParam) *DeathCaseHandler {
	return &DeathCaseHandler{
		Service: handlerParam.Service,
		Logs:    handlerParam.Logs,
	}
}

func (handler *DeathCaseHandler) GetTopThreeProvice(response http.ResponseWriter, request *http.Request) {

	response.Header().Set(requestHeader.ContentType, requestHeader.ApplicationJson)

	res, err := handler.Service.GetTopThreeProvice()
	if err != nil {
		switch *err {
		case covidError.Conflict:
			go handler.Logs.Error(request.URL.Path, "deathCase_handler_getTopThreeProvice_conflict", nil, err)
		default:
			go handler.Logs.Error(request.URL.Path, "deathCase_handler_getTopThreeProvice_internalServerError", nil, err)
		}
		statusCode, errMsg := covidError.MapMessageError(*err, "en")
		response.WriteHeader(statusCode)
		json.NewEncoder(response).Encode(&errMsg)
		return
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(&res)
}