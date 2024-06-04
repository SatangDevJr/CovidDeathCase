package handler_test

import (
	"covid/src/api/deathcase/handler"
	requestHeader "covid/src/api/requestheader"
	"covid/src/pkg/deathcase/mocks"
	"covid/src/pkg/entity"
	"covid/src/pkg/utils/convert"
	covidError "covid/src/pkg/utils/error"
	loggerMocks "covid/src/pkg/utils/logger/mocks"
	"covid/src/pkg/utils/mocker"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	uri            string
	service        *mocks.UseCase
	deathCaseHandler *handler.DeathCaseHandler
	logs           *loggerMocks.Logger
	recorder       *httptest.ResponseRecorder
	request        *http.Request
	router         *mux.Router

	mockServicGetTopThreeProvice            *mocker.MockCall
)

func callServiceGetTopThreeProvice() *mock.Call {
	return service.On("GetTopThreeProvice")
}

func beforeEach() {
	uri = "/deathcase"
	service = &mocks.UseCase{}
	logs = &loggerMocks.Logger{}

	logs.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	deathCaseHandler = &handler.DeathCaseHandler{
		Service: service,
		Logs:    logs,
	}
}

func TestHandler_MakedeathcaseHandler(t *testing.T) {

	t.Run("should return struct death case handler when call make death case handler", func(t *testing.T) {
		beforeEach()
		handlerParam := handler.HandlerParam{
			Service: service,
			Logs:    logs,
		}

		handlerMakeDeathCase := handler.MakeDeathCaseHandler(handlerParam)

		expectedResult := &handler.DeathCaseHandler{
			Service: handlerParam.Service,
			Logs:    handlerParam.Logs,
		}
		assert.Equal(t, expectedResult, handlerMakeDeathCase)
	})
}

func TestHandler_GetTopThreeProvice(t *testing.T) {
	beforeEachGetTopThreeProvice := func() {
		beforeEach()
		uri += uri + "/gettopthreeprovice"
		router = mux.NewRouter()
		router.HandleFunc(uri, deathCaseHandler.GetTopThreeProvice)
		recorder = httptest.NewRecorder()
		request = httptest.NewRequest(http.MethodGet, uri, nil)

		mockServicGetTopThreeProvice = mocker.NewMockCall(callServiceGetTopThreeProvice)
		mockServicGetTopThreeProvice.Return(nil, nil)
	}

	t.Run("should call service search when request get top three provice", func(t *testing.T) {
		beforeEachGetTopThreeProvice()
		request = httptest.NewRequest(http.MethodGet, uri, nil)

		router.ServeHTTP(recorder, request)

		service.AssertCalled(t, "GetTopThreeProvice")
		assert.Equal(t, requestHeader.ApplicationJson, recorder.Header().Get(requestHeader.ContentType))
	})

	t.Run("should response internal server error when service get top three provice", func(t *testing.T) {
		beforeEachGetTopThreeProvice()
		err := covidError.InternalServerError
		mockServicGetTopThreeProvice.Return([]entity.TopThreeDeathCase{}, &err)

		router.ServeHTTP(recorder, request)

		expectedStatusCode, expectedError := covidError.MapMessageError(covidError.InternalServerError, "en")
		var body covidError.Error
		json.NewDecoder(recorder.Body).Decode(&body)
		assert.Equal(t, expectedError, body)
		assert.Equal(t, expectedStatusCode, recorder.Code)
		assert.Equal(t, requestHeader.ApplicationJson, recorder.Header().Get(requestHeader.ContentType))
	})

	t.Run("should response conflict error when service get top three provice failed", func(t *testing.T) {
		beforeEachGetTopThreeProvice()
		mockServicGetTopThreeProvice.Return(nil, convert.ValueToErrorCodePointer(covidError.Conflict))

		router.ServeHTTP(recorder, request)

		expectedStatusCode, expectedError := covidError.MapMessageError(covidError.Conflict, "en")
		var body covidError.Error
		json.NewDecoder(recorder.Body).Decode(&body)
		assert.Equal(t, expectedError, body)
		assert.Equal(t, expectedStatusCode, recorder.Code)
		assert.Equal(t, requestHeader.ApplicationJson, recorder.Header().Get(requestHeader.ContentType))
	})

	t.Run("should response get top three provice when request service get top three provice success", func(t *testing.T) {
		beforeEachGetTopThreeProvice()
		resdeathcase := []entity.TopThreeDeathCase{
			{
				Province: "กรุงเทพมหานคร",
				Cases: 164,
			},
			{
				Province: "ชลบุรี",
				Cases: 43,
			},
			{
				Province: "เชียงใหม่",
				Cases: 42,
			},
		}
		mockServicGetTopThreeProvice.Return(resdeathcase, nil)

		router.ServeHTTP(recorder, request)

		expectedBody := resdeathcase

		var body []entity.TopThreeDeathCase
		json.NewDecoder(recorder.Body).Decode(&body)
		assert.Equal(t, expectedBody, body)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, requestHeader.ApplicationJson, recorder.Header().Get(requestHeader.ContentType))
	})
}