package deathcase_test

import (
	"covid/src/pkg/deathcase"
	"covid/src/pkg/deathcase/mocks"
	"covid/src/pkg/entity"
	ddcMock "covid/src/pkg/external/ddc/mocks"
	"covid/src/pkg/utils/convert"
	covidError "covid/src/pkg/utils/error"
	loggerMocks "covid/src/pkg/utils/logger/mocks"
	"covid/src/pkg/utils/mocker"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	mockUseCase   *mocks.UseCase
	service       *deathcase.Service
	ddcService 	  *ddcMock.UseCase
	logs          *loggerMocks.Logger

	mockDDCServiceGetDeathCaseRound4 *mocker.MockCall
)

func callCDDCServiceGetDeathCaseRound4() *mock.Call {
	return ddcService.On("GetDeathCaseRound4", mock.Anything)
}

func beforeEach() {
	mockUseCase = &mocks.UseCase{}
	ddcService = &ddcMock.UseCase{}
	logs = &loggerMocks.Logger{}

	logs.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	service = &deathcase.Service{
		DDCService: ddcService,
		Logs:          logs,
	}
	service.UseCase = mockUseCase
}

func TestService_NewService(t *testing.T) {
	t.Run("should return struct deathcase service when call new service", func(t *testing.T) {
		beforeEach()

		resService := deathcase.NewService(ddcService, logs)

		expectedService := &deathcase.Service{
			DDCService: ddcService,
			Logs:          logs,
		}

		expectedService.UseCase = expectedService

		assert.Equal(t, expectedService, resService)
	})
}

func TestService_GetTopThreeProvice(t *testing.T) {
	beforeEachGetTopThreeProvice := func() {
		beforeEach()

		mockDDCServiceGetDeathCaseRound4 = mocker.NewMockCall(callCDDCServiceGetDeathCaseRound4)
		mockDDCServiceGetDeathCaseRound4.Return(nil, nil)
	}

	t.Run("should call ddc service get death case round 4 when call service get top three provice", func(t *testing.T) {
		beforeEachGetTopThreeProvice()

		service.GetTopThreeProvice()

		ddcService.AssertCalled(t, "GetDeathCaseRound4")
	})

	t.Run("should return internal server error when ddc service get death case round 4 return error", func(t *testing.T) {
		beforeEachGetTopThreeProvice()
		mockErr := covidError.InternalServerError
		mockDDCServiceGetDeathCaseRound4.Return(nil, &mockErr)

		_, err := service.GetTopThreeProvice()

		expectedError := convert.ValueToErrorCodePointer(covidError.InternalServerError)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should process data and return top 3 provinces with highest death case when ddc service get death case round 4 return success", func(t *testing.T) {
		beforeEachGetTopThreeProvice()

		mockDeathCaseRound4 := mockCovidCase()
		mockDDCServiceGetDeathCaseRound4.Return(mockDeathCaseRound4, nil)

		expected := []entity.TopThreeDeathCase{
			{Province: "กระบี่", Cases: 2},
			{Province: "สตูล", Cases: 2},
			{Province: "ตรัง", Cases: 2},
		}

		result, err := service.GetTopThreeProvice()

		assert.Equal(t, expected, result)
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
		{
			Year:         2021,
			Weeknum:      52,
			Province:     "กระบี่",
			Age:          "86",
			AgeRange:     ">= 70 ปี",
			Occupation:   "เกษตรกร (ปลูกพืช)",
			Type:         "ผู้ป่วยยืนยัน",
			DeathCluster: nil,
			UpdateDate:   "2024-05-27",
		},
	}
}