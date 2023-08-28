package delete_segment_test

import (
	"avito_back_intership/internal/http-server/handlers/segment/delete_segment"
	"avito_back_intership/internal/http-server/handlers/segment/delete_segment/mocks"
	"avito_back_intership/internal/lib/api/slogdiscard"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockSegmentDeleter struct {
	DeleteSegmentFunc func(name string) error
	SegmentInfoFunc   func(segment string) (*sql.Rows, error)
	CreateLogFunc     func(user_id int, seg_name, operation string) error
}

func (m *MockSegmentDeleter) DeleteSegment(name string) error {
	if m.DeleteSegmentFunc != nil {
		return m.DeleteSegmentFunc(name)
	}
	return nil
}

func (m *MockSegmentDeleter) SegmentInfo(segment string) (*sql.Rows, error) {
	if m.SegmentInfoFunc != nil {
		return m.SegmentInfoFunc(segment)
	}
	return nil, nil
}

func (m *MockSegmentDeleter) CreateLog(user_id int, seg_name, operation string) error {
	if m.CreateLogFunc != nil {
		return m.CreateLogFunc(user_id, seg_name, operation)
	}
	return nil
}

func TestDeleteSegmentHandler(t *testing.T) {
	segment := "test_segment"

	t.Run("Successfully deleting segment", func(t *testing.T) {
		mockSegmentDeleter := mocks.NewSegmentDeleter(t)

		handler := delete_segment.New(slogdiscard.NewDiscardLogger(), mockSegmentDeleter)

		mockSegmentDeleter.On("SegmentInfo", segment).Return(nil, nil)

		req := httptest.NewRequest(http.MethodDelete, "/segment/"+segment, nil)
		rec := httptest.NewRecorder()
		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		// Validate the response if needed
	})

	t.Run("Failed to delete segment", func(t *testing.T) {
		errMessage := "failed to delete segment"

		handler := delete_segment.New(nil, &MockSegmentDeleter{
			DeleteSegmentFunc: func(name string) error {
				return errors.New(errMessage)
			},
			SegmentInfoFunc: func(segment string) (*sql.Rows, error) {
				return nil, nil
			},
		})

		req := httptest.NewRequest(http.MethodDelete, "/segment/"+segment, nil)
		rec := httptest.NewRecorder()
		handler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		// Validate the response if needed
	})

	t.Run("Segment not found", func(t *testing.T) {
		errMessage := "Segment not found"

		handler := delete_segment.New(nil, &MockSegmentDeleter{
			DeleteSegmentFunc: func(name string) error {
				return errors.New(errMessage)
			},
			SegmentInfoFunc: func(segment string) (*sql.Rows, error) {
				return nil, errors.New(errMessage)
			},
		})

		req := httptest.NewRequest(http.MethodDelete, "/segment/"+segment, nil)
		rec := httptest.NewRecorder()
		handler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		// Validate the response if needed
	})

	t.Run("Failed to get users in segment", func(t *testing.T) {
		errMessage := "failed to get users in segment " + segment

		handler := delete_segment.New(nil, &MockSegmentDeleter{
			SegmentInfoFunc: func(segment string) (*sql.Rows, error) {
				return nil, errors.New(errMessage)
			},
		})

		req := httptest.NewRequest(http.MethodDelete, "/segment/"+segment, nil)
		rec := httptest.NewRecorder()
		handler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		// Validate the response if needed
	})
}
