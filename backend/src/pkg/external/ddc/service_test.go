package ddc_test

import (
	"bytes"
	"covid/src/pkg/entity"
	"covid/src/pkg/external/ddc"
	"covid/src/pkg/external/ddc/mocks"
	clientMock "covid/src/pkg/utils/client/mocks"
	"covid/src/pkg/utils/convert"
	covidError "covid/src/pkg/utils/error"
	loggerMocks "covid/src/pkg/utils/logger/mocks"
	"covid/src/pkg/utils/mocker"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	mockUseCase   *mocks.UseCase
	service       *ddc.Service
	clientService *clientMock.UseCase
	logs          *loggerMocks.Logger

	mockClientServiceRequest *mocker.MockCall
)

func callClientServiceRequest() *mock.Call {
	return clientService.On("Request", mock.Anything)
}

func beforeEach() {
	mockUseCase = &mocks.UseCase{}
	clientService = &clientMock.UseCase{}
	logs = &loggerMocks.Logger{}

	logs.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	service = &ddc.Service{
		ClientService: clientService,
		Logs:          logs,
	}
	service.UseCase = mockUseCase
}

func TestService_NewService(t *testing.T) {
	t.Run("should return struct dcc service when call new service", func(t *testing.T) {
		beforeEach()

		resService := ddc.NewService(clientService, logs)

		expectedService := &ddc.Service{
			ClientService: clientService,
			Logs:          logs,
		}

		expectedService.UseCase = expectedService

		assert.Equal(t, expectedService, resService)
	})
}

func TestService_GetDeathCaseRound4(t *testing.T) {
	beforeEachGetDeathCaseRound4 := func() {
		beforeEach()

		mockClientServiceRequest = mocker.NewMockCall(callClientServiceRequest)
		mockClientServiceRequest.Return(nil, nil)
	}

	t.Run("should call client service request when call service get death case round 4", func(t *testing.T) {
		beforeEachGetDeathCaseRound4()

		service.GetDeathCaseRound4()

		clientService.AssertCalled(t, "Request", mock.Anything)
	})

	t.Run("should return internal server error when call client service request failed", func(t *testing.T) {
		beforeEachGetDeathCaseRound4()
		mockClientServiceRequest.Return(nil, errors.New("error"))

		_, err := service.GetDeathCaseRound4()

		expectedError := convert.ValueToErrorCodePointer(covidError.InternalServerError)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should return internal server error when call client service request result is nil", func(t *testing.T) {
		beforeEachGetDeathCaseRound4()
		mockClientServiceRequest.Return(nil, nil)

		_, err := service.GetDeathCaseRound4()

		expectedError := convert.ValueToErrorCodePointer(covidError.InternalServerError)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should return internal server error when call client service request result status code is equal 400", func(t *testing.T) {
		beforeEachGetDeathCaseRound4()
		response := &http.Response{
			StatusCode: 400,
		}
		mockClientServiceRequest.Return(response, nil)

		_, err := service.GetDeathCaseRound4()

		expectedError := convert.ValueToErrorCodePointer(covidError.InternalServerError)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should return conflict error when decode client service response massage with wrong format struct", func(t *testing.T) {
		beforeEachGetDeathCaseRound4()
		wrongBody := io.NopCloser(strings.NewReader("wrongBody"))

		response := &http.Response{
			StatusCode: 200,
			Body:       wrongBody,
		}
		mockClientServiceRequest.Return(response, nil)

		_, err := service.GetDeathCaseRound4()

		expectedError := convert.ValueToErrorCodePointer(covidError.Conflict)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should return data covid case array when call service external ddc get death case round 4 success", func(t *testing.T) {
		beforeEachGetDeathCaseRound4()

		dataCovidCase := mockCovidCase()

		resposeFromAPI := ddc.ResponseFromAPICovid{
			Data: dataCovidCase,
		}

		jsonMock, _ := json.Marshal(resposeFromAPI)
		jsonMockBytes := bytes.NewBuffer([]byte(jsonMock))
		response := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(jsonMockBytes),
		}

		mockClientServiceRequest.Return(response, nil)

		res, err := service.GetDeathCaseRound4()

		assert.Equal(t, dataCovidCase, res)
		assert.Nil(t, err)
	})
}

func mockCovidCase() []entity.CovidCase {
	return []entity.CovidCase {
		{
			Year:         2021,
			Weeknum:      52,
			Province:     "กระบี่",
			Age:          "67",
			AgeRange:     "60-69 ปี",
			Occupation:   "พนักงานบริษัท/โรงงาน",
			Type:         "ผู้ป่วยยืนยัน",
			DeathCluster: convert.ValueToStringPointer("ครอบครัว"),
			UpdateDate:   "2024-05-27",
		},
		{
			Year:         2021,
			Weeknum:      52,
			Province:     "นครสวรรค์",
			Age:          "38",
			AgeRange:     "30-39 ปี",
			Occupation:   "รับจ้างทั่วไป / ฟรีแลนซ์",
			Type:         "ผู้ป่วยยืนยัน",
			DeathCluster: nil,
			UpdateDate:   "2024-05-27",
		},
		{
			Year:         2021,
			Weeknum:      52,
			Province:     "สตูล",
			Age:          "64",
			AgeRange:     "60-69 ปี",
			Occupation:   "เกษตรกร (ปลูกพืช)",
			Type:         "ผู้ป่วยยืนยัน",
			DeathCluster: nil,
			UpdateDate:   "2024-05-27",
		},
		{
			Year:         2021,
			Weeknum:      52,
			Province:     "ตรัง",
			Age:          "86",
			AgeRange:     ">= 70 ปี",
			Occupation:   "เกษตรกร (ปลูกพืช)",
			Type:         "ผู้ป่วยยืนยัน",
			DeathCluster: nil,
			UpdateDate:   "2024-05-27",
		},
	}
}