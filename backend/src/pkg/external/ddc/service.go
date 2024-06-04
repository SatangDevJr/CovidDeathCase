package ddc

import (
	requestHeader "covid/src/api/requestheader"
	"covid/src/pkg/entity"
	"covid/src/pkg/utils/client"
	"covid/src/pkg/utils/convert"
	covidError "covid/src/pkg/utils/error"
	"covid/src/pkg/utils/logger"
	"encoding/json"
)

var (
	DDCURL  string
)

type UseCase interface {
	GetDeathCaseRound4() ([]entity.CovidCase, *covidError.ErrorCode)
}

type Service struct {
	UseCase
	ClientService client.UseCase
	Logs          logger.Logger
}

func NewService(clientService client.UseCase, logs logger.Logger) *Service {
	service := &Service{
		ClientService: clientService,
		Logs:          logs,
	}
	service.UseCase = service
	return service
}

func (service *Service) GetDeathCaseRound4() ([]entity.CovidCase, *covidError.ErrorCode) {

	header := make(map[string]string)
	header[requestHeader.ContentType] = requestHeader.ApplicationJson

	ddcAPIURLDeathCaseRound4 := DDCURL + "/Deaths/round-4-line-list"

	option := client.HttpOptions{
		Headers: header,
		URL:     ddcAPIURLDeathCaseRound4,
		Method:  "GET",
	}

	res, err := service.ClientService.Request(option)
	if err != nil {
		go service.Logs.Error("", "ddc_service_get_death_case_round_4_error", "", err)
		return nil, convert.ValueToErrorCodePointer(covidError.InternalServerError)
	}

	if res == nil || res.StatusCode != 200 {
		go service.Logs.Error("", "ddc_service_get_death_case_round_4_error", "", err)
		return nil, convert.ValueToErrorCodePointer(covidError.InternalServerError)
	}

	var resFromAPICovid ResponseFromAPICovid

	errBody := json.NewDecoder(res.Body).Decode(&resFromAPICovid)
	if errBody != nil {
		go service.Logs.Error("", "ddc_decode_get_death_case_round_4_error", "", errBody)
		return nil, convert.ValueToErrorCodePointer(covidError.Conflict)
	}

	covidCase := resFromAPICovid.Data
	
	return covidCase ,nil
}