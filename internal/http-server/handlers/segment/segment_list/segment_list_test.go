package segment_list_test

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"avito_back_intership/internal/http-server/handlers/segment/segment_list"
	mock "avito_back_intership/internal/http-server/handlers/segment/segment_list/mocks"
	"avito_back_intership/internal/lib/api/response"
	"avito_back_intership/internal/lib/api/slogdiscard"
)

func TestNew(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger()
	mockSegmentListGetter := &mock.SegmentListGetter{}

	t.Run("success", func(t *testing.T) {
		// Создаем ожидаемый результат
		expectedResponse := segment_list.Response{
			SegmentList: []string{"segment1", "segment2"},
			Response:    response.OK(),
		}
		_ = expectedResponse

		// Создаем mock результата выполнения функции SegmentList()
		rows := &sql.Rows{}
		mockSegmentListGetter.On("SegmentList").Return(rows, nil)

		// Создаем новый HTTP запрос
		req, err := http.NewRequest("GET", "/segment", nil)
		assert.NoError(t, err)

		// Создаем новый HTTP response recorder
		rr := httptest.NewRecorder()

		fmt.Printf("rr.Body: %v\n", rr.Body)
		fmt.Printf("rr.Code: %v\n", rr.Code)

		// Вызываем хэндлер
		handler := segment_list.New(mockLogger, mockSegmentListGetter)
		fmt.Printf("handler: %v\n", handler)

		handler.ServeHTTP(rr, req)

		fmt.Printf("rr.Body: %v\n", rr.Body)
		fmt.Printf("rr.Code: %v\n", rr.Code)

		// Проверяем код ответа
		assert.Equal(t, http.StatusOK, rr.Code)

		// Проверяем ожидаемый результат
		// var actualResponse segment_list.Response
		// responseBody := rr.Body.Bytes()
		// err = response.Parse(responseBody, &actualResponse)
		// assert.NoError(t, err)
		// assert.Equal(t, expectedResponse, actualResponse)

		// mockLogger.AssertExpectations(t)
		// mockSegmentListGetter.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		// Создаем mock ошибки
		expectedError := errors.New("failed to get list of segments")
		mockSegmentListGetter.On("SegmentList").Return(nil, expectedError)

		// Создаем новый HTTP запрос
		req, err := http.NewRequest("GET", "/segment", nil)
		assert.NoError(t, err)

		// Создаем новый HTTP response recorder
		rr := httptest.NewRecorder()

		// Вызываем хэндлер
		handler := segment_list.New(mockLogger, mockSegmentListGetter)
		handler(rr, req)

		// Проверяем код ответа
		assert.Equal(t, rr.Code, http.StatusOK)

		// Проверяем ошибку в ответе
		actualResponse := response.Error(expectedError.Error())

		assert.Equal(t, response.Error(expectedError.Error()), actualResponse)

		mockSegmentListGetter.AssertExpectations(t)
	})
}
