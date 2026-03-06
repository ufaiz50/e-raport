package api

import (
	"bytes"
	"context"
	"encoding/json"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateStudent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	ctx := context.Background()
	repo := NewStudentRepository(mockDB, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/students", repo.CreateStudent)

	input := models.CreateStudent{Name: "Umar", Email: "umar@example.com", Type: "junior"}
	body, _ := json.Marshal(input)

	mockDB.EXPECT().Create(gomock.Any()).DoAndReturn(func(student *models.Student) *gorm.DB {
		return &gorm.DB{Error: nil}
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "umar@example.com")
	assert.Contains(t, w.Body.String(), "junior")
}

func TestFindStudent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	ctx := context.Background()
	repo := NewStudentRepository(mockDB, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/students/:id", repo.FindStudent)

	expected := models.Student{ID: 1, Name: "Umar", Email: "umar@example.com", Type: "junior"}

	mockDB.EXPECT().Where("id = ? AND school_id = ?", "1", uint(1)).Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).DoAndReturn(func(dest interface{}, conds ...interface{}) database.Database {
		if s, ok := dest.(*models.Student); ok {
			*s = expected
		}
		return mockDB
	})
	mockDB.EXPECT().Error().Return(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/students/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "umar@example.com")
	assert.Contains(t, w.Body.String(), "junior")
}
