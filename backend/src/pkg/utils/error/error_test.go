package error_test

import (
	covid "covid/src/pkg/utils/error"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewError(t *testing.T) {

	t.Run("should return error struct when new error", func(t *testing.T) {
		err := covid.NewError(covid.InternalServerError, covid.InternalServerErrorMessage)

		expectedError := &covid.Error{Code: covid.InternalServerError, Message: covid.InternalServerErrorMessage}
		assert.Equal(t, expectedError, err)
	})

	t.Run("should return error with data when error has data", func(t *testing.T) {
		err := covid.NewError(covid.InternalServerError, covid.InternalServerErrorMessage, "some data")

		expectedError := &covid.Error{
			Code:    covid.InternalServerError,
			Message: covid.InternalServerErrorMessage,
			Data:    "some data",
		}
		assert.Equal(t, expectedError, err)
	})

}
